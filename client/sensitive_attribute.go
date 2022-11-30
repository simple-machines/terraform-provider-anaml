package anaml

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func parseSensitiveAttribute(sensitive *SensitiveAttribute) (map[string]interface{}, error) {
	if sensitive == nil {
		return nil, errors.New("SensitiveAttribute is null")
	}

	res := make(map[string]interface{})
	res["key"] = sensitive.Key

	provider, err := parseSecretProviderConfig(sensitive.ValueConfig)
	if err != nil {
		return nil, err
	}
	for k, v := range provider {
		res[k] = v
	}
	return res, nil
}

func parseSecretProviderConfig(secretProvider *SecretValueConfig) (map[string]interface{}, error) {
	if secretProvider == nil {
		return nil, errors.New("SecretValueConfig is null")
	}

	provider := make(map[string]interface{})

	if secretProvider.Type == "basic" {
		provider["value"] = secretProvider.Secret
	} else if secretProvider.Type == "file" {
		file := make(map[string]interface{})
		file["filepath"] = secretProvider.FilePath

		files := make([]map[string]interface{}, 0, 1)
		files = append(files, file)
		provider["file"] = files
	} else if secretProvider.Type == "awssm" {
		aws := make(map[string]interface{})
		aws["secret_id"] = secretProvider.SecretId

		awss := make([]map[string]interface{}, 0, 1)
		awss = append(awss, aws)
		provider["aws"] = awss
	} else if secretProvider.Type == "gcpsm" {
		gcp := make(map[string]interface{})
		gcp["secret_project"] = secretProvider.SecretProject
		gcp["secret_id"] = secretProvider.SecretId

		gcps := make([]map[string]interface{}, 0, 1)
		gcps = append(gcps, gcp)
		provider["gcp"] = gcps
	} else {
		return nil, fmt.Errorf("SecretValueConfig.Type contains an unexpected value: %s", secretProvider.Type)
	}

	return provider, nil
}

func composeSensitiveAttribute(d map[string]interface{}) (*SensitiveAttribute, error) {
	sensitive := SensitiveAttribute{
		Key: d["key"].(string),
	}

	if d["value"] != nil && d["value"] != "" {
		sensitive.ValueConfig = &SecretValueConfig{
			Type:   "basic",
			Secret: d["value"].(string),
		}
	} else if file, _ := expandSingleMap(d["file"]); file != nil {
		sensitive.ValueConfig = &SecretValueConfig{
			Type:     "file",
			FilePath: file["filepath"].(string),
		}
	} else if aws, _ := expandSingleMap(d["aws"]); aws != nil {
		sensitive.ValueConfig = &SecretValueConfig{
			Type:     "awssm",
			SecretId: aws["secret_id"].(string),
		}
	} else if gcp, _ := expandSingleMap(d["gcp"]); gcp != nil {
		sensitive.ValueConfig = &SecretValueConfig{
			Type:          "gcpsm",
			SecretProject: gcp["secret_project"].(string),
			SecretId:      gcp["secret_id"].(string),
		}
	} else {
		return nil, fmt.Errorf("SensitiveAttribute. Coudn't parse Sensitive Attribute")
	}

	return &sensitive, nil
}

func sensitiveAttributeSchema() *schema.Resource {
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
			"file": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     fileSecretProviderConfigSchema(),
			},
			"aws": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     awsSecretProviderConfigSchema(),
			},
			"gcp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// ExactlyOneOf: []string{"value", "aws", "gcp"},
				Elem: gcpSecretProviderConfigSchema(),
			},
		},
	}
}

func fileSecretProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filepath": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func awsSecretProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func gcpSecretProviderConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secret_project": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"secret_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}
