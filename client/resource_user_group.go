package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,
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
			"members": {
				Type:        schema.TypeSet,
				Description: "Users to include in the user group",
				Required:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	UserGroupID := d.Id()

	UserGroup, err := c.GetUserGroup(UserGroupID)
	if err != nil {
		return err
	}
	if UserGroup == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", UserGroup.Name); err != nil {
		return err
	}
	if err := d.Set("description", UserGroup.Description); err != nil {
		return err
	}
	if err := d.Set("roles", UserGroup.Roles); err != nil {
		return err
	}
	if err := d.Set("members", identifierList(UserGroup.Members)); err != nil {
		return err
	}
	return err
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	UserGroup := UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Roles:       expandStringList(d.Get("roles").([]interface{})),
		Members:     expandIdentifierList(d.Get("members").(*schema.Set).List()),
	}

	ug, err := c.CreateUserGroup(UserGroup)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ug.ID))
	return err
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	UserGroupID := d.Id()

	UserGroup := UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Roles:       expandStringList(d.Get("roles").([]interface{})),
		Members:     expandIdentifierList(d.Get("members").(*schema.Set).List()),
	}

	err := c.UpdateUserGroup(UserGroupID, UserGroup)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	UserGroupID := d.Id()

	err := c.DeleteUserGroup(UserGroupID)
	if err != nil {
		return err
	}

	return nil
}
