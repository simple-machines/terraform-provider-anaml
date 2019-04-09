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
					Type: schema.TypeString,
				},
			},
			"event": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     eventSchema(),
			},
		},
	}
}

func eventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_column": {
				Type:     schema.TypeString,
				Required: true,
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

	if d.Get("expression").(string) != "" {
		table.Type = "view"
		table.Expression = d.Get("expression").(string)
		table.Sources = expandIdentifierList(d.Get("sources").([]interface{}))
	} else {
		table.Type = "root"
	}

	return &table
}

func expandEntityDescription(d *schema.ResourceData) *EventDescription {
	vIR := d.Get("event").(*schema.Set).List()
	ed := EventDescription{}

	if len(vIR) == 1 {
		r := vIR[0].(map[string]interface{})
		number, _ := strconv.Atoi(r["entity"].(string))

		ed = EventDescription{
			Id:     number,
			Column: r["key_column"].(string),
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

		oi["entity"] = strconv.Itoa(ed.Id)
		oi["key_column"] = ed.Column

		td := ed.TimestampInfo
		oi["timestamp_column"] = td.Column

		ois[0] = oi

		return ois
	}

	return make([]interface{}, 0)
}
