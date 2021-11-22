package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const entityDescription = `# Entities

An Entity is an item in the business domain. Common examples of Entities are:

- Customers
- Accounts
- Products
- Orders

Anything that has a unique identifier and would be useful to report or predict on could be an Entity.

In a relational database, the identifiers for Entities will often be used for primary keys.

Tables need to specify one or more columns with entity identifiers in order to be used for Feature definitions.

Features will be generated for a specific Entity. This means the aggregation will be grouped by each Entity identitifer.
`

func ResourceEntity() *schema.Resource {
	return &schema.Resource{
		Description: entityDescription,
		Create:      resourceEntityCreate,
		Read:        resourceEntityRead,
		Update:      resourceEntityUpdate,
		Delete:      resourceEntityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_column": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Labels to attach to the object",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"attribute": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
		},
	}
}

func resourceEntityRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entityID := d.Id()

	entity, err := c.GetEntity(entityID)
	if err != nil {
		return err
	}
	if entity == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", entity.Name); err != nil {
		return err
	}
	if err := d.Set("description", entity.Description); err != nil {
		return err
	}
	if err := d.Set("default_column", entity.DefaultColumn); err != nil {
		return err
	}
	if err := d.Set("labels", entity.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(entity.Attributes)); err != nil {
		return err
	}
	return err
}

func resourceEntityCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entity := Entity{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		DefaultColumn: d.Get("default_column").(string),
		Labels:        expandStringList(d.Get("labels").([]interface{})),
		Attributes:    expandAttributes(d),
	}

	e, err := c.CreateEntity(entity)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceEntityUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entityID := d.Id()
	entity := Entity{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		DefaultColumn: d.Get("default_column").(string),
		Labels:        expandStringList(d.Get("labels").([]interface{})),
		Attributes:    expandAttributes(d),
	}

	err := c.UpdateEntity(entityID, entity)
	if err != nil {
		return err
	}

	return nil
}

func resourceEntityDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entityID := d.Id()

	err := c.DeleteEntity(entityID)
	if err != nil {
		return err
	}

	return nil
}
