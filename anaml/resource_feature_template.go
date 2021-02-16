package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Event windows",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			"rows": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Event windows",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			"open": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "An event window",
				ExactlyOneOf: []string{"days", "rows", "open"},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !val.(bool) {
						errs = append(errs, errors.New("Open must be set to true"))
					}
					return
				},
			},
			"aggregations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"sum", "count", "countdistinct", "avg", "std", "last", "percentagechange", "absolutechange", "standardscore",
					}, true),
				},
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

	days := make([]int, 0, len(feature.Windows))
	rows := make([]int, 0, len(feature.Windows))
	aggs := make([]string, 0, len(feature.Aggregations))

	for _, window := range feature.Windows {
		if window.Type == "daywindow" {
			days = append(days, window.Days)
		}
		if window.Type == "rowwindow" {
			rows = append(rows, window.Rows)
		}
		if window.Type == "openwindow" {
			d.Set("open", true)
		}
	}

	for _, aggregation := range feature.Aggregations {
		aggs = append(aggs, aggregation.Type)
	}

	if err := d.Set("days", days); err != nil {
		return err
	}
	if err := d.Set("rows", rows); err != nil {
		return err
	}
	if err := d.Set("aggregations", aggs); err != nil {
		return err
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

	d.SetId(strconv.Itoa(e.Id))
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
		DataType: DataType{
			Type: d.Get("data_type").(string),
		},
		Select: SQLExpression{
			SQL: d.Get("select").(string),
		},
	}

	if d.Get("filter").(string) != "" {
		template.Filter = &SQLExpression{
			SQL: d.Get("filter").(string),
		}
	}

	if d.Get("table").(string) != "" {
		number, err := strconv.Atoi(d.Get("table").(string))

		if err != nil {
			return nil, err
		}

		hasOpen := d.Get("open").(bool)
		days := d.Get("days").(*schema.Set).List()
		rows := d.Get("rows").(*schema.Set).List()
		aggs := d.Get("aggregations").([]interface{})

		windows := make([]EventWindow, 0, len(rows)+len(days)+1)
		aggregations := make([]AggregateExpression, 0, 1)

		for _, day := range days {
			val, ok := day.(int)
			if ok && val != 0 {
				window := EventWindow{
					Type: "daywindow",
					Days: val,
				}
				windows = append(windows, window)
			}
		}

		for _, row := range rows {
			val, ok := row.(int)
			if ok && val != 0 {
				window := EventWindow{
					Type: "rowwindow",
					Rows: val,
				}
				windows = append(windows, window)
			}
		}

		if hasOpen {
			window := EventWindow{
				Type: "openwindow",
			}
			windows = append(windows, window)
		}

		for _, agg := range aggs {
			val, ok := agg.(string)
			if ok && val != "" {
				aggregation := AggregateExpression{
					Type: val,
				}
				aggregations = append(aggregations, aggregation)
			}
		}

		template.Type = "event"
		template.Table = number
		template.Windows = windows
		template.Aggregations = aggregations
	}

	return &template, nil
}
