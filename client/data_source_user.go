package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"surname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	email := d.Get("email").(string)

	user, err := c.FindUserByEmail(email)
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(user.ID))

	if err := d.Set("name", user.Name); err != nil {
		return err
	}
	if err := d.Set("given_name", user.GivenName); err != nil {
		return err
	}
	if err := d.Set("surname", user.Surname); err != nil {
		return err
	}
	if err := d.Set("roles", mapRolesToFrontend(user.Roles)); err != nil {
		return err
	}
	return nil
}
