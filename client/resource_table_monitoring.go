package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTableMonitoring() *schema.Resource {
	return &schema.Resource{
		Create: resourceTableMonitoringCreate,
		Read:   resourceTableMonitoringRead,
		Update: resourceTableMonitoringUpdate,
		Delete: resourceTableMonitoringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
				Type:        schema.TypeSet,
				Description: "Tables to monitor with this job",
				Required:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
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
		},
	}
}

func resourceTableMonitoringRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableMonitoringID := d.Id()

	TableMonitoring, err := c.GetTableMonitoring(TableMonitoringID)
	if err != nil {
		return err
	}
	if TableMonitoring == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", TableMonitoring.Name); err != nil {
		return err
	}
	if err := d.Set("description", TableMonitoring.Description); err != nil {
		return err
	}
	if err := d.Set("tables", identifierList(TableMonitoring.Tables)); err != nil {
		return err
	}
	if err := d.Set("enabled", TableMonitoring.Enabled); err != nil {
		return err
	}
	if TableMonitoring.Principal != nil {
		if err := d.Set("principal", strconv.Itoa(*TableMonitoring.Principal)); err != nil {
			return err
		}
	}
	if err := d.Set("cluster", strconv.Itoa(TableMonitoring.Cluster)); err != nil {
		return err
	}
	if err := d.Set("cluster_property_sets", identifierList(TableMonitoring.ClusterPropertySets)); err != nil {
		return err
	}
	if TableMonitoring.Schedule.Type == "daily" {
		dailySchedules, err := parseDailySchedule(TableMonitoring.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("daily_schedule", dailySchedules); err != nil {
			return err
		}
	}

	if TableMonitoring.Schedule.Type == "cron" {
		cronSchedules, err := parseCronSchedule(TableMonitoring.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("cron_schedule", cronSchedules); err != nil {
			return err
		}
	}

	return err
}

func resourceTableMonitoringCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableMonitoring, err := composeTableMonitoring(d)
	if err != nil {
		return err
	}

	e, err := c.CreateTableMonitoring(*TableMonitoring)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceTableMonitoringUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableMonitoringID := d.Id()
	TableMonitoring, err := composeTableMonitoring(d)
	if err != nil {
		return err
	}

	err = c.UpdateTableMonitoring(TableMonitoringID, *TableMonitoring)
	if err != nil {
		return err
	}

	return nil
}

func composeTableMonitoring(d *schema.ResourceData) (*TableMonitoring, error) {
	cluster, err := strconv.Atoi(d.Get("cluster").(string))
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

	if dailySchedule, _ := expandSingleMap(d.Get("daily_schedule")); dailySchedule != nil {
		schedule, err := composeDailySchedule(dailySchedule)
		if err != nil {
			return nil, err
		}
		return &TableMonitoring{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			Tables:              expandIdentifierList(d.Get("tables").(*schema.Set).List()),
			Principal:           principal,
			Enabled:             d.Get("enabled").(bool),
			Cluster:             cluster,
			ClusterPropertySets: expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
			Schedule:            schedule,
		}, nil
	}

	if cronSchedule, _ := expandSingleMap(d.Get("cron_schedule")); cronSchedule != nil {
		schedule, err := composeCronSchedule(cronSchedule)
		if err != nil {
			return nil, err
		}
		return &TableMonitoring{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			Tables:              expandIdentifierList(d.Get("tables").(*schema.Set).List()),
			Principal:           principal,
			Enabled:             d.Get("enabled").(bool),
			Cluster:             cluster,
			ClusterPropertySets: expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
			Schedule:            schedule,
		}, nil
	}

	schedule := composeNeverSchedule()

	return &TableMonitoring{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		Principal:           principal,
		Cluster:             cluster,
		ClusterPropertySets: expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
		Schedule:            schedule,
	}, nil
}

func resourceTableMonitoringDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableMonitoringID := d.Id()

	err := c.DeleteTableMonitoring(TableMonitoringID)
	if err != nil {
		return err
	}

	return nil
}
