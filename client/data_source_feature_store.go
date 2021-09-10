package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFeatureStore() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTableRead,
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
				Computed: true,
			},
		},
	}
}

func dataSourceFeatureStoreRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)

	feature, err := c.FindFeatureStoreByName(name)
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
