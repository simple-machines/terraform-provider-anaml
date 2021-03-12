package anaml

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Create:   resourceClusterCreate,
		Read:     resourceClusterRead,
		Update:   resourceClusterUpdate,
		Delete:   resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State:  schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_preview_cluster": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"spark_config": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     sparkConfigSchema(),
			},
			"local": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				Elem:         localSchema(),
				ExactlyOneOf: []string{"local", "spark_server"},
			},
			"spark_server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     sparkServerSchema(),
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
				Type:   schema.TypeMap,
				Elem:   &schema.Schema {
					Type: schema.TypeString,
				},
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make(map[string]interface{}), nil
				},
			},
		},
	}
}

func localSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"anaml_server_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"jwt_token_provider": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     jwtTokenProviderSchema(),
			},
		},
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

func jwtTokenProviderSchema() *schema.Resource {
	schemaMap := loginCredentialsProviderConfigSchema().Schema
	schemaMap["login_server_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
	}
	return &schema.Resource{
		Schema: schemaMap,
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

	local := make(map[string]interface{})
	local["anaml_server_url"] = cluster.AnamlServerURL

	credentialsProvider, err := parseJWTTokenProviderConfig(cluster.CredentialsProvider)
	if credentialsProvider == nil || err != nil {
		return nil, err
	}
	local["jwt_token_provider"] = credentialsProvider

	locals := make([]map[string]interface{}, 0, 1)
	locals = append(locals, local)
	return locals, nil
}

func parseJWTTokenProviderConfig(config *JWTTokenProviderConfig) ([]map[string]interface{}, error) {
	if config == nil {
		return nil, errors.New("JWTTokenProviderConfig is null")
	}

	provider, err := parseLoginCredentialsProviderConfig(config.LoginCredentialsProviderConfig)
	if err != nil {
		return nil, err
	}

	provider["login_server_url"] = config.LoginServerURL

	providers := make([]map[string]interface{}, 0, 1)
	providers = append(providers, provider)
	return providers, nil
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
		credentialsProviderMap, err := expandSingleMap(local["jwt_token_provider"])
		if err != nil {
			return nil, err
		}
		credentialsProvider, err := composeJWTTokenProviderConfig(credentialsProviderMap)
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
		}
		return &cluster, nil
	}

	return nil, errors.New("Invalid cluster type")
}

func composeJWTTokenProviderConfig(d map[string]interface{}) (*JWTTokenProviderConfig, error) {
	loginProvider, err := composeLoginCredentialsProviderConfig(d)
	if loginProvider == nil || err != nil {
		return nil, err
	}

	provider := JWTTokenProviderConfig{
		Type:                           "jwtlogin",
		LoginServerURL:                 d["login_server_url"].(string),
		LoginCredentialsProviderConfig: loginProvider,
	}

	return &provider, nil
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
