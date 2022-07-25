package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"surname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validRoles(), false),
				},
			},
		},
	}
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	userID := d.Id()

	user, err := c.GetUser(userID)
	if err != nil {
		return err
	}
	if user == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", user.Name); err != nil {
		return err
	}
	if err := d.Set("email", user.Email); err != nil {
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

	return err
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	user := User{
		Name:      d.Get("name").(string),
		Email:     getNullableString(d, "email"),
		GivenName: getNullableString(d, "given_name"),
		Surname:   getNullableString(d, "surname"),
		Password:  getNullableString(d, "password"),
		Roles:     mapRolesToBackend(expandStringList(d.Get("roles").([]interface{}))),
	}

	e, err := c.CreateUser(user)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	userID := d.Id()
	user := User{
		Name:      d.Get("name").(string),
		Email:     getNullableString(d, "email"),
		GivenName: getNullableString(d, "given_name"),
		Surname:   getNullableString(d, "surname"),
		Roles:     mapRolesToBackend(expandStringList(d.Get("roles").([]interface{}))),
	}

	err := c.UpdateUser(userID, user)
	if err != nil {
		return err
	}

	if d.HasChange("password") {
		password := getNullableString(d, "password")
		err = c.UpdateUserPassword(userID, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	userID := d.Id()

	err := c.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
