package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTableCreate,
		Read:   resourceTableRead,
		Update: resourceTableUpdate,
		Delete: resourceTableDelete,
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
			"expression": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"sources"},
			},
			"sources": {
				Type:         schema.TypeList,
				Description:  "Tables upon which this view is created",
				Optional:     true,
				RequiredWith: []string{"expression"},

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"event": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     eventSchema(),
			},

			"entity": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAnamlIdentifier(),
				ConflictsWith: []string{"event"},
				RequiredWith:  []string{"pivot_feature"},
			},
			"pivot_feature": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAnamlIdentifier(),
				ConflictsWith: []string{"event"},
			},
			"extra_features": {
				Type:          schema.TypeList,
				Description:   "Tables upon which this view is created",
				Optional:      true,
				ConflictsWith: []string{"event"},

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func eventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateDiagFunc: validateMapKeysAnamlIdentifier(),
			},
			"timestamp_column": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTableRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tableID := d.Id()

	table, err := c.GetTable(tableID)
	if err != nil {
		return err
	}

	if table == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", table.Name); err != nil {
		return err
	}
	if err := d.Set("description", table.Description); err != nil {
		return err
	}

	eventItems := flattenEntityDescription(table.EventInfo)
	if err := d.Set("event", eventItems); err != nil {
		return err
	}

	if err := d.Set("expression", table.Expression); err != nil {
		return err
	}

	if err := d.Set("sources", identifierList(table.Sources)); err != nil {
		return err
	}

	if table.Type == "pivot" {
		if err := d.Set("entity", strconv.Itoa(table.Entity)); err != nil {
			return err
		}

		if err := d.Set("pivot_feature", strconv.Itoa(table.Pivot)); err != nil {
			return err
		}

		if err := d.Set("extra_features", identifierList(table.ExtraFeatures)); err != nil {
			return err
		}
	}

	return nil
}

func resourceTableCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	table := buildTable(d)
	e, err := c.CreateTable(*table)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.Id))
	return err
}

func resourceTableUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tableID := d.Id()
	table := buildTable(d)

	err := c.UpdateTable(tableID, *table)
	if err != nil {
		return err
	}

	return nil
}

func resourceTableDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tableID := d.Id()

	err := c.DeleteTable(tableID)
	if err != nil {
		return err
	}

	return nil
}

func buildTable(d *schema.ResourceData) *Table {
	table := Table{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		EventInfo:   expandEntityDescription(d),
	}

	entity, _ := strconv.Atoi(d.Get("entity").(string))
	pivot, _ := strconv.Atoi(d.Get("pivot_feature").(string))

	if d.Get("expression").(string) != "" {
		table.Type = "view"
		table.Expression = d.Get("expression").(string)
		table.Sources = expandIdentifierList(d.Get("sources").([]interface{}))
	} else if d.Get("pivot_feature").(string) != "" {
		table.Type = "pivot"
		table.Entity = entity
		table.Pivot = pivot
		table.ExtraFeatures = expandIdentifierList(d.Get("extra_features").([]interface{}))
	} else {
		table.Type = "root"
	}

	return &table
}

func expandEntityDescription(d *schema.ResourceData) *EventDescription {
	vIR := d.Get("event").([]interface{})
	ed := EventDescription{}

	if len(vIR) == 1 {
		r := vIR[0].(map[string]interface{})

		entities := make(map[string]string)

		for k, v := range r["entities"].(map[string]interface{}) {
			entities[k] = v.(string)
		}

		ed = EventDescription{
			Entities: entities,
			TimestampInfo: &TimestampInfo{
				Column: r["timestamp_column"].(string),
			},
		}
	} else {
		return nil
	}

	return &ed
}

func flattenEntityDescription(ed *EventDescription) []interface{} {
	if ed != nil {
		ois := make([]interface{}, 1, 1)

		oi := make(map[string]interface{})

		oi["entities"] = ed.Entities

		td := ed.TimestampInfo
		oi["timestamp_column"] = td.Column

		ois[0] = oi

		return ois
	}

	return make([]interface{}, 0)
}
