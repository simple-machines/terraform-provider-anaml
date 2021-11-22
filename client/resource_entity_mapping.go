package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const entityMappingDescription = `# Entity Mappings

An Entity Mapping is a relationship between one or more Entities so that features can
automatically be aggregated at different levels without defining the Feature multiple times.

Each Entity Mapping specifies:

- From Entity - The Entity that the Feature is initially defined on
- To Entity - The Entity that can be converted to
- Mapping Feature - A Feature that is defined on the *from* Entity and has a value of the *to* Entity.

For example, given a customer and account Entities where each customer has multiple accounts
you could define an Entity Mapping from Account to Customer:

- From Entity = Account
- To Entity = Customer
- Mapping Feature = Last Customer Id for each Account Id
`

func ResourceEntityMapping() *schema.Resource {
	return &schema.Resource{
		Description: entityMappingDescription,
		Create:      resourceEntityMappingCreate,
		Read:        resourceEntityMappingRead,
		Update:      resourceEntityMappingUpdate,
		Delete:      resourceEntityMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"from": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"to": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"mapping": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
		},
	}
}

func resourceEntityMappingRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	mappingID := d.Id()

	mapping, err := c.GetEntityMapping(mappingID)
	if err != nil {
		return err
	}
	if mapping == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("from", strconv.Itoa(mapping.From)); err != nil {
		return err
	}
	if err := d.Set("to", strconv.Itoa(mapping.To)); err != nil {
		return err
	}
	if err := d.Set("mapping", strconv.Itoa(mapping.Mapping)); err != nil {
		return err
	}
	return err
}

func resourceEntityMappingCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	from, _ := strconv.Atoi(d.Get("from").(string))
	to, _ := strconv.Atoi(d.Get("to").(string))
	feat, _ := strconv.Atoi(d.Get("mapping").(string))
	mapping := EntityMapping{
		From:    from,
		To:      to,
		Mapping: feat,
	}

	e, err := c.CreateEntityMapping(mapping)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceEntityMappingUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	mappingID := d.Id()
	from, _ := strconv.Atoi(d.Get("from").(string))
	to, _ := strconv.Atoi(d.Get("to").(string))
	feat, _ := strconv.Atoi(d.Get("mapping").(string))

	mapping := EntityMapping{
		From:    from,
		To:      to,
		Mapping: feat,
	}
	err := c.UpdateEntityMapping(mappingID, mapping)
	if err != nil {
		return err
	}

	return nil
}

func resourceEntityMappingDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	mappingID := d.Id()

	err := c.DeleteEntityMapping(mappingID)
	if err != nil {
		return err
	}

	return nil
}
