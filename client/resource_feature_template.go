package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceFeatureTemplate ...
func ResourceFeatureTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureTemplateCreate,
		Read:   resourceFeatureTemplateRead,
		Update: resourceFeatureTemplateUpdate,
		Delete: resourceFeatureTemplateDelete,
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
				Default:  "root",
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "string",
			},
			"table": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A reference to a Table ID the feature is derived from",
				ValidateFunc: validateAnamlIdentifier(),
			},
			"select": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An SQL expression for the column to aggregate",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An SQL column expression to filter with",
			},
			"days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "An event window",
				ConflictsWith: []string{"rows"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"rows": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "An event window",
				ConflictsWith: []string{"days"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"aggregation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sum", "count", "countdistinct", "avg", "std", "last", "percentagechange", "absolutechange",
					"standardscore", "basketsum", "basketlast", "collectlist", "collectset",
				}, true),
			},
			"over": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "A list of Features this row feature depends on",
				AtLeastOneOf: []string{"table", "over"},

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"entity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAnamlIdentifier(),
				RequiredWith: []string{"over"},
			},
		},
	}
}

func resourceFeatureTemplateRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureID := d.Id()

	feature, err := c.GetFeatureTemplate(featureID)
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
	}

	if feature.Type == "event" {
		if feature.Window.Type == "daywindow" {
			if err := d.Set("days", feature.Window.Days); err != nil {
				return err
			}
			if err = d.Set("rows", nil); err != nil {
				return err
			}
		} else if feature.Window.Type == "rowwindow" {
			if err := d.Set("rows", feature.Window.Rows); err != nil {
				return err
			}
			if err = d.Set("days", nil); err != nil {
				return err
			}
		} else if feature.Window.Type == "openwindow" {
			if err = d.Set("days", nil); err != nil {
				return err
			}
			if err = d.Set("rows", nil); err != nil {
				return err
			}
		}

		if err := d.Set("table", strconv.Itoa(feature.Table)); err != nil {
			return err
		}

		if err := d.Set("aggregation", feature.Aggregate.Type); err != nil {
			return err
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

	return nil
}

func resourceFeatureTemplateCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	template, err := buildFeatureTemplate(d)
	if err != nil {
		return err
	}

	e, err := c.CreateFeatureTemplate(*template)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceFeatureTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	templateID := d.Id()
	template, err := buildFeatureTemplate(d)
	if err != nil {
		return err
	}

	err = c.UpdateFeatureTemplate(templateID, *template)
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureTemplateDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	templateID := d.Id()

	err := c.DeleteFeatureTemplate(templateID)
	if err != nil {
		return err
	}

	return nil
}

func buildFeatureTemplate(d *schema.ResourceData) (*FeatureTemplate, error) {
	template := FeatureTemplate{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Select: SQLExpression{
			SQL: d.Get("select").(string),
		},
	}

	if d.Get("filter").(string) != "" {
		template.Filter = &SQLExpression{
			SQL: d.Get("filter").(string),
		}
	}

	if d.Get("aggregation").(string) != "" {
		template.Aggregate = &AggregateExpression{
			Type: d.Get("aggregation").(string),
		}
	}

	if d.Get("table").(string) != "" {
		number, err := strconv.Atoi(d.Get("table").(string))

		if err != nil {
			return nil, err
		}

		window := EventWindow{}
		if d.Get("days").(int) != 0 {
			window.Type = "daywindow"
			window.Days = d.Get("days").(int)
		} else if d.Get("rows").(int) != 0 {
			window.Type = "rowwindow"
			window.Rows = d.Get("rows").(int)
		} else {
			window.Type = "openwindow"
		}

		template.Type = "event"
		template.Table = number
		template.Window = &window
	} else {
		template.Type = "row"
		template.Over = expandIdentifierList(d.Get("over").([]interface{}))
		number, _ := strconv.Atoi(d.Get("entity").(string))
		template.EntityID = number
	}

	return &template, nil
}
