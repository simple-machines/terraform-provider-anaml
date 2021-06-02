package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTableCaching() *schema.Resource {
	return &schema.Resource{
		Create: resourceTableCachingCreate,
		Read:   resourceTableCachingRead,
		Update: resourceTableCachingUpdate,
		Delete: resourceTableCachingDelete,
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
			"spec": {
				Type:        schema.TypeList,
				Description: "Table and entity specifications to cache with this job",
				Optional:    true,
				Elem:        specSchema(),
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
		},
	}
}

func specSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:         schema.TypeInt,
				Required:     true,
			},
			"entity": {
				Type:         schema.TypeInt,
				Required:     true,
			},
		},
	}
}

func resourceTableCachingRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableCachingID := d.Id()

	TableCaching, err := c.GetTableCaching(TableCachingID)
	if err != nil {
		return err
	}
	if TableCaching == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", TableCaching.Name); err != nil {
		return err
	}
	if err := d.Set("description", TableCaching.Description); err != nil {
		return err
	}
	if err := d.Set("spec", flattenTableCachingSpecs(TableCaching.Specs)); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(TableCaching.Cluster)); err != nil {
		return err
	}

	if TableCaching.Schedule.Type == "daily" {
		dailySchedules, err := parseDailySchedule(TableCaching.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("daily_schedule", dailySchedules); err != nil {
			return err
		}
	}

	if TableCaching.Schedule.Type == "cron" {
		cronSchedules, err := parseCronSchedule(TableCaching.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("cron_schedule", cronSchedules); err != nil {
			return err
		}
	}

	return err
}

func resourceTableCachingCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableCaching, err := composeTableCaching(d)
	if err != nil {
		return err
	}

	e, err := c.CreateTableCaching(*TableCaching)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceTableCachingUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableCachingID := d.Id()
	TableCaching, err := composeTableCaching(d)
	if err != nil {
		return err
	}

	err = c.UpdateTableCaching(TableCachingID, *TableCaching)
	if err != nil {
		return err
	}

	return nil
}

func composeTableCaching(d *schema.ResourceData) (*TableCaching, error) {
	cluster, err := strconv.Atoi(d.Get("cluster").(string))
	if err != nil {
		return nil, err
	}

	if dailySchedule, _ := expandSingleMap(d.Get("daily_schedule")); dailySchedule != nil {
		schedule, err := composeDailySchedule(dailySchedule)
		if err != nil {
			return nil, err
		}
		return &TableCaching{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
      Specs:       expandTableCachingSpecs(d),
			Cluster:     cluster,
			Schedule:    schedule,
		}, nil
	}

	if cronSchedule, _ := expandSingleMap(d.Get("cron_schedule")); cronSchedule != nil {
		schedule, err := composeCronSchedule(cronSchedule)
		if err != nil {
			return nil, err
		}
		return &TableCaching{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
      Specs:       expandTableCachingSpecs(d),
			Cluster:     cluster,
			Schedule:    schedule,
		}, nil
	}

	schedule := composeNeverSchedule()

	return &TableCaching{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
    Specs:       expandTableCachingSpecs(d),
		Cluster:     cluster,
		Schedule:    schedule,
	}, nil
}

func resourceTableCachingDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	TableCachingID := d.Id()

	err := c.DeleteTableCaching(TableCachingID)
	if err != nil {
		return err
	}

	return nil
}

func expandTableCachingSpecs(d *schema.ResourceData) []TableCachingSpec {
	drs := d.Get("spec").([]interface{})
	res := make([]TableCachingSpec, 0, len(drs))

	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})
		table, _  := val["table"].(int)
		entity, _ := val["entity"].(int)

		parsed := TableCachingSpec{
			Table:  table,
      Entity: entity,
		}
		res = append(res, parsed)
	}

	return res
}

func flattenTableCachingSpecs(specs []TableCachingSpec) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(specs))

	for _, spec := range specs {
		single := make(map[string]interface{})
		single["table"] = spec.Table
		single["entity"] = spec.Entity
		res = append(res, single)
	}

	return res
}
