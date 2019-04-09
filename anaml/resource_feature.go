package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFeature() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureCreate,
		Read:   resourceFeatureRead,
		Update: resourceFeatureUpdate,
		Delete: resourceFeatureDelete,

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "A reference to a Table ID the feature is derived from",
			},
			"select": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A reference to a Table ID the feature is derived from",
			},
			"window": {
				Type:        schema.TypeSet,
				Elem:        windowSchema(),
				Required:    true,
				Description: "An event window",
				MinItems:    1,
				MaxItems:    1,
			},
			"aggregation": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func windowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"days": {
				Type:     schema.TypeInt,
				Optional: true,
				// ExactlyOneOf: []string{"rows", "days", "open"},
			},
			"rows": {
				Type:     schema.TypeInt,
				Optional: true,
				// ExactlyOneOf: []string{"rows", "days", "open"},
			},
			"open": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
				// ExactlyOneOf: []string{"rows", "days", "open"},
			},
		},
	}
}

func resourceFeatureRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tableID := d.Id()

	feature, err := c.GetFeature(tableID)
	if err != nil {
		return err
	}
	if feature == nil {
		d.SetId("")
		return nil
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
		Aggregate: AggregateExpression{
			Type: d.Get("aggregation").(string),
		},
	}

	if d.Get("table").(string) != "" {
		number, err := strconv.Atoi(d.Get("table").(string))

		if err != nil {
			return nil, err
		}

		window, err := expandWindowList(d)
		if err != nil {
			return nil, err
		}

		feature.Type = "event"
		feature.Table = number
		feature.Window = *window
	} else {
		feature.Type = "row"
		return nil, errors.New("Rows not quite implemented yet")
	}

	return &feature, nil
}

func expandWindowList(d *schema.ResourceData) (*EventWindow, error) {
	vIR := d.Get("window").(*schema.Set).List()
	ew := EventWindow{}

	if len(vIR) == 1 {
		r := vIR[0].(map[string]interface{})

		if r["days"].(int) != 0 {
			days, _ := r["days"].(int)
			ew = EventWindow{
				Type: "daywindow",
				Days: days,
			}
		} else if r["rows"].(int) != 0 {
			rows, _ := r["rows"].(int)
			ew = EventWindow{
				Type: "rowwindow",
				Rows: rows,
			}
		} else if r["open"] != nil {
			ew = EventWindow{
				Type: "openwindow",
			}
		} else {
			return nil, errors.New("No rows, days or open specified")
		}

	} else {
		return nil, errors.New("More than one specified")
	}

	return &ew, nil
}
