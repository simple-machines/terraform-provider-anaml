package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: resourceClusterRead,
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

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)

	cluster, err := c.FindCluster(name)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		return nil
	} else {
		d.SetId(strconv.Itoa(cluster.Id))
	}

	if err := d.Set("name", cluster.Name); err != nil {
		return err
	}
	if err := d.Set("description", cluster.Description); err != nil {
		return err
	}
	return err
}
