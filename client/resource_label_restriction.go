package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const labelDescription = `# Label Restrictions

A Label is a user-defined piece of metadata that allows for classification of resources.
Restrictions limit the labels that can be applied to a given object.
`

func ResourceLabelRestriction() *schema.Resource {
	return &schema.Resource{
		Description: labelDescription,
		Create:      resourceLabelRestrictionCreate,
		Read:        resourceLabelRestrictionRead,
		Update:      resourceLabelRestrictionUpdate,
		Delete:      resourceLabelRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"text": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"emoji": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"colour": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLabelRestrictionRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	labelID := d.Id()

	label, err := c.GetLabelRestriction(labelID)
	if err != nil {
		return err
	}
	if label == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("text", label.Text); err != nil {
		return err
	}
	if err := d.Set("emoji", label.Emoji); err != nil {
		return err
	}
	if err := d.Set("colour", label.Colour); err != nil {
		return err
	}

	return err
}

func resourceLabelRestrictionCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	label := composeLabel(d)
	l, err := c.CreateLabelRestriction(*label)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(l.ID))
	return err
}

func resourceLabelRestrictionUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	labelID := d.Id()
	label := composeLabel(d)
	err := c.UpdateLabelRestriction(labelID, *label)
	if err != nil {
		return err
	}

	return nil
}

func resourceLabelRestrictionDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	labelID := d.Id()

	err := c.DeleteLabelRestriction(labelID)
	if err != nil {
		return err
	}

	return nil
}

func composeLabel(d *schema.ResourceData) *LabelRestriction {
	label := LabelRestriction{
		Text:   d.Get("text").(string),
		Emoji:  getNullableString(d, "emoji"),
		Colour: getNullableString(d, "colour"),
	}
	return &label
}
