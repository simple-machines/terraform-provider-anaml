package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceTable ...
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				Elem:         sourceSchema(),
				ExactlyOneOf: []string{"source", "expression", "entity_mapping"},
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

			"entity_mapping": {
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

func sourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"folder": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"table_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"topic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
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

	if table.Type == "root" {
		if err := d.Set("source", flattenSourceReferences(table.Source)); err != nil {
			return err
		}
	}

	if table.Type == "pivot" {
		if err := d.Set("entity_mapping", strconv.Itoa(table.EntityMapping)); err != nil {
			return err
		}

		if err := d.Set("extra_features", identifierList(table.ExtraFeatures)); err != nil {
			return err
		}
	}

	if err := d.Set("labels", table.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(table.Attributes)); err != nil {
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

	d.SetId(strconv.Itoa(e.ID))
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
		Labels:      expandStringList(d.Get("labels").([]interface{})),
		Attributes:  expandAttributes(d),
	}

	if d.Get("expression").(string) != "" {
		table.Type = "view"
		table.Expression = d.Get("expression").(string)
		table.Sources = expandIdentifierList(d.Get("sources").([]interface{}))
	} else if d.Get("entity_mapping").(string) != "" {
		table.Type = "pivot"
		entity, _ := strconv.Atoi(d.Get("entity_mapping").(string))

		table.EntityMapping = entity
		table.ExtraFeatures = expandIdentifierList(d.Get("extra_features").([]interface{}))
	} else {
		table.Type = "root"
		table.Source = expandSourceReferences(d)
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

func expandSourceReferences(d *schema.ResourceData) *SourceReference {
	srs := d.Get("source").([]interface{})

	for _, sr := range srs {
		val, _ := sr.(map[string]interface{})
		sourceID, _ := strconv.Atoi(val["source"].(string))

		source_type := ""
		if v, ok := val["folder"].(string); ok && v != "" {
			source_type = "folder"
		}
		if v, ok := val["table_name"].(string); ok && v != "" {
			source_type = "table"
		}
		if v, ok := val["topic"].(string); ok && v != "" {
			source_type = "topic"
		}

		parsed := SourceReference{
			Type:      source_type,
			SourceID:  sourceID,
			Folder:    val["folder"].(string),
			TableName: val["table_name"].(string),
			Topic:     val["topic"].(string),
		}
		return &parsed
	}

	return nil
}

func flattenSourceReferences(source *SourceReference) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if source == nil {
		return res
	}

	single := make(map[string]interface{})
	single["source"] = strconv.Itoa(source.SourceID)
	single["folder"] = source.Folder
	single["table_name"] = source.TableName
	single["topic"] = source.Topic
	res = append(res, single)

	return res
}
