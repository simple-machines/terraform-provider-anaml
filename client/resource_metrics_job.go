package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const metricsJobDescription = `
# Metrics Jobs (Schedule)

A Metrics Job schedules the run of a Metrics Set and describes its
output locations.
`

func ResourceMetricsJob() *schema.Resource {
	return &schema.Resource{
		Description: metricsJobDescription,
		Create:      resourceMetricsJobCreate,
		Read:        resourceMetricsJobRead,
		Update:      resourceMetricsJobUpdate,
		Delete:      resourceMetricsJobDelete,
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
			"metrics_set": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"principal": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"commit_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Commit to run metrics set on.",
			},
			"branch_target": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Branch to run metrics set on.",
				ConflictsWith: []string{"commit_target"},
			},
		},
	}
}

func resourceMetricsJobRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsJobID := d.Id()

	MetricsJob, err := c.GetMetricsJob(MetricsJobID)
	if err != nil {
		return err
	}
	if MetricsJob == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", MetricsJob.Name); err != nil {
		return err
	}
	if err := d.Set("description", MetricsJob.Description); err != nil {
		return err
	}

	if err := d.Set("principal", strconv.Itoa(MetricsJob.Principal)); err != nil {
		return err
	}

	destinations, err := flattenDestinationReferences(MetricsJob.Destinations)
	if err != nil {
		return err
	}

	if err := d.Set("metrics_set", strconv.Itoa(MetricsJob.MetricsSet)); err != nil {
		return err
	}
	if err := d.Set("enabled", MetricsJob.Enabled); err != nil {
		return err
	}
	if err := d.Set("destination", destinations); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(MetricsJob.Cluster)); err != nil {
		return err
	}
	if err := d.Set("cluster_property_sets", identifierList(MetricsJob.ClusterPropertySets)); err != nil {
		return err
	}
	if err := d.Set("labels", MetricsJob.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(MetricsJob.Attributes)); err != nil {
		return err
	}

	daily, cron, dependency, err := parseSchedule(MetricsJob.Schedule)
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

	if MetricsJob.VersionTarget != nil {
		if MetricsJob.VersionTarget.Commit != nil {
			if err := d.Set("commit_target", MetricsJob.VersionTarget.Commit); err != nil {
				return err
			}
			if err := d.Set("branch_target", nil); err != nil {
				return err
			}
		} else if MetricsJob.VersionTarget.Branch != nil {
			if err := d.Set("branch_target", MetricsJob.VersionTarget.Branch); err != nil {
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

func resourceMetricsJobCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsJob, err := composeMetricsJob(d)
	if err != nil {
		return err
	}

	e, err := c.CreateMetricsJob(*MetricsJob)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceMetricsJobUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsJobID := d.Id()
	MetricsJob, err := composeMetricsJob(d)
	if err != nil {
		return err
	}

	err = c.UpdateMetricsJob(MetricsJobID, *MetricsJob)
	if err != nil {
		return err
	}

	return nil
}

func composeMetricsJob(d *schema.ResourceData) (*MetricsJob, error) {
	metricsSet, err := getAnamlId(d, "metrics_set")
	if err != nil {
		return nil, err
	}
	cluster, err := getAnamlId(d, "cluster")
	if err != nil {
		return nil, err
	}
	principal, err := getAnamlId(d, "principal")
	if err != nil {
		return nil, err
	}
	schedule, err := composeSchedule(d)
	if err != nil {
		return nil, err
	}
	versionTarget := composeVersionTarget(d)
	destinations, err := expandDestinationReferences(d.Get("destination").([]interface{}))
	if err != nil {
		return nil, err
	}

	job := MetricsJob{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		MetricsSet:          metricsSet,
		Principal:           principal,
		Enabled:             d.Get("enabled").(bool),
		Destinations:        destinations,
		Cluster:             cluster,
		ClusterPropertySets: expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
		Schedule:            schedule,
		Labels:              expandLabels(d),
		Attributes:          expandAttributes(d),
		VersionTarget:       versionTarget,
	}

	return &job, nil
}

func resourceMetricsJobDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsJobID := d.Id()

	err := c.DeleteMetricsJob(MetricsJobID)
	if err != nil {
		return err
	}

	return nil
}
