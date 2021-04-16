package anaml

import (
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
			"aggregation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sum", "count", "countdistinct", "avg", "std", "last", "percentagechange", "absolutechange", "standardscore", "basketsum", "basketlast",
				}, true),
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
	if feature.Aggregate != nil {
		if err := d.Set("aggregation", feature.Aggregate.Type); err != nil {
			return err
		}
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

		template.Type = "event"
		template.Table = number
	}

	return &template, nil
}
