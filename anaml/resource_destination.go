package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDestination() *schema.Resource {
	return &schema.Resource{
		Read: resourceDestinationRead,
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

func resourceDestinationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)

	destination, err := c.FindDestination(name)
	if err != nil {
		return err
	}

	if destination == nil {
		d.SetId("")
		return nil
	} else {
		d.SetId(strconv.Itoa(destination.Id))
	}

	if err := d.Set("name", destination.Name); err != nil {
		return err
	}
	if err := d.Set("description", destination.Description); err != nil {
		return err
	}
	return err
}
