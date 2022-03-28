package anaml

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const featureStoreDescription = `
# Feature Stores (Schedules)

A Schedule is the output of a feature set run at a specific time and written to one or more destination.
Generally the output in the Destination will be a table with a timestamp, entity identifier, and one
column per features in the Feature Set.

Schedules can either be a historical run which covers a range of dates, or a daily run where new data
is generated on a daily basis.

An entity population can be used to further refine the entities and dates for feature generation.
`

func ResourceFeatureStore() *schema.Resource {
	return &schema.Resource{
		Description: featureStoreDescription,
		Create:      resourceFeatureStoreCreate,
		Read:        resourceFeatureStoreRead,
		Update:      resourceFeatureStoreUpdate,
		Delete:      resourceFeatureStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"run_date_offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"feature_set": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"principal": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"daily_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          dailyScheduleSchema(),
				ConflictsWith: []string{"cron_schedule"},
			},
			"cron_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          cronScheduleSchema(),
				ConflictsWith: []string{"daily_schedule"},
			},
			"destination": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     destinationSchema(),
			},
			"cluster": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"entity_population": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Labels to attach to the object",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"attribute": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
			"commit_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Commit to run feature set (and population) for.",
			},
			"branch_target": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Branch to run feature set (and population) for.",
				ConflictsWith: []string{"commit_target"},
			},
		},
	}
}

func destinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"folder": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     folderDestinationSchema(),
			},
			"table": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     tableDestinationSchema(),
			},
			"topic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     topicDestinationSchema(),
			},
		},
	}
}

func folderDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"partitioning_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func tableDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func topicDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceFeatureStoreRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStoreID := d.Id()

	FeatureStore, err := c.GetFeatureStore(FeatureStoreID)
	if err != nil {
		return err
	}
	if FeatureStore == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", FeatureStore.Name); err != nil {
		return err
	}
	if err := d.Set("description", FeatureStore.Description); err != nil {
		return err
	}
	if FeatureStore.StartDate != nil {
		if err := d.Set("start_date", *FeatureStore.StartDate); err != nil {
			return err
		}
	}
	if FeatureStore.RunDateOffset != nil {
		if err := d.Set("run_date_offset", *FeatureStore.RunDateOffset); err != nil {
			return err
		}
	}
	if FeatureStore.EndDate != nil {
		if err := d.Set("end_date", *FeatureStore.EndDate); err != nil {
			return err
		}
	}
	if FeatureStore.Principal != nil {
		if err := d.Set("principal", strconv.Itoa(*FeatureStore.Principal)); err != nil {
			return err
		}
	}

	destinations, err := flattenDestinationReferences(FeatureStore.Destinations)
	if err != nil {
		return err
	}

	if err := d.Set("feature_set", strconv.Itoa(FeatureStore.FeatureSet)); err != nil {
		return err
	}
	if err := d.Set("enabled", FeatureStore.Enabled); err != nil {
		return err
	}
	if err := d.Set("destination", destinations); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(FeatureStore.Cluster)); err != nil {
		return err
	}
	if err := d.Set("labels", FeatureStore.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(FeatureStore.Attributes)); err != nil {
		return err
	}
	if FeatureStore.Population != nil {
		if err := d.Set("entity_population", strconv.Itoa(*FeatureStore.Population)); err != nil {
			return err
		}
	}

	if FeatureStore.Schedule.Type == "daily" {
		dailySchedules, err := parseDailySchedule(FeatureStore.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("daily_schedule", dailySchedules); err != nil {
			return err
		}
	}

	if FeatureStore.Schedule.Type == "cron" {
		cronSchedules, err := parseCronSchedule(FeatureStore.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("cron_schedule", cronSchedules); err != nil {
			return err
		}
	}

	if FeatureStore.VersionTarget != nil {
		if FeatureStore.VersionTarget.Commit != nil {
			if err := d.Set("commit_target", FeatureStore.VersionTarget.Commit); err != nil {
				return err
			}
			if err := d.Set("branch_target", nil); err != nil {
				return err
			}
		} else if FeatureStore.VersionTarget.Branch != nil {
			if err := d.Set("branch_target", FeatureStore.VersionTarget.Branch); err != nil {
				return err
			}
			if err := d.Set("commit_target", nil); err != nil {
				return err
			}
		}
	}

	return err
}

func resourceFeatureStoreCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStore, err := composeFeatureStore(d)
	if err != nil {
		return err
	}

	e, err := c.CreateFeatureStore(*FeatureStore)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceFeatureStoreUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStoreID := d.Id()
	FeatureStore, err := composeFeatureStore(d)
	if err != nil {
		return err
	}

	err = c.UpdateFeatureStore(FeatureStoreID, *FeatureStore)
	if err != nil {
		return err
	}

	return nil
}

func composeFeatureStore(d *schema.ResourceData) (*FeatureStore, error) {
	featureSet, err := strconv.Atoi(d.Get("feature_set").(string))
	if err != nil {
		return nil, err
	}

	var principal (*int) = nil
	principalRaw, principalOk := d.GetOk("principal")
	if principalOk {
		principal_, err := strconv.Atoi(principalRaw.(string))
		if err != nil {
			return nil, err
		}
		principal = &principal_
	}

	cluster, err := strconv.Atoi(d.Get("cluster").(string))
	if err != nil {
		return nil, err
	}

	var population (*int) = nil
	if d.Get("entity_population").(string) != "" {
		population_, err := strconv.Atoi(d.Get("entity_population").(string))
		if err != nil {
			return nil, err
		}
		population = &population_
	}

	var schedule = composeNeverSchedule()
	if dailySchedule, _ := expandSingleMap(d.Get("daily_schedule")); dailySchedule != nil {
		schedule, err = composeDailySchedule(dailySchedule)
		if err != nil {
			return nil, err
		}
	}
	if cronSchedule, _ := expandSingleMap(d.Get("cron_schedule")); cronSchedule != nil {
		schedule, err = composeCronSchedule(cronSchedule)
		if err != nil {
			return nil, err
		}
	}
	var versionTarget (*VersionTarget) = nil
	if commit, _ := d.Get("commit_target").(string); commit != "" {
		versionTarget = &VersionTarget{
			Type:   "commit",
			Commit: &commit,
		}
	}
	if branch, _ := d.Get("branch_target").(string); branch != "" {
		versionTarget = &VersionTarget{
			Type:   "branch",
			Branch: &branch,
		}
	}

	destinations, err := expandDestinationReferences(d)
	if err != nil {
		return nil, err
	}

	return &FeatureStore{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		RunDateOffset: getNullableInt(d, "run_date_offset"),
		StartDate:     getNullableString(d, "start_date"),
		EndDate:       getNullableString(d, "end_date"),
		FeatureSet:    featureSet,
		Principal:     principal,
		Enabled:       d.Get("enabled").(bool),
		Destinations:  destinations,
		Cluster:       cluster,
		Population:    population,
		Schedule:      schedule,
		Labels:        expandStringList(d.Get("labels").([]interface{})),
		Attributes:    expandAttributes(d),
		VersionTarget: versionTarget,
	}, nil
}

func resourceFeatureStoreDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStoreID := d.Id()

	err := c.DeleteFeatureStore(FeatureStoreID)
	if err != nil {
		return err
	}

	return nil
}

func expandDestinationReferences(d *schema.ResourceData) ([]DestinationReference, error) {
	drs := d.Get("destination").([]interface{})
	res := make([]DestinationReference, 0, len(drs))

	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})

		destID, _ := strconv.Atoi(val["destination"].(string))
		parsed := DestinationReference{
			DestinationID: destID,
		}

		if folder, _ := expandSingleMap(val["folder"]); folder != nil {
			if path, ok := folder["path"].(string); ok {
				parsed.Type = "folder"
				parsed.Folder = path
				enabled := folder["partitioning_enabled"].(bool)
				parsed.FolderPartitioningEnabled = &enabled
			} else {
				return nil, fmt.Errorf("error casting table.path %i", folder["path"])
			}
		}

		if table, _ := expandSingleMap(val["table"]); table != nil {
			if tableName, ok := table["name"].(string); ok {
				parsed.Type = "table"
				parsed.TableName = tableName
			} else {
				return nil, fmt.Errorf("error casting table.name %i", table["name"])
			}
		}

		if topic, _ := expandSingleMap(val["topic"]); topic != nil {
			if topicName, ok := topic["name"].(string); ok {
				parsed.Type = "topic"
				parsed.Topic = topicName
			} else {
				return nil, fmt.Errorf("error casting topic.name %i", topic["name"])
			}
		}

		res = append(res, parsed)
	}

	return res, nil
}

func flattenDestinationReferences(destinations []DestinationReference) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0, len(destinations))

	for _, destination := range destinations {
		single := make(map[string]interface{})
		single["destination"] = strconv.Itoa(destination.DestinationID)

		if destination.Type == "folder" {
			folder := make(map[string]interface{})
			folder["path"] = destination.Folder
			folder["partitioning_enabled"] = destination.FolderPartitioningEnabled

			folders := make([]map[string]interface{}, 0, 1)
			folders = append(folders, folder)
			single["folder"] = folders
		}

		if destination.Type == "table" {
			table := make(map[string]interface{})
			table["name"] = destination.TableName

			tables := make([]map[string]interface{}, 0, 1)
			tables = append(tables, table)
			single["table"] = tables
		}

		if destination.Type == "topic" {
			topic := make(map[string]interface{})
			topic["name"] = destination.Topic

			topics := make([]map[string]interface{}, 0, 1)
			topics = append(topics, topic)
			single["topic"] = topics
		}

		res = append(res, single)
	}

	return res, nil
}
