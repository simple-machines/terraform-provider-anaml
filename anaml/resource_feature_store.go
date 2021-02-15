package anaml

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceFeatureStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureStoreCreate,
		Read:   resourceFeatureStoreRead,
		Update: resourceFeatureStoreUpdate,
		Delete: resourceFeatureStoreDelete,
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
			"feature_set": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"daily", "historical",
				}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToTitle(old) == strings.ToTitle(new)
				},
			},
			"destination": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     destinationSchema(),
			},
			"cluster": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"table_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceFeatureStoreRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStoreID := d.Id()

	FeatureStore, err := c.GetFeatureStore(FeatureStoreID)
	if err != nil {
		return err
	}
	if FeatureStore == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", FeatureStore.Name); err != nil {
		return err
	}
	if err := d.Set("description", FeatureStore.Description); err != nil {
		return err
	}
	if err := d.Set("feature_set", strconv.Itoa(FeatureStore.FeatureSet)); err != nil {
		return err
	}
	if err := d.Set("enabled", FeatureStore.Enabled); err != nil {
		return err
	}
	if err := d.Set("mode", FeatureStore.Mode); err != nil {
		return err
	}
	if err := d.Set("destination", flattenDestinationReferences(FeatureStore.Destinations)); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(FeatureStore.Cluster)); err != nil {
		return err
	}
	return err
}

func resourceFeatureStoreCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureSet, _ := strconv.Atoi(d.Get("feature_set").(string))
	cluster, _ := strconv.Atoi(d.Get("cluster").(string))

	FeatureStore := FeatureStore{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		FeatureSet:   featureSet,
		Enabled:      d.Get("enabled").(bool),
		Mode:         strings.ToTitle(d.Get("mode").(string)),
		Destinations: expandDestinationReferences(d),
		Cluster:      cluster,
	}

	e, err := c.CreateFeatureStore(FeatureStore)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.Id))
	return err
}

func resourceFeatureStoreUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureSet, _ := strconv.Atoi(d.Get("feature_set").(string))
	cluster, _ := strconv.Atoi(d.Get("cluster").(string))
	FeatureStoreID := d.Id()

	FeatureStore := FeatureStore{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		FeatureSet:   featureSet,
		Enabled:      d.Get("enabled").(bool),
		Mode:         strings.ToTitle(d.Get("mode").(string)),
		Destinations: expandDestinationReferences(d),
		Cluster:      cluster,
	}

	err := c.UpdateFeatureStore(FeatureStoreID, FeatureStore)
	if err != nil {
		return err
	}

	return nil
}

func resourceFeatureStoreDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	FeatureStoreID := d.Id()

	err := c.DeleteFeatureStore(FeatureStoreID)
	if err != nil {
		return err
	}

	return nil
}

func expandDestinationReferences(d *schema.ResourceData) []DestinationReference {
	drs := d.Get("destination").([]interface{})
	res := make([]DestinationReference, 0, len(drs))

	for _, dr := range drs {
		val, _ := dr.(map[string]interface{})
		destId, _ := strconv.Atoi(val["destination"].(string))

		parsed := DestinationReference{
			DestinationId: destId,
			Folder:        val["folder"].(string),
			TableName:     val["table_name"].(string),
		}
		res = append(res, parsed)
	}

	return res
}

func flattenDestinationReferences(destinations []DestinationReference) []interface{} {
	res := make([]interface{}, len(destinations))

	for _, destination := range destinations {
		single := make(map[string]interface{})
		single["destination"] = strconv.Itoa(destination.DestinationId)
		if destination.Folder != "" {
			single["folder"] = destination.Folder
		}
		if single["table_name"] != "" {
			single["table_name"] = destination.TableName
		}
		res = append(res, single)
	}

	return res
}
