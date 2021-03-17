package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSourceRead,
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

func dataSourceSourceRead(d *schema.ResourceData, m interface{}) error {
			 c := m.(*Client)
			 name := d.Get("name").(string)

			 source, err := c.FindSource(name)
			 if err != nil {
							 return err
			 }

			 if source == nil {
							 d.SetId("")
							 return nil
			 } else {
							 d.SetId(strconv.Itoa(source.ID))
			 }

			 if err := d.Set("name", source.Name); err != nil {
							 return err
			 }
			 if err := d.Set("description", source.Description); err != nil {
							 return err
			 }
			 return err
}
