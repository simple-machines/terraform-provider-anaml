package main

import (
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"anaml-operations_cluster":     anaml.DataSourceCluster(),
			"anaml-operations_destination": anaml.DataSourceDestination(),
			"anaml-operations_source":      anaml.DataSourceSource(),
			"anaml-operations_feature":     anaml.DataSourceFeature(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"anaml-operations_cluster":     anaml.ResourceCluster(),
			"anaml-operations_destination": anaml.ResourceDestination(),
			"anaml-operations_source":      anaml.ResourceSource(),
			"anaml-operations_user":        anaml.ResourceUser(),
			"anaml-operations_monitoring":  anaml.ResourceTableMonitoring(),
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

	c, err := anaml.NewClient(host, &username, &password, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
