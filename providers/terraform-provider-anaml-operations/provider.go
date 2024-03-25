package main

import (
	"time"

	anaml "anaml.io/terraform/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANAML_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANAML_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ANAML_PASSWORD", nil),
			},
			"request_timeout": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "30s",
				ValidateFunc: anaml.ValidateDuration(),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"anaml-operations_cluster":       anaml.DataSourceCluster(),
			"anaml-operations_destination":   anaml.DataSourceDestination(),
			"anaml-operations_source":        anaml.DataSourceSource(),
			"anaml-operations_feature_store": anaml.DataSourceFeatureStore(),
			"anaml-operations_user":          anaml.DataSourceUser(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"anaml-operations_access_token":             anaml.ResourceAccessToken(),
			"anaml-operations_attribute_restriction":    anaml.ResourceAttributeRestriction(),
			"anaml-operations_branch_protection":        anaml.ResourceBranchProtection(),
			"anaml-operations_caching":                  anaml.ResourceTableCaching(),
			"anaml-operations_cluster":                  anaml.ResourceCluster(),
			"anaml-operations_destination":              anaml.ResourceDestination(),
			"anaml-operations_event_store":              anaml.ResourceEventStore(),
			"anaml-operations_feature_store":            anaml.ResourceFeatureStore(),
			"anaml-operations_metrics_job":              anaml.ResourceMetricsJob(),
			"anaml-operations_label_restriction":        anaml.ResourceLabelRestriction(),
			"anaml-operations_monitoring":               anaml.ResourceTableMonitoring(),
			"anaml-operations_source":                   anaml.ResourceSource(),
			"anaml-operations_user_group":               anaml.ResourceUserGroup(),
			"anaml-operations_user":                     anaml.ResourceUser(),
			"anaml-operations_view_materialisation_job": anaml.ResourceViewMaterialisationJob(),
			"anaml-operations_webhook":                  anaml.ResourceWebhook(),
		},

		ConfigureFunc: providerConfigure,
	}
	return &provider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	timeout, err := time.ParseDuration(d.Get("request_timeout").(string))

	if err != nil {
		return nil, err
	}

	c, err := anaml.NewClient(host, &username, &password, nil, timeout)
	if err != nil {
		return nil, err
	}

	return c, nil
}
