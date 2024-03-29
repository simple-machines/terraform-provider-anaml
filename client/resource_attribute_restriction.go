package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const attributeDescription = `# Attribute Restrictions

An Attribute is a key/value pair for user-defined metadata. Restrictions limit the attributes
that can be applied to a given object, and what values they can take.

Multiple different types of attributes are supported:

- Enum ("Choice")
- Free Text
- Boolean
- Integer
- User (Assign an Anaml user id to the attribute)
- User Group (Assign an Anaml user group id to the attribute)
`

func ResourceAttributeRestriction() *schema.Resource {
	return &schema.Resource{
		Description: attributeDescription,
		Create:      resourceAttributeRestrictionCreate,
		Read:        resourceAttributeRestrictionRead,
		Update:      resourceAttributeRestrictionUpdate,
		Delete:      resourceAttributeRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mandatory": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enum": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				Elem:         enumAttributeSchema(),
				ExactlyOneOf: []string{"enum", "freetext", "boolean", "integer", "user", "user_group"},
			},
			"freetext": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"boolean": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"integer": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"user": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"user_group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"applies_to": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"cluster", "destination", "entity", "feature", "feature_set",
						"feature_store", "feature_template", "source", "table",
					}, false),
				},
			},
		},
	}
}

func enumAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"choice": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     enumChoiceSchema(),
			},
		},
	}
}

func enumChoiceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"display_emoji": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_colour": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAttributeRestrictionRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	attributeID := d.Id()

	attribute, err := c.GetAttributeRestriction(attributeID)
	if err != nil {
		return err
	}
	if attribute == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("key", attribute.Key); err != nil {
		return err
	}
	if err := d.Set("description", attribute.Description); err != nil {
		return err
	}
	if err := d.Set("mandatory", attribute.Mandatory); err != nil {
		return err
	}
	if err := d.Set("default_value", attribute.DefaultValue); err != nil {
		return err
	}

	if attribute.Type == "enumattribute" {
		e, err := parseEnumAttribute(attribute)
		if err != nil {
			return err
		}
		if err := d.Set("enum", e); err != nil {
			return err
		}
	} else {
		if err := d.Set("enum", nil); err != nil {
			return err
		}
	}

	if err := d.Set("freetext", buildEmpty(attribute.Type == "freetextattribute")); err != nil {
		return err
	}
	if err := d.Set("boolean", buildEmpty(attribute.Type == "booleanattribute")); err != nil {
		return err
	}
	if err := d.Set("integer", buildEmpty(attribute.Type == "integerattribute")); err != nil {
		return err
	}
	if err := d.Set("user", buildEmpty(attribute.Type == "userattribute")); err != nil {
		return err
	}
	if err := d.Set("user_group", buildEmpty(attribute.Type == "usergroupattribute")); err != nil {
		return err
	}
	if err := d.Set("applies_to", mapTargetsToFrontend(attribute.AppliesTo)); err != nil {
		return err
	}
	return err
}

func resourceAttributeRestrictionCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	attribute, err := composeAttribute(d)
	if attribute == nil || err != nil {
		return err
	}

	a, err := c.CreateAttributeRestriction(*attribute)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(a.ID))
	return err
}

func resourceAttributeRestrictionUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	attributeID := d.Id()
	attribute, err := composeAttribute(d)
	if attribute == nil || err != nil {
		return err
	}

	err = c.UpdateAttributeRestriction(attributeID, *attribute)
	if err != nil {
		return err
	}

	return nil
}

func resourceAttributeRestrictionDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	attributeID := d.Id()

	err := c.DeleteAttributeRestriction(attributeID)
	if err != nil {
		return err
	}

	return nil
}

func composeAttribute(d *schema.ResourceData) (*AttributeRestriction, error) {
	appliesTo := mapTargetsToBackend(expandStringList(d.Get("applies_to").(*schema.Set).List()))

	attribute := AttributeRestriction{
		Key:         d.Get("key").(string),
		Description: d.Get("description").(string),
		Mandatory:   d.Get("mandatory").(bool),
		AppliesTo:   appliesTo,
	}

	if default_value, set := d.GetOk("default_value"); set {
		default_value := default_value.(string)
		attribute.DefaultValue = &default_value
	}

	if e, _ := expandSingleMap(d.Get("enum")); e != nil {
		choices, err := expandEnumChoices(e["choice"].(*schema.Set).List())
		if err != nil {
			return nil, err
		}
		attribute.Type = "enumattribute"
		attribute.Choices = &choices
		return &attribute, nil
	}

	if existsEmpty(d.Get("freetext").([]interface{})) {
		attribute.Type = "freetextattribute"
		return &attribute, nil
	}

	if existsEmpty(d.Get("boolean").([]interface{})) {
		attribute.Type = "booleanattribute"
		return &attribute, nil
	}

	if existsEmpty(d.Get("integer").([]interface{})) {
		attribute.Type = "integerattribute"
		return &attribute, nil
	}

	if existsEmpty(d.Get("user").([]interface{})) {
		attribute.Type = "userattribute"
		return &attribute, nil
	}

	if existsEmpty(d.Get("user_group").([]interface{})) {
		attribute.Type = "usergroupattribute"
		return &attribute, nil
	}
	return nil, errors.New("Invalid attribute type")
}

func parseEnumAttribute(attribute *AttributeRestriction) ([]map[string]interface{}, error) {
	if attribute == nil {
		return nil, errors.New("Attribute Restriction is null")
	}

	e := make(map[string]interface{})
	choices := flattenEnumChoices(*attribute.Choices)
	e["choice"] = choices

	es := make([]map[string]interface{}, 0, 1)
	es = append(es, e)
	return es, nil
}

func parseNonEnumAttribute(attribute *AttributeRestriction) ([]map[string]interface{}, error) {
	if attribute == nil {
		return nil, errors.New("Attribute Restriction is null")
	}

	nea := make(map[string]interface{})

	neas := make([]map[string]interface{}, 0, 1)
	neas = append(neas, nea)
	return neas, nil
}

func mapTargetsToFrontend(backend []AttributeTarget) []string {
	vs := make([]string, 0, len(backend))
	for _, v := range backend {
		if v.Type == "cluster" {
			vs = append(vs, "cluster")
		} else if v.Type == "destination" {
			vs = append(vs, "destination")
		} else if v.Type == "entity" {
			vs = append(vs, "entity")
		} else if v.Type == "feature" {
			vs = append(vs, "feature")
		} else if v.Type == "featureset" {
			vs = append(vs, "feature_set")
		} else if v.Type == "featurestore" {
			vs = append(vs, "feature_store")
		} else if v.Type == "featuretemplate" {
			vs = append(vs, "feature_template")
		} else if v.Type == "source" {
			vs = append(vs, "source")
		} else if v.Type == "table" {
			vs = append(vs, "table")
		}
		// TODO: We should raise an error if we fall through the cases.
	}
	return vs
}

func mapTargetsToBackend(frontend []string) []AttributeTarget {
	vs := make([]AttributeTarget, 0, len(frontend))
	for _, v := range frontend {
		if v == "cluster" {
			vs = append(vs, AttributeTarget{"cluster"})
		} else if v == "destination" {
			vs = append(vs, AttributeTarget{"destination"})
		} else if v == "entity" {
			vs = append(vs, AttributeTarget{"entity"})
		} else if v == "feature" {
			vs = append(vs, AttributeTarget{"feature"})
		} else if v == "feature_set" {
			vs = append(vs, AttributeTarget{"featureset"})
		} else if v == "feature_store" {
			vs = append(vs, AttributeTarget{"featurestore"})
		} else if v == "feature_template" {
			vs = append(vs, AttributeTarget{"featuretemplate"})
		} else if v == "source" {
			vs = append(vs, AttributeTarget{"source"})
		} else if v == "table" {
			vs = append(vs, AttributeTarget{"table"})
		}
		// TODO: We should raise an error if we fall through the cases.
	}
	return vs
}

func expandEnumChoices(choices []interface{}) ([]EnumAttributeChoice, error) {
	res := make([]EnumAttributeChoice, 0, len(choices))

	for _, choice := range choices {
		val, _ := choice.(map[string]interface{})

		var parsed EnumAttributeChoice
		var display EnumAttributeDisplay
		display_emoji := ""
		display_colour := ""

		if de, ok := val["display_emoji"]; ok {
			display_emoji = de.(string)
		}
		if dc, ok := val["display_colour"]; ok {
			display_colour = dc.(string)
		}

		if display_emoji != "" || display_colour != "" {
			display = EnumAttributeDisplay{
				Emoji:  display_emoji,
				Colour: display_colour,
			}
			parsed = EnumAttributeChoice{
				Value:   val["value"].(string),
				Display: &display,
			}
		} else {
			parsed = EnumAttributeChoice{
				Value:   val["value"].(string),
				Display: nil,
			}
		}

		res = append(res, parsed)
	}

	return res, nil
}

func flattenEnumChoices(choices []EnumAttributeChoice) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(choices))
	for _, choice := range choices {
		single := make(map[string]interface{})
		single["value"] = choice.Value
		if choice.Display != nil {
			single["display_emoji"] = choice.Display.Emoji
			single["display_colour"] = choice.Display.Colour
		}
		res = append(res, single)
	}
	return res
}

func existsEmpty(drs []interface{}) bool {
	if len(drs) == 0 {
		return false
	} else {
		return true
	}
}

func buildEmpty(specs bool) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if specs {
		single := make(map[string]interface{})
		res = append(res, single)
	}

	return res
}
