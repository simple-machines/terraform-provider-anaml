package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceFeature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFeatureRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAnamlIdentifier(),
				ConflictsWith: []string{"name"},
			},
			"days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "An event window",
				ConflictsWith: []string{"name", "rows"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"rows": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "An event window",
				ConflictsWith: []string{"name", "days"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
		},
	}
}

func dataSourceFeatureRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureName := d.Get("name").(string)
	templateStr := d.Get("template").(string)

	if featureName != "" {
		feature, err := c.FindFeatureByName(featureName)
		if err != nil {
			return err
		}
		if feature != nil {
			d.SetId(strconv.Itoa(feature.Id))
		}
	} else if templateStr != "" {
		template, _ := strconv.Atoi(templateStr)
		feature, err := c.FindFeatureByTemplate(template, d.Get("rows").(int), d.Get("days").(int))
		if err != nil {
			return err
		}
		if feature != nil {
			d.SetId(strconv.Itoa(feature.Id))
			if err := d.Set("name", feature.Name); err != nil {
				return err
			}
		}
	} else {
		return errors.New("Feature could not be found")
	}
	return nil
}
