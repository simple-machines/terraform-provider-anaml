package anaml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"fmt"
	"strconv"
)

func getAnamlId(d *schema.ResourceData, key string) (int, error) {
	if raw, ok := d.GetOk(key); ok {
		if i, err := strconv.Atoi(raw.(string)); err == nil {
			return i, nil
		} else {
			return 0, err
		}
	}
	return 0, fmt.Errorf("Required Identifier %s is missing", key)
}

func getAnamlIdPointer(d *schema.ResourceData, key string) *int {
	if raw, ok := d.GetOk(key); ok {
		if i, err := strconv.Atoi(raw.(string)); err == nil {
			return &i
		} else {
			return nil
		}
	}
	return nil
}

func getStringPointer(d *schema.ResourceData, key string) *string {
	if raw, ok := d.GetOk(key); ok {
		if i, ok := raw.(string); ok {
			return &i
		} else {
			return nil
		}
	}
	return nil
}

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

func expandDestinationReferences(drs []interface{}) ([]DestinationReference, error) {
	res := make([]DestinationReference, 0, len(drs))

	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})

		destID, _ := strconv.Atoi(val["destination"].(string))
		options := expandAttributesFromInterfaces(val["option"].(*schema.Set).List())

		parsed := DestinationReference{
			DestinationID: destID,
			Options:       options,
		}

		if folder, _ := expandSingleMap(val["folder"]); folder != nil {
			if path, ok := folder["path"].(string); ok {
				parsed.Type = "folder"
				parsed.Folder = path
				enabled := folder["partitioning_enabled"].(bool)
				parsed.FolderPartitioningEnabled = &enabled
				mode := folder["save_mode"].(string)
				parsed.Mode = mode
			} else {
				return nil, fmt.Errorf("error casting table.path %i", folder["path"])
			}
		}

		if table, _ := expandSingleMap(val["table"]); table != nil {
			if tableName, ok := table["name"].(string); ok {
				parsed.Type = "table"
				parsed.TableName = tableName
				if mode, _ := table["save_mode"].(string); mode != "" {
					parsed.Mode = mode
				}
			} else {
				return nil, fmt.Errorf("error casting table.name %i", table["name"])
			}
		}

		if topic, _ := expandSingleMap(val["topic"]); topic != nil {
			if topicName, ok := topic["name"].(string); ok {
				parsed.Type = "topic"
				parsed.Topic = topicName
				parsed.Format = &KafkaFormat{
					Type: topic["format"].(string),
				}
			} else {
				return nil, fmt.Errorf("error casting topic.name %i", topic["name"])
			}
		}

		res = append(res, parsed)
	}

	return res, nil
}

func flattenDestinationReferences(destinations []DestinationReference) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0, len(destinations))

	for _, destination := range destinations {
		single := make(map[string]interface{})
		single["destination"] = strconv.Itoa(destination.DestinationID)
		if destination.Options != nil {
			single["option"] = flattenAttributes(destination.Options)
		}

		if destination.Type == "folder" {
			folder := make(map[string]interface{})
			folder["path"] = destination.Folder
			folder["partitioning_enabled"] = destination.FolderPartitioningEnabled
			folder["save_mode"] = destination.Mode

			folders := make([]map[string]interface{}, 0, 1)
			folders = append(folders, folder)
			single["folder"] = folders
		}

		if destination.Type == "table" {
			table := make(map[string]interface{})
			table["name"] = destination.TableName
			table["save_mode"] = destination.Mode

			tables := make([]map[string]interface{}, 0, 1)
			tables = append(tables, table)
			single["table"] = tables
		}

		if destination.Type == "topic" {
			topic := make(map[string]interface{})
			topic["name"] = destination.Topic
			topic["format"] = destination.Format.Type

			topics := make([]map[string]interface{}, 0, 1)
			topics = append(topics, topic)
			single["topic"] = topics
		}

		res = append(res, single)
	}

	return res, nil
}
