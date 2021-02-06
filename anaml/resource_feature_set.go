package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFeatureSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureSetCreate,
		Read:   resourceFeatureSetRead,
		Update: resourceFeatureSetUpdate,
		Delete: resourceFeatureSetDelete,
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
			"entity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"features": {
				Type:        schema.TypeSet,
				Description: "Features to include in the feature set",
				Required:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func resourceFeatureSetRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureSetID := d.Id()

	FeatureSet, err := c.GetFeatureSet(FeatureSetID)
	if err != nil {
		return err
	}
	if FeatureSet == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", FeatureSet.Name); err != nil {
		return err
	}
	if err := d.Set("description", FeatureSet.Description); err != nil {
		return err
	}
	if err := d.Set("entity", strconv.Itoa(FeatureSet.EntityId)); err != nil {
		return err
	}
	if err := d.Set("features", identifierList(FeatureSet.Features)); err != nil {
		return err
	}
	return err
}

func resourceFeatureSetCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entity, _ := strconv.Atoi(d.Get("entity").(string))

	FeatureSet := FeatureSet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		EntityId:    entity,
		Features:    expandIdentifierList(d.Get("features").(*schema.Set).List()),
	}

	e, err := c.CreateFeatureSet(FeatureSet)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.Id))
	return err
}

func resourceFeatureSetUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entity, _ := strconv.Atoi(d.Get("entity").(string))
	FeatureSetID := d.Id()

	FeatureSet := FeatureSet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		EntityId:    entity,
		Features:    expandIdentifierList(d.Get("features").(*schema.Set).List()),
	}

	err := c.UpdateFeatureSet(FeatureSetID, FeatureSet)
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureSetDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureSetID := d.Id()

	err := c.DeleteFeatureSet(FeatureSetID)
	if err != nil {
		return err
	}

	return nil
}
