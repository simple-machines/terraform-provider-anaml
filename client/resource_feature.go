package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const featureDescription = `# Features

A Feature is a time specific observation of the input data. Features define how data
is transformed into information that is useful for analytics or machine learning. At
a more concrete level, a Feature defines how data from tables is selected or aggregated
into a useful output. Each Feature selects data from a single source table but can use
one or more columns from that table. Each Feature is generated for a single entity.

There are two types of Features:
- Event Features
- Row Features


`

func ResourceFeature() *schema.Resource {
	return &schema.Resource{
		Description: featureDescription,
		Create:      resourceFeatureCreate,
		Read:        resourceFeatureRead,
		Update:      resourceFeatureUpdate,
		Delete:      resourceFeatureDelete,
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
			"table": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A reference to a Table ID the feature is derived from.",
				ValidateFunc: validateAnamlIdentifier(),
				RequiredWith: []string{"aggregation"},
			},
			"select": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "An SQL expression for the column to aggregate.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "An SQL column expression to filter with.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"hours": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The event window description for the number of days to aggregate over.",
				ConflictsWith: []string{"days", "rows", "months"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The event window description for the number of days to aggregate over.",
				ConflictsWith: []string{"hours", "rows", "months"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"months": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The event window description for the number of months to aggregate over.",
				ConflictsWith: []string{"hours", "days", "rows"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"rows": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The event window description for the number of rows (events) to aggregate over.",
				ConflictsWith: []string{"hours", "days", "months"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"aggregation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The aggregation to perform.",
				ValidateFunc: validation.StringInSlice([]string{
					"sum", "count", "countdistinct", "avg", "std", "min", "max", "minby", "maxby",
					"first", "last", "percentagechange", "absolutechange", "standardscore", "basketsum",
					"basketlast", "basketmax", "basketmin", "collectlist", "collectset",
				}, false),
				RequiredWith: []string{"table"},
			},
			"post_aggregation": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "An SQL expression to apply to the result of the feature aggregation.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"entity_restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of entity Id's that the feature is restricted to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"over": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "A list of Features this row feature depends on",
				AtLeastOneOf: []string{"table", "over"},
				RequiredWith: []string{"entity"},

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"entity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Entity to map a row feature over.",
				ValidateFunc: validateAnamlIdentifier(),
				RequiredWith: []string{"over"},
			},
			"template": {
				Type:         schema.TypeString,
				Description:  "The feature template this feature is derived from.",
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
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
		},
	}
}

func resourceFeatureRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureID := d.Id()

	feature, err := c.GetFeature(featureID)
	if err != nil {
		return err
	}
	if feature == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", feature.Name); err != nil {
		return err
	}
	if err := d.Set("description", feature.Description); err != nil {
		return err
	}
	if err := d.Set("select", feature.Select.SQL); err != nil {
		return err
	}
	if feature.Filter != nil {
		if err := d.Set("filter", feature.Filter.SQL); err != nil {
			return err
		}
	} else {
		d.Set("filter", nil)
	}

	if feature.PostAggExpr != nil {
		if err := d.Set("post_aggregation", feature.PostAggExpr.SQL); err != nil {
			return err
		}
	} else {
		d.Set("post_aggregation", nil)
	}

	if feature.TemplateID != nil {
		if err := d.Set("template", strconv.Itoa(*feature.TemplateID)); err != nil {
			return err
		}
	} else {
		d.Set("template", nil)
	}

	if feature.Type == "event" {
		if feature.Window.Type == "hourwindow" {
			if err := d.Set("hours", feature.Window.Hours); err != nil {
				return err
			}
		} else {
			if err = d.Set("hours", nil); err != nil {
				return err
			}
		}
		if feature.Window.Type == "daywindow" {
			if err := d.Set("days", feature.Window.Days); err != nil {
				return err
			}
		} else {
			if err = d.Set("days", nil); err != nil {
				return err
			}
		}
		if feature.Window.Type == "rowwindow" {
			if err := d.Set("rows", feature.Window.Rows); err != nil {
				return err
			}
		} else {
			if err := d.Set("rows", nil); err != nil {
				return err
			}
		}
		if feature.Window.Type == "monthwindow" {
			if err := d.Set("months", feature.Window.Months); err != nil {
				return err
			}
		} else {
			if err := d.Set("months", nil); err != nil {
				return err
			}
		}

		if err := d.Set("table", strconv.Itoa(feature.Table)); err != nil {
			return err
		}

		if err := d.Set("aggregation", feature.Aggregate.Type); err != nil {
			return err
		}

		if feature.EntityRestr != nil {
			if err := d.Set("entity_restrictions", identifierList(*feature.EntityRestr)); err != nil {
				return err
			}
		} else {
			d.Set("entity_restrictions", nil)
		}
	} else if feature.Type == "row" {
		if err := d.Set("over", identifierList(feature.Over)); err != nil {
			return err
		}
		if err := d.Set("entity", strconv.Itoa(feature.EntityID)); err != nil {
			return err
		}
	} else {
		return errors.New("Unrecognised ADT type for feature")
	}

	if err := d.Set("labels", feature.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(feature.Attributes)); err != nil {
		return err
	}
	return nil
}

func resourceFeatureCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	feature, err := buildFeature(d)
	if err != nil {
		return err
	}

	e, err := c.CreateFeature(*feature)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceFeatureUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureID := d.Id()
	table, err := buildFeature(d)
	if err != nil {
		return err
	}

	err = c.UpdateFeature(featureID, *table)
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureID := d.Id()

	err := c.DeleteFeature(featureID)
	if err != nil {
		return err
	}

	return nil
}

func buildFeature(d *schema.ResourceData) (*Feature, error) {
	template := getAnamlIdPointer(d, "template")
	feature := Feature{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Select: SQLExpression{
			SQL: d.Get("select").(string),
		},
		Aggregate: &AggregateExpression{
			Type: d.Get("aggregation").(string),
		},
		Labels:     expandLabels(d),
		Attributes: expandAttributes(d),
		TemplateID: template,
	}

	if filter, ok := d.GetOk("filter"); ok {
		feature.Filter = &SQLExpression{
			SQL: filter.(string),
		}
	}

	if post, ok := d.GetOk("post_aggregation"); ok {
		feature.PostAggExpr = &SQLExpression{
			SQL: post.(string),
		}
	}

	if table, ok := d.GetOk("table"); ok {
		number, err := strconv.Atoi(table.(string))

		if err != nil {
			return nil, err
		}

		window := EventWindow{}
		if hours, ok := d.GetOk("hours"); ok {
			window.Type = "hourwindow"
			window.Hours = hours.(int)
		} else if days, ok := d.GetOk("days"); ok {
			window.Type = "daywindow"
			window.Days = days.(int)
		} else if months, ok := d.GetOk("months"); ok {
			window.Type = "monthwindow"
			window.Months = months.(int)
		} else if rows, ok := d.GetOk("rows"); ok {
			window.Type = "rowwindow"
			window.Rows = rows.(int)
		} else {
			window.Type = "openwindow"
		}

		feature.Type = "event"
		feature.Table = number
		feature.Window = &window
		entity_restrictions := d.Get("entity_restrictions").([]interface{})
		if len(entity_restrictions) > 0 {
			listVal := expandIdentifierList(entity_restrictions)
			feature.EntityRestr = &listVal
		} else {
			feature.EntityRestr = nil
		}
	} else {
		number, err := getAnamlId(d, "entity")
		if err != nil {
			return nil, err
		}

		feature.EntityID = number
		feature.Type = "row"
		feature.Over = expandIdentifierList(d.Get("over").([]interface{}))
	}

	return &feature, nil
}
