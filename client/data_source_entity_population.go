package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEntityPopulation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEntityPopulationRead,
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

func dataSourceEntityPopulationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)

	population, err := c.FindEntityPopulationByName(name)
	if err != nil {
		return err
	}

	if population == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(population.ID))

	if err := d.Set("name", population.Name); err != nil {
		return err
	}
	if err := d.Set("description", population.Description); err != nil {
		return err
	}
	return nil
}
