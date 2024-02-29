package anaml

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const viewMaterialisationDescription = `
# View Materialisation

View materialisation offer a simple way of scheduling arbitrary SQL table expressions to be executed and written to a destination.
`

func ResourceViewMaterialisationJob() *schema.Resource {
	return &schema.Resource{
		Description: viewMaterialisationDescription,
		Create:      resourceViewMaterialisationJobCreate,
		Read:        resourceViewMaterialisationJobRead,
		Update:      resourceViewMaterialisationJobUpdate,
		Delete:      resourceViewMaterialisationJobDelete,
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
			"principal": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"include_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"daily_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          dailyScheduleSchema(),
				ConflictsWith: []string{"cron_schedule", "dependency_schedule"},
			},
			"cron_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          cronScheduleSchema(),
				ConflictsWith: []string{"daily_schedule", "dependency_schedule"},
			},
			"dependency_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          dependencyScheduleSchema(),
				ConflictsWith: []string{"daily_schedule", "cron_schedule"},
			},
			"usagettl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"view": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     viewMaterialisationSpecSchema(),
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
				Description: "Commit to run view materialisation for.",
			},
			"branch_target": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Branch to run view materialisation for.",
				ConflictsWith: []string{"commit_target"},
			},
		},
	}
}

func viewMaterialisationSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     destinationSchema(),
			},
		},
	}
}

func resourceViewMaterialisationJobRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	ViewMaterialisationJobID := d.Id()

	ViewMaterialisationJob, err := c.GetViewMaterialisationJob(ViewMaterialisationJobID)
	if err != nil {
		return err
	}
	if ViewMaterialisationJob == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", ViewMaterialisationJob.Name); err != nil {
		return err
	}
	if err := d.Set("description", ViewMaterialisationJob.Description); err != nil {
		return err
	}

	if ViewMaterialisationJob.Type == "batch" {
		daily, cron, dependency, err := parseSchedule(ViewMaterialisationJob.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("daily_schedule", daily); err != nil {
			return err
		}
		if err := d.Set("cron_schedule", cron); err != nil {
			return err
		}
		if err := d.Set("dependency_schedule", dependency); err != nil {
			return err
		}
	}

	if ViewMaterialisationJob.Principal != nil {
		if err := d.Set("principal", strconv.Itoa(*ViewMaterialisationJob.Principal)); err != nil {
			return err
		}
	}

	views, err := flattenViewMaterialisationSpec(ViewMaterialisationJob.Views)
	if err != nil {
		return err
	}

	if err := d.Set("view", views); err != nil {
		return err
	}

	if ViewMaterialisationJob.UsageTTL != nil {
		if err := d.Set("usagettl", ViewMaterialisationJob.UsageTTL); err != nil {
			return err
		}
	} else {
		d.Set("usagettl", nil)
	}
	if err := d.Set("cluster", strconv.Itoa(ViewMaterialisationJob.Cluster)); err != nil {
		return err
	}
	if err := d.Set("cluster_property_sets", identifierList(ViewMaterialisationJob.ClusterPropertySets)); err != nil {
		return err
	}
	if err := d.Set("additional_spark_properties", ViewMaterialisationJob.AdditionalSparkProperties); err != nil {
		return err
	}
	if err := d.Set("labels", ViewMaterialisationJob.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(ViewMaterialisationJob.Attributes)); err != nil {
		return err
	}

	if ViewMaterialisationJob.VersionTarget != nil {
		if ViewMaterialisationJob.VersionTarget.Commit != nil {
			if err := d.Set("commit_target", ViewMaterialisationJob.VersionTarget.Commit); err != nil {
				return err
			}
			if err := d.Set("branch_target", nil); err != nil {
				return err
			}
		} else if ViewMaterialisationJob.VersionTarget.Branch != nil {
			if err := d.Set("branch_target", ViewMaterialisationJob.VersionTarget.Branch); err != nil {
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

func resourceViewMaterialisationJobCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	ViewMaterialisationJob, err := composeViewMaterialisationJob(d)
	if err != nil {
		return err
	}

	e, err := c.CreateViewMaterialisationJob(*ViewMaterialisationJob)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceViewMaterialisationJobUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	ViewMaterialisationJobID := d.Id()
	vm, err := composeViewMaterialisationJob(d)
	if err != nil {
		return err
	}

	err = c.UpdateViewMaterialisationJob(ViewMaterialisationJobID, *vm)
	if err != nil {
		return err
	}

	return nil
}

func composeViewMaterialisationJob(d *schema.ResourceData) (*ViewMaterialisationJob, error) {
	principal := getAnamlIdPointer(d, "principal")
	cluster, err := strconv.Atoi(d.Get("cluster").(string))
	if err != nil {
		return nil, err
	}

	source := d.Get("additional_spark_properties").(map[string]interface{})
	additionalSparkProperties := make(map[string]string)

	for k, v := range source {
		additionalSparkProperties[k] = v.(string)
	}

	usageTTL := getStringPointer(d, "usagettl")

	views, err := expandViewMaterialisationSpec(d)
	if err != nil {
		return nil, err
	}
	schedule, err := composeSchedule(d)
	if err != nil {
		return nil, err
	}
	versionTarget := composeVersionTarget(d)

	viewMateralisationJob := ViewMaterialisationJob{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		Type:                      "batch",
		Schedule:                  schedule,
		IncludeMetadata:           d.Get("include_metadata").(bool),
		Principal:                 principal,
		UsageTTL:                  usageTTL,
		Views:                     views,
		Cluster:                   cluster,
		ClusterPropertySets:       expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
		AdditionalSparkProperties: additionalSparkProperties,
		Labels:                    expandLabels(d),
		Attributes:                expandAttributes(d),
		VersionTarget:             versionTarget,
	}

	return &viewMateralisationJob, nil
}

func resourceViewMaterialisationJobDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	id := d.Id()

	err := c.DeleteViewMaterialisationJob(id)
	if err != nil {
		return err
	}

	return nil
}

func expandViewMaterialisationSpec(d *schema.ResourceData) ([]ViewMaterialisationSpec, error) {
	views := d.Get("view").([]interface{})
	res := make([]ViewMaterialisationSpec, 0, len(views))

	for _, view := range views {
		val, _ := view.(map[string]interface{})
		table, err := strconv.Atoi(val["table"].(string))
		if err != nil {
			return nil, err
		}

		destinations, err := expandDestinationReferences(val["destination"].([]interface{}))
		if err != nil {
			return nil, err
		}

		if len(destinations) != 1 {
			return nil, fmt.Errorf("Incorrect number of destinations parsed. This is an internal error to the provider.")
		}

		viewMaterialisationSpec := ViewMaterialisationSpec{
			Table:       table,
			Destination: destinations[0],
		}

		res = append(res, viewMaterialisationSpec)
	}

	return res, nil
}

func flattenViewMaterialisationSpec(views []ViewMaterialisationSpec) ([]interface{}, error) {
	res := make([]interface{}, 0, len(views))

	for _, view := range views {
		single := make(map[string]interface{})
		single["table"] = strconv.Itoa(view.Table)
		destinations, err := flattenDestinationReferences([]DestinationReference{view.Destination})
		if err != nil {
			return nil, err
		}
		single["destination"] = destinations
		res = append(res, single)
	}

	return res, nil
}
