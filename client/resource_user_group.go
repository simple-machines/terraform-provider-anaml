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
				Type:        schema.TypeList,
				Description: "Users to include in the user group",
				Optional:    true,
				Elem:        userGroupMemberSchema(),
			},
			"external_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func userGroupMemberSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(validGroupMemberSource(), false),
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
	if err := d.Set("roles", mapRolesToFrontend(UserGroup.Roles)); err != nil {
		return err
	}
	if err := d.Set("members", flattenUserGroupMembers(UserGroup.Members)); err != nil {
		return err
	}
	if err := d.Set("external_group_id", UserGroup.ExternalGroupID); err != nil {
		return err
	}
	return err
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	groupMembers, err := expandUserGroupMembers(d.Get("members").([]interface{}))
	if err != nil {
		return err
	}
	UserGroup := UserGroup{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Roles:           mapRolesToBackend(expandStringList(d.Get("roles").([]interface{}))),
		Members:         groupMembers,
		ExternalGroupID: getNullableString(d, "external_group_id"),
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

	groupMembers, err := expandUserGroupMembers(d.Get("members").([]interface{}))
	if err != nil {
		return err
	}
	UserGroup := UserGroup{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Roles:           mapRolesToBackend(expandStringList(d.Get("roles").([]interface{}))),
		Members:         groupMembers,
		ExternalGroupID: getNullableString(d, "external_group_id"),
	}

	err = c.UpdateUserGroup(UserGroupID, UserGroup)
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

func flattenUserGroupMembers(members []UserGroupMember) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(members))
	for _, member := range members {
		single := make(map[string]interface{})
		single["user_id"] = member.UserID
		single["source"] = member.Source.Type
		res = append(res, single)
	}
	return res
}

func expandUserGroupMembers(members []interface{}) ([]UserGroupMember, error) {
	res := make([]UserGroupMember, 0, len(members))

	for _, member := range members {
		val, _ := member.(map[string]interface{})
		var source UserGroupMemberSource
		if val["source"] == "anaml" {
			source = UserGroupMemberSource{
				Type: "anaml",
			}
		} else {
			source = UserGroupMemberSource{
				Type: "external",
			}
		}

		parsed := UserGroupMember{
			UserID: val["user_id"].(int),
			Source: source,
		}
		res = append(res, parsed)
	}

	return res, nil
}
