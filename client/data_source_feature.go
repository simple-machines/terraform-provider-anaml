package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFeature() *schema.Resource {
	return &schema.Resource{
		Description: "A single Feature",

		Read: dataSourceFeatureRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The Feature's name",
				Required:    true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFeatureRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)

	feature, err := c.FindFeatureByName(name)
	if err != nil {
		return err
	}

	if feature == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(feature.ID))

	if err := d.Set("name", feature.Name); err != nil {
		return err
	}
	if err := d.Set("description", feature.Description); err != nil {
		return err
	}
	return nil
}
