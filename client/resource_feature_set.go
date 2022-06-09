package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const featureSetDescription = `# Feature Sets

A Feature Set is collection of features that are generated at the same time. A Feature Set would usually comprise of:

* the Features required to train and score a machine learning model; or
* the Features required in a report or dashboard

Feature Sets are often re-used over multiple Feature Stores to generate historical, daily or online outputs.

Each Feature Set is specific to an Entity. Once the Entity is selected, the list of Features
available to be chosen is restricted to Features for that Entity.
`

func ResourceFeatureSet() *schema.Resource {
	return &schema.Resource{
		Description: featureSetDescription,
		Create:      resourceFeatureSetCreate,
		Read:        resourceFeatureSetRead,
		Update:      resourceFeatureSetUpdate,
		Delete:      resourceFeatureSetDelete,
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
			"post_expression": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Expressions to apply to output post feature generation",
				Elem:        postExpressionSchema(),
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
	if err := d.Set("entity", strconv.Itoa(FeatureSet.EntityID)); err != nil {
		return err
	}
	if err := d.Set("features", identifierList(FeatureSet.Features)); err != nil {
		return err
	}
	if err := d.Set("labels", FeatureSet.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(FeatureSet.Attributes)); err != nil {
		return err
	}
	if err := d.Set("post_expressions", flattenPostExpressions(FeatureSet.PostExpressions)); err != nil {
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
		EntityID:    entity,
		Features:    expandIdentifierList(d.Get("features").(*schema.Set).List()),
		Labels:      expandStringList(d.Get("labels").([]interface{})),
		Attributes:  expandAttributes(d),
	}

	e, err := c.CreateFeatureSet(FeatureSet)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceFeatureSetUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entity, _ := strconv.Atoi(d.Get("entity").(string))
	FeatureSetID := d.Id()

	FeatureSet := FeatureSet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		EntityID:    entity,
		Features:    expandIdentifierList(d.Get("features").(*schema.Set).List()),
		Labels:      expandStringList(d.Get("labels").([]interface{})),
		Attributes:  expandAttributes(d),
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

func postExpressionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     filterPostExpressionSchema(),
			},
			"mask": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     maskPostExpressionSchema(),
			},
		},
	}
}

func filterPostExpressionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"expression": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func maskPostExpressionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"expression": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func flattenPostExpressions(postExpressions []PostExpression) []map[string]([]map[string]interface{}) {
	res := make([]map[string]([]map[string]interface{}), 0, len(postExpressions))
	for _, expr := range postExpressions {
		single := make(map[string]([]map[string]interface{}))
		if expr.Type == "filter" {
			nest := make(map[string]interface{})
			nest["expression"] = expr.Expression
			single["filter"] = []map[string]interface{}{nest}
		}
		if expr.Type == "mask" {
			nest := make(map[string]interface{})
			nest["column"] = expr.Column
			nest["expression"] = expr.Expression
			single["mask"] = []map[string]interface{}{nest}
		}
		res = append(res, single)
	}
	return res
}
