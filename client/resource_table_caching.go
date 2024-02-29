package anaml

import (
	"context"
	"errors"
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
			"include": {
				Type:        schema.TypeList,
				Description: "Table and entity specifications to cache with this job",
				Optional:    true,
				MaxItems:    1,
				Elem:        planSchema(),
			},
			"auto": {
				Type:         schema.TypeList,
				Description:  "Table and entity specifications to cache with this job",
				Optional:     true,
				MaxItems:     1,
				Elem:         excludeSchema(),
				ExactlyOneOf: []string{"include", "auto"},
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
			"principal": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
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

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceTableCachingV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceTableCachingUpgradeV0,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec": {
							Type:        schema.TypeList,
							Description: "Table and entity specifications to cache with this job",
							Optional:    true,
							Elem:        specSchema(),
						},
					},
				},
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

func resourceTableCachingUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	plan := make(map[string]interface{})
	plan["include"] = rawState["spec"]

	rawState["include"] = plan
	delete(rawState, "spec")

	return rawState, nil
}

func planSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Description: "Table and entity specifications to cache with this job",
				Required:    true,
				Elem:        specSchema(),
				MinItems:    1,
			},
		},
	}
}

func excludeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude": {
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

	if err := d.Set("include", nil); err != nil {
		return err
	}
	if err := d.Set("auto", nil); err != nil {
		return err
	}
	if TableCaching.Principal != nil {
		if err := d.Set("principal", strconv.Itoa(*TableCaching.Principal)); err != nil {
			return err
		}
	}
	loc, plan := flattenTableCachingPlan(TableCaching.Plan)
	if err := d.Set(loc, plan); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(TableCaching.Cluster)); err != nil {
		return err
	}
	if err := d.Set("cluster_property_sets", identifierList(TableCaching.ClusterPropertySets)); err != nil {
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

	daily, cron, dependency, err := parseSchedule(TableCaching.Schedule)
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
	cluster, err := getAnamlId(d, "cluster")
	if err != nil {
		return nil, err
	}
	schedule, err := composeSchedule(d)
	if err != nil {
		return nil, err
	}
	principal := getAnamlIdPointer(d, "principal")
	retainment := getStringPointer(d, "retainment")
	plan, err := expandTableCachingPlan(d)
	if err != nil {
		return nil, err
	}

	return &TableCaching{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		PrefixURI:           d.Get("prefix_url").(string),
		Principal:           principal,
		Plan:                plan,
		Retainement:         retainment,
		Cluster:             cluster,
		ClusterPropertySets: expandIdentifierList(d.Get("cluster_property_sets").([]interface{})),
		Schedule:            schedule,
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

func expandTableCachingPlan(d *schema.ResourceData) (*CachingPlan, error) {
	if inclusion, _ := expandSingleMap(d.Get("include")); inclusion != nil {
		plan := CachingPlan{
			Type:  "inclusion",
			Specs: expandTableCachingSpecs(inclusion["spec"].([]interface{})),
		}
		return &plan, nil
	}

	if auto, _ := expandSingleMap(d.Get("auto")); auto != nil {
		plan := CachingPlan{
			Type:     "auto",
			Excluded: expandTableCachingSpecs(auto["exclude"].([]interface{})),
		}
		return &plan, nil
	}

	return nil, errors.New("Invalid caching plan type")
}

func expandTableCachingSpecs(drs []interface{}) []TableCachingSpec {
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

func flattenTableCachingPlan(plan *CachingPlan) (string, []map[string]interface{}) {
	res := make([]map[string]interface{}, 0, 1)
	loc := ""

	if plan.Type == "inclusion" {
		single := make(map[string]interface{})
		single["spec"] = flattenTableCachingSpecs(plan.Specs)
		res = append(res, single)
		loc = "include"
	} else {
		single := make(map[string]interface{})
		single["exclude"] = flattenTableCachingSpecs(plan.Excluded)
		res = append(res, single)
		loc = "auto"
	}

	return loc, res
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
