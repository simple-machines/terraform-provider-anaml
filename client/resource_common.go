package anaml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func destinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"folder": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     folderDestinationSchema(),
			},
			"table": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     tableDestinationSchema(),
			},
			"topic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     topicDestinationSchema(),
			},
			"option": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
		},
	}
}

func folderDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"partitioning_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"save_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"overwrite", "ignore", "append", "errorifexists",
				}, false),
			},
		},
	}
}

func tableDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"save_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"overwrite", "ignore", "append", "errorifexists",
				}, false),
			},
		},
	}
}

func topicDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"json", "avro",
				}, false),
			},
		},
	}
}

func expandAttributes(d *schema.ResourceData) []Attribute {
	drs := d.Get("attribute").(*schema.Set).List()
	return expandAttributesFromInterfaces(drs)
}

func expandAttributesFromInterfaces(drs []interface{}) []Attribute {
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
