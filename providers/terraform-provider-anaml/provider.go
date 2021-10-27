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
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANAML_DEFAULT_BRANCH", nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"anaml_entity":            anaml.DataSourceEntity(),
			"anaml_entity_population": anaml.DataSourceEntityPopulation(),
			"anaml_table":             anaml.DataSourceTable(),
			"anaml_feature":           anaml.DataSourceFeature(),
			"anaml_feature_set":       anaml.DataSourceFeatureSet(),
			"anaml_feature_template":  anaml.DataSourceFeatureTemplate(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"anaml_entity":            anaml.ResourceEntity(),
			"anaml_entity_mapping":    anaml.ResourceEntityMapping(),
			"anaml_entity_population": anaml.ResourceEntityPopulation(),
			"anaml_table":             anaml.ResourceTable(),
			"anaml_feature":           anaml.ResourceFeature(),
			"anaml_feature_set":       anaml.ResourceFeatureSet(),
			"anaml_feature_template":  anaml.ResourceFeatureTemplate(),
		},

		ConfigureFunc: providerConfigure,
	}
	return &provider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	branch := d.Get("branch").(string)

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	c, err := anaml.NewClient(host, &username, &password, &branch)
	if err != nil {
		return nil, err
	}

	return c, nil
}
