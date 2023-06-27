package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"table": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"run_date_offset", "start_date", "end_date"},
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
				Optional: true,
				Default:  true,
			},
			"include_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"cluster_property_sets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"additional_spark_properties": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					return make(map[string]interface{}), nil
				},
			},
			"entity_population": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Labels to attach to the object",
				Elem:        labelSchema(),
			},
			"attribute": {
				Type:        schema.TypeSet,
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

	if FeatureStore.Type == "batch" {
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
	}
	if FeatureStore.Type == "streaming" {
		if FeatureStore.Table == nil {
			return errors.New("Required field is missing for streaming feature store: table")
		}
		if err := d.Set("table", *FeatureStore.Table); err != nil {
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
	if err := d.Set("include_metadata", FeatureStore.IncludeMetadata); err != nil {
		return err
	}
	if err := d.Set("destination", destinations); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(FeatureStore.Cluster)); err != nil {
		return err
	}
	if err := d.Set("cluster_property_sets", identifierList(FeatureStore.ClusterPropertySets)); err != nil {
		return err
	}
	if err := d.Set("additional_spark_properties", FeatureStore.AdditionalSparkProperties); err != nil {
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
	} else {
		if err := d.Set("commit_target", nil); err != nil {
			return err
		}
		if err := d.Set("branch_target", nil); err != nil {
			return err
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

	source := d.Get("additional_spark_properties").(map[string]interface{})
	additionalSparkProperties := make(map[string]string)

	for k, v := range source {
		additionalSparkProperties[k] = v.(string)
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

	destinations, err := expandDestinationReferences(d.Get("destination").([]interface{}))
	if err != nil {
		return nil, err
	}

	featureStore := FeatureStore{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		FeatureSet:                featureSet,
		Principal:                 principal,
		Enabled:                   d.Get("enabled").(bool),
		Destinations:              destinations,
		Cluster:                   cluster,
		ClusterPropertySets:       expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
		AdditionalSparkProperties: additionalSparkProperties,
		Population:                population,
		Schedule:                  schedule,
		Labels:                    expandLabels(d),
		Attributes:                expandAttributes(d),
		IncludeMetadata:           d.Get("include_metadata").(bool),
		VersionTarget:             versionTarget,
	}

	table := getNullableInt(d, "table")
	if table != nil {
		featureStore.Type = "streaming"
		featureStore.Table = table
	} else {
		featureStore.Type = "batch"
		featureStore.RunDateOffset = getNullableInt(d, "run_date_offset")
		featureStore.StartDate = getNullableString(d, "start_date")
		featureStore.EndDate = getNullableString(d, "end_date")
	}

	return &featureStore, nil
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
