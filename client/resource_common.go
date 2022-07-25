package anaml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func labelSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeString,
	}
}

func expandLabels(d *schema.ResourceData) []string {
	return expandStringList(d.Get("labels").(*schema.Set).List())
}

func attributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func expandAttributes(d *schema.ResourceData) []Attribute {
	drs := d.Get("attribute").(*schema.Set).List()
	res := make([]Attribute, 0, len(drs))
	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})
		parsed := Attribute{
			Key:   val["key"].(string),
			Value: val["value"].(string),
		}
		res = append(res, parsed)
	}
	return res
}

func flattenAttributes(attributes []Attribute) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(attributes))
	for _, attribute := range attributes {
		single := make(map[string]interface{})
		single["key"] = attribute.Key
		single["value"] = attribute.Value
		res = append(res, single)
	}
	return res
}
