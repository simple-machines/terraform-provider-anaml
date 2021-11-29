package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const apiKeyDescription = `
# AccessTokens

AccessTokens are how to authenticate with Anaml programmatically.
`

func ResourceAccessToken() *schema.Resource {
	return &schema.Resource{
		Description: webhooksDescription,
		Create:      resourceAccessTokenCreate,
		Read:        resourceAccessTokenRead,
		Delete:      resourceAccessTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
				ForceNew:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validRoles(), false),
				},
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAccessTokenRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tokenId := d.Id()

	ownerKey, _ := d.GetChange("owner")
	owner, _ := strconv.Atoi(ownerKey.(string))

	token, err := c.GetAccessToken(owner, tokenId)
	if err != nil {
		return err
	}
	if token == nil {
		d.SetId("")
		return nil
	}

	return err
}

func resourceAccessTokenCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	request := AccessToken{
		Description: d.Get("description").(string),
		Roles:       mapRolesToBackend(expandStringList(d.Get("roles").([]interface{}))),
	}
	owner, _ := strconv.Atoi(d.Get("owner").(string))

	token, err := c.CreateAccessToken(owner, request)
	if err != nil {
		return err
	}

	d.SetId(token.ID)
	d.Set("secret", token.Secret)
	return err
}

func resourceAccessTokenDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	tokenID := d.Id()
	ownerKey, _ := d.GetChange("owner")
	owner, _ := strconv.Atoi(ownerKey.(string))

	err := c.DeleteAccessToken(owner, tokenID)
	if err != nil {
		return err
	}

	return nil
}
