package anaml

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"prefix_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:        schema.TypeList,
				Description: "Table and entity specifications to cache with this job",
				Required:    true,
				MaxItems:    1,

				Elem: planSchema(),
			},
			"retainment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
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

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceTableCachingV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceExampleInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceTableCachingV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix_url": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceExampleInstanceStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	plan := make(map[string]interface{})
	plan["spec"] = rawState["spec"]

	rawState["plan"] = plan
	delete(rawState, "spec")

	return rawState, nil
}

func planSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Description: "Table and entity specifications to cache with this job",
				Optional:    true,
				Elem:        specSchema(),
			},
		},
	}
}

func specSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"entity": {
				Type:     schema.TypeInt,
				Required: true,
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
	if err := d.Set("plan", flattenTableCachingPlan(TableCaching.Plan)); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(TableCaching.Cluster)); err != nil {
		return err
	}
	if err := d.Set("prefix_url", TableCaching.PrefixURI); err != nil {
		return err
	}

	if TableCaching.Retainement != nil {
		if err := d.Set("retainment", TableCaching.Retainement); err != nil {
			return err
		}
	} else {
		d.Set("retainment", nil)
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

	var retainment *string
	if d.Get("retainment").(string) != "" {
		retainmentstr := d.Get("retainment").(string)
		retainment = &retainmentstr
	}

	return &TableCaching{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PrefixURI:   d.Get("prefix_url").(string),
		Plan:        expandTableCachingPlan(d),
		Retainement: retainment,
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

func expandTableCachingPlan(d *schema.ResourceData) *CachingPlan {
	drs, _ := expandSingleMap(d.Get("plan"))
	res := CachingPlan{
		Type: "inclusion",
		Specs: expandTableCachingSpecs(drs),
	}

	return &res
}

func expandTableCachingSpecs(d map[string]interface{}) []TableCachingSpec {
	drs := d["spec"].([]interface{})
	res := make([]TableCachingSpec, 0, len(drs))

	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})
		table, _ := val["table"].(int)
		entity, _ := val["entity"].(int)

		parsed := TableCachingSpec{
			Table:  table,
			Entity: entity,
		}
		res = append(res, parsed)
	}

	return res
}

func flattenTableCachingPlan(plan *CachingPlan) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	single := make(map[string]interface{})
	single["spec"] = flattenTableCachingSpecs(plan.Specs)
	res = append(res, single)

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
