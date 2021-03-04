package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Required: true,
			},
		},
	}
}

func dataSourceFeatureRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureName := d.Get("name").(string)

	feature, err := c.FindFeatureByName(featureName)
	if err != nil {
		return err
	}
	if feature != nil {
		d.SetId(strconv.Itoa(feature.ID))
	}
	return nil
}
