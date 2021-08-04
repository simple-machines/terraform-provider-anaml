package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"merge_requests": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"merge_request_comments": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"commits": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"feature_store_runs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"monitoring_runs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
			"caching_runs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{},
			},
		},
	}
}

func resourceWebhookRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	webhookID := d.Id()

	webhook, err := c.GetWebhook(webhookID)
	if err != nil {
		return err
	}
	if webhook == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", webhook.Name); err != nil {
		return err
	}
	if err := d.Set("description", webhook.Description); err != nil {
		return err
	}
	if err := d.Set("url", webhook.URL); err != nil {
		return err
	}
	if err := d.Set("merge_requests", flattenEmpty(webhook.MergeRequests)); err != nil {
		return err
	}
	if err := d.Set("merge_request_comments", flattenEmpty(webhook.MergeRequestComments)); err != nil {
		return err
	}
	if err := d.Set("commits", flattenEmpty(webhook.Commits)); err != nil {
		return err
	}
	if err := d.Set("feature_store_runs", flattenEmpty(webhook.FeatureStoreRuns)); err != nil {
		return err
	}
	if err := d.Set("monitoring_runs", flattenEmpty(webhook.MonitoringRuns)); err != nil {
		return err
	}
	if err := d.Set("caching_runs", flattenEmpty(webhook.CachingRuns)); err != nil {
		return err
	}
	return err
}

func resourceWebhookCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	webhook := Webhook{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		URL:                  d.Get("url").(string),
		MergeRequests:        expandEmpty(d.Get("merge_requests").([]interface{})),
		MergeRequestComments: expandEmpty(d.Get("merge_request_comments").([]interface{})),
		Commits:              expandEmpty(d.Get("commits").([]interface{})),
		FeatureStoreRuns:     expandEmpty(d.Get("feature_store_runs").([]interface{})),
		MonitoringRuns:       expandEmpty(d.Get("monitoring_runs").([]interface{})),
		CachingRuns:          expandEmpty(d.Get("caching_runs").([]interface{})),
	}

	e, err := c.CreateWebhook(webhook)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	webhookID := d.Id()
	webhook := Webhook{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		URL:                  d.Get("url").(string),
		MergeRequests:        expandEmpty(d.Get("merge_requests").([]interface{})),
		MergeRequestComments: expandEmpty(d.Get("merge_request_comments").([]interface{})),
		Commits:              expandEmpty(d.Get("commits").([]interface{})),
		FeatureStoreRuns:     expandEmpty(d.Get("feature_store_runs").([]interface{})),
		MonitoringRuns:       expandEmpty(d.Get("monitoring_runs").([]interface{})),
		CachingRuns:          expandEmpty(d.Get("caching_runs").([]interface{})),
	}

	err := c.UpdateWebhook(webhookID, webhook)
	if err != nil {
		return err
	}

	return nil
}

func resourceWebhookDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	webhookID := d.Id()

	err := c.DeleteWebhook(webhookID)
	if err != nil {
		return err
	}

	return nil
}

func expandEmpty(drs []interface{}) *struct{} {
	if len(drs) == 0 {
		return nil
	} else {
		return &struct{}{}
	}
}

func flattenEmpty(specs *(struct{})) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if specs != nil {
		single := make(map[string]interface{})
		res = append(res, single)
	}

	return res
}
