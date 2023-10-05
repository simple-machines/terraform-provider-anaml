package anaml

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const clusterDesc = `# Clusters

A Cluster represents an external compute resource that can be used to compute feature values,
run monitoring jobs and and generate previews.

There are two types of Clusters, a local cluster, and an Anaml Spark Server.

### Local Clusters

Local Clusters are rarely used outside of a development environment, and use the Anaml Server
application as a Spark Cluster in local mode.

### Anaml Spark Server

This form of cluster uses a separate microservice which runs a web application and launches
jobs on a Spark cluster. For this form of Cluster, you're required to provide the URL of the
spark server application.

Anaml Spark Server Clusters can be:

- Google Dataproc clusters
- Amazon EMR clusters
- Azure HD Insight clusters
- Hadoop Yarn clusters
- Spark on Kubernetes clusters
`

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Description: clusterDesc,
		Create:      resourceClusterCreate,
		Read:        resourceClusterRead,
		Update:      resourceClusterUpdate,
		Delete:      resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the cluster.",
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_preview_cluster": {
				Type:        schema.TypeBool,
				Description: "Whether this cluster can be used for preview generation.",
				Required:    true,
			},
			"spark_config": {
				Type:        schema.TypeList,
				Description: "Additional configuration which is passed to Spark when performing feature generation runs.",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        sparkConfigSchema(),
			},
			"property_set": {
				Type:        schema.TypeSet,
				Description: "Property Set with Additional configuration which is passed to Spark when performing feature generation runs.",
				Optional:    true,
				Elem:        propertySetSchema(),
			},
			"local": {
				Type:         schema.TypeList,
				Description:  "Set up for a local cluster. When this setting is used, a local spark session will be launched within the JVM process of the web server. Not recommended for production deployments.",
				Optional:     true,
				MaxItems:     1,
				Elem:         localSchema(),
				ExactlyOneOf: []string{"local", "spark_server"},
			},
			"spark_server": {
				Type:        schema.TypeList,
				Description: "Set up for a remote cluster.",
				Optional:    true,
				MaxItems:    1,
				Elem:        sparkServerSchema(),
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Labels to attach to the object",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"attribute": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
		},
	}
}

func sparkConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_hive_support": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"hive_metastore_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"additional_spark_properties": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					return make(map[string]interface{}), nil
				},
			},
		},
	}
}

func propertySetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of cluster property set. If specified, all values must be unique.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the cluster property set.",
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"additional_spark_properties": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				DefaultFunc: func() (interface{}, error) {
					return make(map[string]interface{}), nil
				},
			},
		},
	}
}

func localSchema() *schema.Resource {
	schemaMap := loginCredentialsProviderConfigSchema().Schema
	schemaMap["anaml_server_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
	}
	return &schema.Resource{
		Schema: schemaMap,
	}
}

func sparkServerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spark_server_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func loginCredentialsProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"basic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     basicCredentialsProviderConfigSchema(),
			},
			"file": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     fileCredentialsProviderConfigSchema(),
			},
			"aws": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     awsCredentialsProviderConfigSchema(),
			},
			"gcp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     gcpCredentialsProviderConfigSchema(),
			},
		},
	}
}

func basicCredentialsProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func fileCredentialsProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"filepath": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func awsCredentialsProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"password_secret_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func gcpCredentialsProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"password_secret_project": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"password_secret_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	clusterID := d.Id()

	cluster, err := c.GetCluster(clusterID)
	if err != nil {
		return err
	}
	if cluster == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", cluster.Name); err != nil {
		return err
	}
	if err := d.Set("description", cluster.Description); err != nil {
		return err
	}
	if err := d.Set("is_preview_cluster", cluster.IsPreviewCluster); err != nil {
		return err
	}
	sparkConfig, err := parseSparkConfig(cluster.SparkConfig)
	if sparkConfig == nil || err != nil {
		return err
	}
	if err := d.Set("spark_config", sparkConfig); err != nil {
		return err
	}

	if err := d.Set("property_set", flattenPropertySet(cluster.PropertySet)); err != nil {
		return err
	}

	if cluster.Type == "local" {
		local, err := parseLocal(cluster)
		if local == nil || err != nil {
			return err
		}
		if err := d.Set("local", local); err != nil {
			return err
		}
	}

	if cluster.Type == "sparkserver" {
		sparkServer, err := parseSparkServer(cluster)
		if sparkServer == nil || err != nil {
			return err
		}
		if err := d.Set("spark_server", sparkServer); err != nil {
			return err
		}
	}

	if err := d.Set("labels", cluster.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(cluster.Attributes)); err != nil {
		return err
	}
	return err
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	cluster, err := composeCluster(d)
	if cluster == nil || err != nil {
		return err
	}

	e, err := c.CreateCluster(*cluster)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	resourceClusterRead(d, m)
	return err
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	clusterID := d.Id()
	cluster, err := composeCluster(d)
	if cluster == nil || err != nil {
		return err
	}

	err = c.UpdateCluster(clusterID, *cluster)
	if err != nil {
		return err
	}

	resourceClusterRead(d, m)
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	clusterID := d.Id()

	err := c.DeleteCluster(clusterID)
	if err != nil {
		return err
	}

	return nil
}

func parseLocal(cluster *Cluster) ([]map[string]interface{}, error) {
	if cluster == nil {
		return nil, errors.New("Cluster is null")
	}

	local, err := parseLoginCredentialsProviderConfig(cluster.CredentialsProvider)
	if local == nil || err != nil {
		return nil, err
	}
	local["anaml_server_url"] = cluster.AnamlServerURL

	locals := make([]map[string]interface{}, 0, 1)
	locals = append(locals, local)
	return locals, nil
}

func parseLoginCredentialsProviderConfig(credentials *LoginCredentialsProviderConfig) (map[string]interface{}, error) {
	if credentials == nil {
		return nil, errors.New("LoginCredentialsProviderConfig is null")
	}

	provider := make(map[string]interface{})

	if credentials.Type == "basic" {
		basic := make(map[string]interface{})
		basic["username"] = credentials.Username
		basic["password"] = credentials.Password

		basics := make([]map[string]interface{}, 0, 1)
		basics = append(basics, basic)
		provider["basic"] = basics
	} else if credentials.Type == "file" {
		file := make(map[string]interface{})
		file["username"] = credentials.Username
		file["filepath"] = credentials.FilePath

		files := make([]map[string]interface{}, 0, 1)
		files = append(files, file)
		provider["file"] = files
	} else if credentials.Type == "aws" {
		aws := make(map[string]interface{})
		aws["username"] = credentials.Username
		aws["password_secret_id"] = credentials.PasswordSecretId

		awss := make([]map[string]interface{}, 0, 1)
		awss = append(awss, aws)
		provider["aws"] = awss
	} else if credentials.Type == "gcp" {
		gcp := make(map[string]interface{})
		gcp["username"] = credentials.Username
		gcp["password_secret_project"] = credentials.PasswordSecretProject
		gcp["password_secret_id"] = credentials.PasswordSecretId

		gcps := make([]map[string]interface{}, 0, 1)
		gcps = append(gcps, gcp)
		provider["gcp"] = gcps
	} else {
		return nil, fmt.Errorf("LoginCredentialsProviderConfig.Type contains an unexpected value: %s", credentials.Type)
	}

	return provider, nil
}

func parseSparkServer(cluster *Cluster) ([]map[string]interface{}, error) {
	if cluster == nil {
		return nil, errors.New("Cluster is null")
	}

	sparkServer := make(map[string]interface{})
	sparkServer["spark_server_url"] = cluster.SparkServerURL

	sparkServers := make([]map[string]interface{}, 0, 1)
	sparkServers = append(sparkServers, sparkServer)
	return sparkServers, nil
}

func parseSparkConfig(config *SparkConfig) ([]map[string]interface{}, error) {
	if config == nil {
		return nil, errors.New("SparkConfig is null")
	}

	sparkConfig := make(map[string]interface{})
	sparkConfig["enable_hive_support"] = config.EnableHiveSupport
	sparkConfig["hive_metastore_url"] = config.HiveMetastoreURL
	sparkConfig["additional_spark_properties"] = config.AdditionalSparkProperties

	sparkConfigs := make([]map[string]interface{}, 0, 1)
	sparkConfigs = append(sparkConfigs, sparkConfig)
	return sparkConfigs, nil
}

func composeCluster(d *schema.ResourceData) (*Cluster, error) {
	sparkConfigMap, err := expandSingleMap(d.Get("spark_config"))
	if err != nil {
		return nil, err
	}
	sparkConfig := composeSparkConfig(sparkConfigMap)

	if local, _ := expandSingleMap(d.Get("local")); local != nil {
		credentialsProvider, err := composeLoginCredentialsProviderConfig(local)
		if credentialsProvider == nil || err != nil {
			return nil, err
		}

		cluster := Cluster{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			Type:                "local",
			IsPreviewCluster:    d.Get("is_preview_cluster").(bool),
			AnamlServerURL:      local["anaml_server_url"].(string),
			CredentialsProvider: credentialsProvider,
			SparkConfig:         &sparkConfig,
			PropertySet:         expandPropertySet(d),
			Labels:              expandLabels(d),
			Attributes:          expandAttributes(d),
		}
		return &cluster, nil
	}

	if sparkServer, _ := expandSingleMap(d.Get("spark_server")); sparkServer != nil {
		cluster := Cluster{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Type:             "sparkserver",
			IsPreviewCluster: d.Get("is_preview_cluster").(bool),
			SparkServerURL:   sparkServer["spark_server_url"].(string),
			SparkConfig:      &sparkConfig,
			PropertySet:      expandPropertySet(d),
			Labels:           expandLabels(d),
			Attributes:       expandAttributes(d),
		}
		return &cluster, nil
	}

	return nil, errors.New("Invalid cluster type")
}

func composeLoginCredentialsProviderConfig(d map[string]interface{}) (*LoginCredentialsProviderConfig, error) {
	if basic, _ := expandSingleMap(d["basic"]); basic != nil {
		provider := LoginCredentialsProviderConfig{
			Type:     "basic",
			Username: basic["username"].(string),
			Password: basic["password"].(string),
		}
		return &provider, nil
	}

	if file, _ := expandSingleMap(d["file"]); file != nil {
		provider := LoginCredentialsProviderConfig{
			Type:     "file",
			Username: file["username"].(string),
			FilePath: file["filepath"].(string),
		}
		return &provider, nil
	}

	if aws, _ := expandSingleMap(d["aws"]); aws != nil {
		provider := LoginCredentialsProviderConfig{
			Type:             "aws",
			Username:         aws["username"].(string),
			PasswordSecretId: aws["password_secret_id"].(string),
		}
		return &provider, nil
	}

	if gcp, _ := expandSingleMap(d["gcp"]); gcp != nil {
		provider := LoginCredentialsProviderConfig{
			Type:                  "gcp",
			Username:              gcp["username"].(string),
			PasswordSecretProject: gcp["password_secret_project"].(string),
			PasswordSecretId:      gcp["password_secret_id"].(string),
		}
		return &provider, nil
	}

	return nil, errors.New("Invalid login credentials provider config type")
}

func composeSparkConfig(d map[string]interface{}) SparkConfig {
	source := d["additional_spark_properties"].(map[string]interface{})
	additionalSparkProperties := make(map[string]string)

	for k, v := range source {
		additionalSparkProperties[k] = v.(string)
	}
	return SparkConfig{
		EnableHiveSupport:         d["enable_hive_support"].(bool),
		HiveMetastoreURL:          d["hive_metastore_url"].(string),
		AdditionalSparkProperties: additionalSparkProperties,
	}
}

func expandPropertySet(d *schema.ResourceData) []PropertySet {
	prev, new := d.GetChange("property_set")
	prevX := prev.(*schema.Set).List()
	drs := new.(*schema.Set).List()
	res := make([]PropertySet, 0, len(drs))
	for _, dr := range drs {
		vals := dr.(map[string]interface{})
		source := vals["additional_spark_properties"].(map[string]interface{})
		additionalSparkProperties := make(map[string]string)

		for k, v := range source {
			additionalSparkProperties[k] = v.(string)
		}

		if vals["name"].(string) != "" {
			var identifier *int
			name := vals["name"].(string)
			explicit, hasId := vals["id"].(int)

			if hasId && explicit != 0 {
				identifier = &explicit
			} else {
				// No identifier was explicitly specified.
				// Try to find a property set from the last read
				// with the same name, we'll use its ID.
				for _, drX := range prevX {
					oldVals, _ := drX.(map[string]interface{})
					implicit, hasOld := oldVals["id"].(int)
					oldName := oldVals["name"].(string)
					if oldName == name && hasOld && implicit != 0 {
						identifier = &implicit
					}
				}
			}

			parsed := PropertySet{
				ID:                        identifier,
				Name:                      name,
				AdditionalSparkProperties: additionalSparkProperties,
			}

			res = append(res, parsed)
		}
	}
	return res
}

func flattenPropertySet(ps []PropertySet) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(ps))
	for _, ps := range ps {
		single := make(map[string]interface{})
		single["id"] = ps.ID
		single["name"] = ps.Name
		single["additional_spark_properties"] = ps.AdditionalSparkProperties
		res = append(res, single)
	}
	return res
}
