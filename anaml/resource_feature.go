package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceFeature() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureCreate,
		Read:   resourceFeatureRead,
		Update: resourceFeatureUpdate,
		Delete: resourceFeatureDelete,
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
				Default:  "root",
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "string",
			},
			"table": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "A reference to a Table ID the feature is derived from",
				ValidateFunc:  validateAnamlIdentifier(),
				ConflictsWith: []string{"over"},
			},
			"select": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A reference to a Table ID the feature is derived from",
			},
			"days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "An event window",
				ExactlyOneOf: []string{"days", "rows", "open", "over"},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rows": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "An event window",
				ExactlyOneOf: []string{"days", "rows", "open", "over"},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"open": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "An event window",
				ExactlyOneOf: []string{"days", "rows", "open", "over"},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !val.(bool) {
						errs = append(errs, errors.New("Open must be set to true"))
					}
					return
				},
			},
			"aggregation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sum", "count", "countdistinct", "avg", "std", "last", "percentagechange", "absolutechange", "standardscore",
				}, true),
				ConflictsWith: []string{"over"},
			},
			"over": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "A list of Features this row feature depends on",
				ConflictsWith: []string{"table"},

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"entity": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAnamlIdentifier(),
				ConflictsWith: []string{"table"},
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

	if feature.Type == "event" {
		if feature.Window.Type == "daywindow" {
			err := d.Set("days", feature.Window.Days)
			if err != nil {
				return err
			}
		} else if feature.Window.Type == "rowwindow" {
			err := d.Set("rows", feature.Window.Rows)
			if err != nil {
				return err
			}
		} else if feature.Window.Type == "openwindow" {
			err := d.Set("open", true)
			if err != nil {
				return err
			}
		}

		if err := d.Set("table", strconv.Itoa(feature.Table)); err != nil {
			return err
		}

		if err := d.Set("aggregation", feature.Aggregate.Type); err != nil {
			return err
		}
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

	d.SetId(strconv.Itoa(e.Id))
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
	feature := Feature{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DataType: DataType{
			Type: d.Get("data_type").(string),
		},
		Select: SQLExpression{
			SQL: d.Get("select").(string),
		},
		Aggregate: &AggregateExpression{
			Type: d.Get("aggregation").(string),
		},
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
		} else if d.Get("open").(bool) {
			window.Type = "openwindow"
		} else {
			return nil, errors.New("Open window not set to true")
		}

		feature.Type = "event"
		feature.Table = number
		feature.Window = &window
	} else {
		feature.Type = "row"
		feature.Over = expandIdentifierList(d.Get("over").([]interface{}))
		number, _ := strconv.Atoi(d.Get("entity").(string))
		feature.EntityId = number
	}

	return &feature, nil
}
