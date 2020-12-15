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
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
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
	if err := d.Set("namespace", FeatureStore.Namespace); err != nil {
		return err
	}
	if err := d.Set("mode", FeatureStore.Mode); err != nil {
		return err
	}
	return err
}

func resourceFeatureStoreCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	featureSet, _ := strconv.Atoi(d.Get("feature_set").(string))

	FeatureStore := FeatureStore{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		FeatureSet:  featureSet,
		Namespace:   d.Get("namespace").(string),
		Mode:        strings.ToTitle(d.Get("mode").(string)),
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
	FeatureStoreID := d.Id()

	FeatureStore := FeatureStore{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		FeatureSet:  featureSet,
		Namespace:   d.Get("namespace").(string),
		Mode:        strings.ToTitle(d.Get("mode").(string)),
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
