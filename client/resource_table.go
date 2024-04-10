package anaml

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const tableDescription = `# Tables

A Table represents a source of data for feature generation. A Table can be one of five types:

- External Table
- View Table
- Join Table
- Pivot Table
- Event Store Table

### External Tables

An External Table is the representation of a source table. Therefore you have to specify the
underlying data sources that the table is stored in.

### View Tables

A View Table is a pre-defined query over one or more other Root Tables or View Tables.
They function in the same way as views in relational databases. They can be used to transform
or join tables using arbitary SQL.

**When to use View Tables and when to use Features?**

In general, features should be used wherever possible to transform or aggregate columns as required.
This allows for better documentation, re-use and collaboration of features as well as allowing
for better optimise generation runs.

View Tables are mainly useful for joining data on keys other than entity id's such as reference data lookups.


### Join Tables

Join tables are similar to View tables, in that they perform operations other tables, but these perform
time aware joins between events and dimensional tables, or dimensional and dimensional.

Interestingly, Join tables perform efficient and correct joining of SCD2 tables with correctness properties
which are challenging to achieve with SQL.


### Pivot Tables

Pivot Tables allow features to be re-keyed to different entities.

This can be very useful, as often in business applications one may have different levels of entity; for example,
plans and customers. Data may be organised and keyed by plan, and some attributes and campaigns may target plans;
but also, plans are owned by customers, and when doing analysis on customers, knowing information about their
plans is crucial.

Pivot tables help to solve this issue by allowing features which are written for plans to be repurposed at the
customer level.

When construcing a pivot table, one needs to specify an entity mapping, which shows how to go from plans to
customers (this uses a feature query which outputs the customer for a particular plan), and the plan level
features which are to be aggregated to the customer.

Usually, the features one writes on a pivot table are simple aggregations, such as number of plans (count) or
average of some column with some filtering. Day and Row windows are not required.


### Event Store Tables

Event store tables are a robust store of tables backed and managed by Anaml. Usually, these will ingest data
from a Kafka topic, and describe mappings to events.

Table definitions for event store tables reference a managed event store, and the entity for which the data
should be interpreted.


## Time and Entity Descriptions

To be used in feature generation a Table must have one or more Entities as well as the semantics
for how to interpret timestamp columns for the table.

To achieve this, one should use one of the 'event', 'scd2', or 'point_in_time' blocks (as described below).

All of these blocks are able to accept a map of entities can be used as keys for this table in feature
generation, as well as their timestamp columns.
`

// ResourceTable ...
func ResourceTable() *schema.Resource {
	return &schema.Resource{
		Description: tableDescription,
		Create:      resourceTableCreate,
		Read:        resourceTableRead,
		Update:      resourceTableUpdate,
		Delete:      resourceTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the table in Anaml.",
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:         schema.TypeList,
				Description:  "Source information for a Root table.",
				Optional:     true,
				MaxItems:     1,
				Elem:         sourceSchema(),
				ExactlyOneOf: []string{"source", "event_store", "view", "pivot", "join"},
			},

			"view": {
				Type:        schema.TypeList,
				Description: "Define a View table using sources and an expression",
				Optional:    true,
				MaxItems:    1,
				Elem:        viewTableSchema(),
			},

			"join": {
				Type:          schema.TypeList,
				Description:   "Create a Join table, which performs time aware joins between tables.",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"event", "scd2", "point_in_time"},
				Elem:          joinTableSchema(),
			},

			"pivot": {
				Type:          schema.TypeList,
				Description:   "Create a Pivot table, which allows features to be aggregated for different entities.",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"event", "scd2", "point_in_time"},
				Elem:          pivotTableSchema(),
			},

			"event_store": {
				Type:        schema.TypeList,
				Description: "Information for how to interpret an Event store topic as a Table",
				Optional:    true,
				MaxItems:    1,
				Elem:        eventStoreSchema(),
			},

			"event": {
				Type:        schema.TypeList,
				Description: "This table contains events which occurred at a particular time",
				Optional:    true,
				MaxItems:    1,
				Elem:        eventSchema(),
			},

			"scd2": {
				Type:          schema.TypeList,
				Description:   "This table is a Slowly Changing Dimensional table",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"event", "point_in_time"},
				Elem:          scd2Schema(),
			},

			"point_in_time": {
				Type:          schema.TypeList,
				Description:   "This table is a Dimensional table with updates at particular times",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"event", "scd2"},
				Elem:          pointInTimeSchema(),
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

			"domain_modelling": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				ConflictsWith:    []string{"pivot"},
				Description:      "Model dimensions and measures for tables, and add virtual columns as simple SQL expressions",
				Elem:             domainModellingSchema(),
				DiffSuppressFunc: domainModellingSuppressFunc(),
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"timezone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func scd2Schema() *schema.Resource {
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
			"primary_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"from_column": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"valid_to_column": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"timezone": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
		},
	}
}

func pointInTimeSchema() *schema.Resource {
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
			"primary_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"timestamp_column": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"timezone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func joinTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:         schema.TypeString,
				Description:  "The root tables on the Left of the Join.",
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"joins": {
				Type:        schema.TypeList,
				Description: "Dimensions tables to join to.",
				Optional:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func pivotTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entity_mapping": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},

			"features": {
				Type:        schema.TypeList,
				Description: "Features to be included in this pivot table",
				Required:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func domainModellingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base": {
				Type:        schema.TypeList,
				Description: "An existing column to annotate",
				Optional:    true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Description:  "Name of the Table",
							Required:     true,
							ValidateFunc: validateAnamlName(),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dimension": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{},
						},
						"measure": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:        schema.TypeString,
										Description: "Units for the measure",
										Optional:    true,
									},
								},
							},
						},
						"not_null": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{},
						},
					},
				},
			},
			"virtual": {
				Type:        schema.TypeList,
				Description: "Dimensions tables to join to.",
				Optional:    true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Description:  "Name of the Table",
							Required:     true,
							ValidateFunc: validateAnamlName(),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expression": {
							Type:         schema.TypeString,
							Description:  "Name of the Table",
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"dimension": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{},
						},
						"measure": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:        schema.TypeString,
										Description: "Units for the measure",
										Optional:    true,
									},
								},
							},
						},
						"not_null": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{},
						},
					},
				},
			},
		},
	}
}

// We don't want to emit a diff if there is an empty
// domain modelling block. It's just for the prettiness
// of the terraform module.
func domainModellingSuppressFunc() schema.SchemaDiffSuppressFunc {
	sizeOfSlice := func(item interface{}) int {
		if slice, ok := item.([]interface{}); ok {
			return len(slice)
		}
		return 0
	}

	domainModellingSize := func(model interface{}) (int, error) {
		modelSlice, ok := model.([]interface{})
		if !ok {
			return 0, fmt.Errorf("expected []interface{}, got %T", model)
		}
		size := 0
		for _, model := range modelSlice {
			if mapped, ok := model.(map[string]interface{}); ok {
				size += sizeOfSlice(mapped["base"])
				size += sizeOfSlice(mapped["virtual"])
			}
		}
		return size, nil
	}

	return func(k, old, new string, d *schema.ResourceData) bool {
		oldModel, newModel := d.GetChange("domain_modelling")
		oldSize, err := domainModellingSize(oldModel)
		if err != nil {
			return false
		}
		newSize, err := domainModellingSize(newModel)
		if err != nil {
			return false
		}
		return oldSize == 0 && newSize == 0
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
				ExactlyOneOf: []string{"source.0.folder", "source.0.table_name", "source.0.topic"},
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

func eventStoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"store": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"entity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
		},
	}
}

func viewTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"expression": {
				Type:        schema.TypeString,
				Description: "Expression for a View table.",
				Required:    true,
			},
			"sources": {
				Type:        schema.TypeList,
				Description: "Tables upon which this View is created",
				Required:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
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

	flattenEntityDescription(d, table.EventInfo)

	if table.Type == "root" {
		if err := d.Set("source", flattenSourceReferences(table.Source)); err != nil {
			return err
		}
	} else {
		if err := d.Set("source", nil); err != nil {
			return err
		}
	}

	if table.Type == "view" {
		if err := d.Set("view", flattenViewReferences(table)); err != nil {
			return err
		}
	} else {
		if err := d.Set("view", nil); err != nil {
			return err
		}
	}

	if table.Type == "join" {
		if err := d.Set("join", flattenJoinTableSpecification(table)); err != nil {
			return err
		}
	} else {
		if err := d.Set("join", nil); err != nil {
			return err
		}
	}

	if table.Type == "pivot" {
		if err := d.Set("pivot", flattenPivotReferences(table)); err != nil {
			return err
		}
	} else {
		if err := d.Set("pivot", nil); err != nil {
			return err
		}
	}

	if table.Type == "eventstore" {
		if err := d.Set("event_store", flattenEventStoreReferences(table.Source)); err != nil {
			return err
		}
	} else {
		if err := d.Set("event_store", nil); err != nil {
			return err
		}
	}

	if err := d.Set("domain_modelling", flattenColumnInfo(table.Columns)); err != nil {
		return err
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
	table, err := buildTable(d)
	if err != nil {
		return err
	}
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
	table, err := buildTable(d)
	if err != nil {
		return err
	}

	err = c.UpdateTable(tableID, *table)
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

func buildTable(d *schema.ResourceData) (*Table, error) {
	table := Table{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		EventInfo:   expandEntityDescription(d),
		Labels:      expandLabels(d),
		Attributes:  expandAttributes(d),
	}

	if source, _ := expandSingleMap(d.Get("source")); source != nil {
		table.Type = "root"
		table.Source = expandSourceReferences(source)
	} else if view, _ := expandSingleMap(d.Get("view")); view != nil {
		table.Type = "view"
		expression, sources := expandViewSpecification(view)
		table.Expression = expression
		table.Sources = sources
	} else if join, _ := expandSingleMap(d.Get("join")); join != nil {
		table.Type = "join"
		base, others := expandJoinSpecification(join)
		table.Base = base
		table.Joins = others
	} else if pivot, _ := expandSingleMap(d.Get("pivot")); pivot != nil {
		table.Type = "pivot"
		entity_mapping, extra_features := expandPivotSpecification(pivot)
		table.EntityMapping = entity_mapping
		table.ExtraFeatures = extra_features
	} else if event_store, _ := expandSingleMap(d.Get("event_store")); event_store != nil {
		table.Type = "eventstore"
		table.Source = expandEventStoreReferences(event_store)
	} else {
		return nil, fmt.Errorf("Table did not appear to be one of the expected variants")
	}

	columns, err := expandColumnInfo(d)
	if err != nil {
		return nil, err
	}
	table.Columns = columns

	return &table, nil
}

func expandEntityDescription(d *schema.ResourceData) *EventDescription {
	events := d.Get("event").([]interface{})
	if len(events) == 1 {
		r := events[0].(map[string]interface{})
		entities := make(map[string]string)

		for k, v := range r["entities"].(map[string]interface{}) {
			entities[k] = v.(string)
		}

		return &EventDescription{
			Entities: entities,
			TimestampInfo: &TimestampInfo{
				Type:   "event",
				Column: r["timestamp_column"].(string),
				Zone:   r["timezone"].(string),
			},
		}
	}

	scd2s := d.Get("scd2").([]interface{})
	if len(scd2s) == 1 {
		r := scd2s[0].(map[string]interface{})
		entities := make(map[string]string)

		for k, v := range r["entities"].(map[string]interface{}) {
			entities[k] = v.(string)
		}

		primary, _ := strconv.Atoi(r["primary_key"].(string))
		return &EventDescription{
			Entities: entities,
			TimestampInfo: &TimestampInfo{
				Type:       "scd2",
				PrimaryKey: &primary,
				From:       r["from_column"].(string),
				ValidTo:    r["valid_to_column"].(string),
				Zone:       r["timezone"].(string),
			},
		}
	}

	pits := d.Get("point_in_time").([]interface{})
	if len(pits) == 1 {
		r := pits[0].(map[string]interface{})
		entities := make(map[string]string)

		for k, v := range r["entities"].(map[string]interface{}) {
			entities[k] = v.(string)
		}

		primary, _ := strconv.Atoi(r["primary_key"].(string))
		return &EventDescription{
			Entities: entities,
			TimestampInfo: &TimestampInfo{
				Type:       "snapshot",
				PrimaryKey: &primary,
				Column:     r["timestamp_column"].(string),
				Zone:       r["timezone"].(string),
			},
		}
	}

	return nil
}

func flattenEntityDescription(d *schema.ResourceData, ed *EventDescription) error {
	if ed != nil {
		single := make(map[string]interface{})
		single["entities"] = ed.Entities

		td := ed.TimestampInfo
		if td.Type == "event" {
			single["timestamp_column"] = td.Column
			single["timezone"] = td.Zone

			if err := d.Set("event", []interface{}{single}); err != nil {
				return err
			}
			if err := d.Set("point_in_time", nil); err != nil {
				return err
			}
			if err := d.Set("scd2", nil); err != nil {
				return err
			}
		}

		if td.Type == "scd2" && td.PrimaryKey != nil {
			single["primary_key"] = strconv.Itoa(*td.PrimaryKey)
			single["from_column"] = td.From
			single["valid_to_column"] = td.ValidTo
			single["timezone"] = td.Zone

			if err := d.Set("scd2", []interface{}{single}); err != nil {
				return err
			}
			if err := d.Set("point_in_time", nil); err != nil {
				return err
			}
			if err := d.Set("event", nil); err != nil {
				return err
			}
		}

		if td.Type == "snapshot" && td.PrimaryKey != nil {
			single["primary_key"] = strconv.Itoa(*td.PrimaryKey)
			single["timestamp_column"] = td.Column
			single["timezone"] = td.Zone

			if err := d.Set("point_in_time", []interface{}{single}); err != nil {
				return err
			}
			if err := d.Set("scd2", nil); err != nil {
				return err
			}
			if err := d.Set("event", nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func expandViewSpecification(val map[string]interface{}) (string, []int) {
	expression := val["expression"].(string)
	sourcesList := expandIdentifierList(val["sources"].([]interface{}))
	return expression, sourcesList
}

func expandSourceReferences(val map[string]interface{}) *SourceReference {
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

func expandEventStoreReferences(val map[string]interface{}) *SourceReference {
	store, _ := strconv.Atoi(val["store"].(string))
	entity, _ := strconv.Atoi(val["entity"].(string))
	topic, _ := val["topic"].(string)

	parsed := SourceReference{
		EventStoreId: store,
		Entity:       entity,
		Topic:        topic,
	}
	return &parsed
}

func expandJoinSpecification(val map[string]interface{}) (*int, []int) {
	store, _ := strconv.Atoi(val["table"].(string))
	joinList := expandIdentifierList(val["joins"].([]interface{}))
	return &store, joinList
}

func expandPivotSpecification(val map[string]interface{}) (int, []int) {
	entity_mapping, _ := strconv.Atoi(val["entity_mapping"].(string))
	extra_features := expandIdentifierList(val["features"].([]interface{}))
	return entity_mapping, extra_features
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

func flattenViewReferences(table *Table) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if table == nil {
		return res
	}

	single := make(map[string]interface{})
	single["expression"] = table.Expression
	single["sources"] = identifierList(table.Sources)
	res = append(res, single)

	return res
}

func flattenPivotReferences(table *Table) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if table == nil {
		return res
	}

	single := make(map[string]interface{})
	single["entity_mapping"] = strconv.Itoa(table.EntityMapping)
	single["features"] = identifierList(table.ExtraFeatures)
	res = append(res, single)

	return res
}

func flattenEventStoreReferences(source *SourceReference) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if source == nil {
		return res
	}

	single := make(map[string]interface{})
	single["store"] = strconv.Itoa(source.EventStoreId)
	single["entity"] = strconv.Itoa(source.Entity)
	single["topic"] = source.Topic
	res = append(res, single)

	return res
}

func flattenJoinTableSpecification(table *Table) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)
	if table.Base == nil {
		return res
	}

	single := make(map[string]interface{})
	single["table"] = strconv.Itoa(*table.Base)
	single["joins"] = identifierList(table.Joins)
	res = append(res, single)

	return res
}

func expandColumnKind(info map[string]interface{}) *ColumnKind {
	dimensions := info["dimension"].([]interface{})
	measures := info["measure"].([]interface{})
	for _, _ = range dimensions {
		return &ColumnKind{
			Type: "dimension",
		}
	}
	for _, measureRaw := range measures {
		ret := ColumnKind{
			Type: "measure",
		}
		if measure, ok := measureRaw.(map[string]interface{}); ok {
			if fetched, ok := measure["units"].(string); ok && fetched != "" {
				ret.Units = &fetched
			}
		}
		return &ret
	}

	return nil
}

func expandColumnConstraints(info map[string]interface{}) []ColumnConstraint {
	res := make([]ColumnConstraint, 0, 1)
	notnulls := info["not_null"].([]interface{})
	for _, _ = range notnulls {
		res = append(res, ColumnConstraint{
			Type: "notnull",
		})
	}

	return res
}

func expandColumnInfo(d *schema.ResourceData) (map[string]ColumnInfo, error) {
	modelling := d.Get("domain_modelling").([]interface{})
	res := make(map[string]ColumnInfo)

	for _, domain := range modelling {
		val := domain.(map[string]interface{})
		bases := val["base"].([]interface{})
		virtuals := val["virtual"].([]interface{})

		for _, base := range bases {
			name, baseInfo, err := createColumnInfo(base, "base")
			if err != nil {
				return nil, err
			}
			res[name] = baseInfo
		}
		for _, virtual := range virtuals {
			name, virtualInfo, err := createColumnInfo(virtual, "virtual")
			if err != nil {
				return nil, err
			}
			res[name] = virtualInfo
		}
	}
	return res, nil
}

func createColumnInfo(column interface{}, columnType string) (string, ColumnInfo, error) {
	value, ok := column.(map[string]interface{})
	if !ok {
		return "", ColumnInfo{}, fmt.Errorf("Expected Column info, couldn't derive")
	}
	name := value["name"].(string)
	description := value["description"].(string)
	kind := expandColumnKind(value)
	constraints := expandColumnConstraints(value)

	columnRepresentation := ColumnRepresentation{
		Type: columnType,
	}
	if columnType == "virtual" {
		if expr, ok := value["expression"].(string); ok {
			columnRepresentation.Expression = &expr
		} else {
			return "", ColumnInfo{}, fmt.Errorf("Expected expression for virtual column")
		}
	}

	return name, ColumnInfo{
		Description: description,
		Column:      &columnRepresentation,
		Kind:        kind,
		Constraints: constraints,
	}, nil
}

func flattenColumnKind(kind *ColumnKind) ([]map[string]interface{}, []map[string]interface{}) {
	if kind != nil {
		if kind.Type == "dimension" {
			single := make(map[string]interface{})
			return []map[string]interface{}{single}, nil
		}
		if kind.Type == "measure" {
			single := make(map[string]interface{})
			if kind.Units != nil {
				single["units"] = *kind.Units
			}
			return nil, []map[string]interface{}{single}
		}
	}
	return nil, nil
}

func flattenColumnConstraints(constraints []ColumnConstraint) []map[string]interface{} {
	notnulls := make([]map[string]interface{}, 0, 1)
	for _, constraint := range constraints {
		if constraint.Type == "notnull" {
			notnulls = append(notnulls, make(map[string]interface{}))
		}
	}
	return notnulls
}

func flattenColumnInfo(infos map[string]ColumnInfo) interface{} {
	res := make([]map[string]interface{}, 0, 1)
	bases := make([]map[string]interface{}, 0, len(infos))
	virtuals := make([]map[string]interface{}, 0, len(infos))

	for k, info := range infos {
		single := make(map[string]interface{})
		dimensions, measures := flattenColumnKind(info.Kind)
		notnulls := flattenColumnConstraints(info.Constraints)
		single["dimension"] = dimensions
		single["measure"] = measures
		single["not_null"] = notnulls

		if info.Column.Type == "base" {
			single["name"] = k
			single["description"] = info.Description
			bases = append(bases, single)
		} else {
			single["name"] = k
			single["expression"] = info.Column.Expression
			single["description"] = info.Description
			virtuals = append(virtuals, single)
		}
	}

	single := make(map[string]interface{})
	single["base"] = bases
	single["virtual"] = virtuals
	res = append(res, single)
	return res
}
