package anaml

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceDestinationCreate,
		Read:   resourceDestinationRead,
		Update: resourceDestinationUpdate,
		Delete: resourceDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"s3": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				Elem:         s3SourceDestinationSchema(),
				ExactlyOneOf: []string{"s3", "s3a", "jdbc", "hive", "big_query", "gcs", "local", "hdfs", "kafka", "snowflake"},
			},
			"s3a": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     s3aSourceDestinationSchema(),
			},
			"jdbc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     jdbcSourceDestinationSchema(),
			},
			"hive": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     hiveSourceDestinationSchema(),
			},
			"big_query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     bigQueryDestinationSchema(),
			},
			"gcs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     gcsSourceDestinationSchema(),
			},
			"local": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     localSourceDestinationSchema(),
			},
			"hdfs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     hdfsSourceDestinationSchema(),
			},
			"kafka": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     kafkaSourceDestinationSchema(),
			},
			"snowflake": {
            			Type:     schema.TypeList,
            			Optional: true,
            			MaxItems: 1,
            			Elem:     snowflakeSourceDestinationSchema(),
            		},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Labels to attach to the object",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"attribute": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
		},
	}
}

func bigQueryDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"temporary_staging_area": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     bigQueryTemporaryStagingAreaSchema(),
			},
			"persistent_staging_area": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     bigQueryPersistentStagingAreaSchema(),
			},
		},
	}
}

func bigQueryTemporaryStagingAreaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func bigQueryPersistentStagingAreaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceDestinationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	destinationID := d.Id()

	destination, err := c.GetDestination(destinationID)
	if err != nil {
		return err
	}
	if destination == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", destination.Name); err != nil {
		return err
	}
	if err := d.Set("description", destination.Description); err != nil {
		return err
	}

	if destination.Type == "s3" {
		s3, err := parseS3Destination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("s3", s3); err != nil {
			return err
		}
	}

	if destination.Type == "s3a" {
		s3a, err := parseS3ADestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("s3a", s3a); err != nil {
			return err
		}
	}

	if destination.Type == "gcs" {
		gcs, err := parseS3Destination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("gcs", gcs); err != nil {
			return err
		}
	}

	if destination.Type == "local" {
		local, err := parseLocalDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("local", local); err != nil {
			return err
		}
	}

	if destination.Type == "hdfs" {
		hdfs, err := parseLocalDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("hdfs", hdfs); err != nil {
			return err
		}
	}

	if destination.Type == "jdbc" {
		jdbc, err := parseJDBCDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("jdbc", jdbc); err != nil {
			return err
		}
	}

	if destination.Type == "hive" {
		hive, err := parseHiveDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("hive", hive); err != nil {
			return err
		}
	}

	if destination.Type == "bigquery" {
		bigQuery, err := parseBigQueryDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("big_query", bigQuery); err != nil {
			return err
		}
	}

	if destination.Type == "kafka" {
		kafka, err := parseKafkaDestination(destination)
		if err != nil {
			return err
		}
		if err := d.Set("kafka", kafka); err != nil {
			return err
		}
	}

    	if destination.Type == "snowflake" {
	    	snowflake, err := parseSnowflakeDestination(destination)
	    	if err != nil {
	    		return err
	    	}
	    	if err := d.Set("snowflake", snowflake); err != nil {
	    		return err
	    	}
	    }
	
	if err := d.Set("labels", destination.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(destination.Attributes)); err != nil {
		return err
	}
	return err
}

func resourceDestinationCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	destination, err := composeDestination(d)
	if destination == nil || err != nil {
		return err
	}

	e, err := c.CreateDestination(*destination)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceDestinationUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	destinationID := d.Id()
	destination, err := composeDestination(d)
	if destination == nil || err != nil {
		return err
	}

	err = c.UpdateDestination(destinationID, *destination)
	if err != nil {
		return err
	}

	return nil
}

func resourceDestinationDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	destinationID := d.Id()

	err := c.DeleteDestination(destinationID)
	if err != nil {
		return err
	}

	return nil
}

// Used for both S3 and GCS destinations
func parseS3Destination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	s3 := make(map[string]interface{})
	s3["bucket"] = destination.Bucket
	s3["path"] = destination.Path

	fileFormat := parseFileFormat(destination.FileFormat)
	for k, v := range fileFormat {
		s3[k] = v
	}

	s3s := make([]map[string]interface{}, 0, 1)
	s3s = append(s3s, s3)
	return s3s, nil
}

func parseS3ADestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	s3a := make(map[string]interface{})
	s3a["bucket"] = destination.Bucket
	s3a["path"] = destination.Path
	s3a["endpoint"] = destination.Endpoint
	s3a["access_key"] = destination.AccessKey
	s3a["secret_key"] = destination.SecretKey

	fileFormat := parseFileFormat(destination.FileFormat)
	for k, v := range fileFormat {
		s3a[k] = v
	}

	s3as := make([]map[string]interface{}, 0, 1)
	s3as = append(s3as, s3a)
	return s3as, nil
}

// Used for local and HDFS destinations
func parseLocalDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	local := make(map[string]interface{})
	local["path"] = destination.Path

	fileFormat := parseFileFormat(destination.FileFormat)
	for k, v := range fileFormat {
		local[k] = v
	}

	locals := make([]map[string]interface{}, 0, 1)
	locals = append(locals, local)
	return locals, nil
}

func parseJDBCDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	jdbc := make(map[string]interface{})
	jdbc["url"] = destination.URL
	jdbc["schema"] = destination.Schema

	credentialsProvider, err := parseLoginCredentialsProviderConfig(destination.CredentialsProvider)
	if err != nil {
		return nil, err
	}
	jdbc["credentials_provider"] = []map[string]interface{}{credentialsProvider}

	jdbcs := make([]map[string]interface{}, 0, 1)
	jdbcs = append(jdbcs, jdbc)
	return jdbcs, nil
}

func parseBigQueryDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	bigQuery, err := parseBigQueryStagingArea(destination.StagingArea)
	if err != nil {
		return nil, err
	}

	bigQuery["path"] = destination.Path

	return []map[string]interface{}{bigQuery}, nil
}

func parseBigQueryStagingArea(stagingArea *GCSStagingArea) (map[string]interface{}, error) {
	if stagingArea == nil {
		return nil, errors.New("GCSStagingArea is null")
	}

	stagingAreaMap := make(map[string]interface{})

	if stagingArea.Type == "temporary" {
		temporaryMap := map[string]interface{}{
			"bucket": stagingArea.Bucket,
		}
		stagingAreaMap["temporary_staging_area"] = []map[string]interface{}{temporaryMap}
	} else if stagingArea.Type == "persistent" {
		persistentMap := map[string]interface{}{
			"bucket": stagingArea.Bucket,
			"path":   stagingArea.Path,
		}
		stagingAreaMap["persistent_staging_area"] = []map[string]interface{}{persistentMap}
	} else {
		return nil, fmt.Errorf("Type contains an unrecognised value: '%s'", stagingArea.Type)
	}

	return stagingAreaMap, nil
}

func parseHiveDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	hive := make(map[string]interface{})
	hive["database"] = destination.Database

	hives := make([]map[string]interface{}, 0, 1)
	hives = append(hives, hive)
	return hives, nil
}

func parseKafkaDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	kafka := make(map[string]interface{})
	kafka["bootstrap_servers"] = destination.BootstrapServers
	kafka["schema_registry_url"] = destination.SchemaRegistryURL

	sensitives := make([]map[string]interface{}, len(destination.KafkaProperties))
	for i, v := range destination.KafkaProperties {
		sa, err := parseSensitiveAttribute(&v)
		if err != nil {
			return nil, err
		}

		sensitives[i] = sa
	}

	kafka["property"] = sensitives

	kafkas := make([]map[string]interface{}, 0, 1)
	kafkas = append(kafkas, kafka)
	return kafkas, nil
}

func parseSnowflakeDestination(destination *Destination) ([]map[string]interface{}, error) {
	if destination == nil {
		return nil, errors.New("Destination is null")
	}

	snowflake := make(map[string]interface{})
	snowflake["url"] = destination.URL
	snowflake["schema"] = destination.Schema
        snowflake["database"] = destination.Database
        snowflake["warehouse"] = destination.Warehouse

	credentialsProvider, err := parseLoginCredentialsProviderConfig(destination.CredentialsProvider)
	if err != nil {
		return nil, err
	}
	snowflake["credentials_provider"] = []map[string]interface{}{credentialsProvider}

	snowflakes := make([]map[string]interface{}, 0, 1)
	snowflakes = append(snowflakes, snowflake)
	return snowflakes, nil
}

func composeDestination(d *schema.ResourceData) (*Destination, error) {
	if s3, _ := expandSingleMap(d.Get("s3")); s3 != nil {
		fileFormat := composeFileFormat(s3)
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "s3",
			Bucket:      s3["bucket"].(string),
			Path:        s3["path"].(string),
			FileFormat:  fileFormat,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if s3a, _ := expandSingleMap(d.Get("s3a")); s3a != nil {
		fileFormat := composeFileFormat(s3a)
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "s3a",
			Bucket:      s3a["bucket"].(string),
			Path:        s3a["path"].(string),
			Endpoint:    s3a["endpoint"].(string),
			AccessKey:   s3a["access_key"].(string),
			SecretKey:   s3a["secret_key"].(string),
			FileFormat:  fileFormat,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if jdbc, _ := expandSingleMap(d.Get("jdbc")); jdbc != nil {
		credentialsProviderMap, err := expandSingleMap(jdbc["credentials_provider"])
		if err != nil {
			return nil, err
		}

		credentialsProvider, err := composeLoginCredentialsProviderConfig(credentialsProviderMap)
		if err != nil {
			return nil, err
		}

		destination := Destination{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			Type:                "jdbc",
			URL:                 jdbc["url"].(string),
			Schema:              jdbc["schema"].(string),
			CredentialsProvider: credentialsProvider,
			Labels:              expandStringList(d.Get("labels").([]interface{})),
			Attributes:          expandAttributes(d),
		}
		return &destination, nil
	}

	if hive, _ := expandSingleMap(d.Get("hive")); hive != nil {
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "hive",
			Database:    hive["database"].(string),
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if bigQuery, _ := expandSingleMap(d.Get("big_query")); bigQuery != nil {
		stagingArea, err := composeGCSStagingArea(bigQuery)
		if err != nil {
			return nil, err
		}

		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "bigquery",
			Path:        bigQuery["path"].(string),
			StagingArea: stagingArea,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if gcs, _ := expandSingleMap(d.Get("gcs")); gcs != nil {
		fileFormat := composeFileFormat(gcs)
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "gcs",
			Bucket:      gcs["bucket"].(string),
			Path:        gcs["path"].(string),
			FileFormat:  fileFormat,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if local, _ := expandSingleMap(d.Get("local")); local != nil {
		fileFormat := composeFileFormat(local)
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "local",
			Path:        local["path"].(string),
			FileFormat:  fileFormat,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if hdfs, _ := expandSingleMap(d.Get("hdfs")); hdfs != nil {
		fileFormat := composeFileFormat(hdfs)
		destination := Destination{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        "hdfs",
			Path:        hdfs["path"].(string),
			FileFormat:  fileFormat,
			Labels:      expandStringList(d.Get("labels").([]interface{})),
			Attributes:  expandAttributes(d),
		}
		return &destination, nil
	}

	if kafka, _ := expandSingleMap(d.Get("kafka")); kafka != nil {
		value := kafka["property"]

		array, ok := kafka["property"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("Kafka Properties Value is not an array. Value: %v", value)
		}

		sensitives := make([]SensitiveAttribute, len(array))
		for i, v := range array {

			prop, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Kafka Properties Value is not a map interfaces. Value: %v.", v)
			}
			sa, err := composeSensitiveAttribute(prop)
			if err != nil {
				return nil, err
			}
			sensitives[i] = *sa
		}

		destination := Destination{
			Name:              d.Get("name").(string),
			Description:       d.Get("description").(string),
			Type:              "kafka",
			BootstrapServers:  kafka["bootstrap_servers"].(string),
			SchemaRegistryURL: kafka["schema_registry_url"].(string),
			KafkaProperties:   sensitives,
			Labels:            expandStringList(d.Get("labels").([]interface{})),
			Attributes:        expandAttributes(d),
		}
		return &destination, nil
	}

	if snowflake, _ := expandSingleMap(d.Get("snowflake")); snowflake != nil {
    	credentialsProviderMap, err := expandSingleMap(snowflake["credentials_provider"])
    	if err != nil {
    		return nil, err
    	}

    	credentialsProvider, err := composeLoginCredentialsProviderConfig(credentialsProviderMap)
    	if err != nil {
    		return nil, err
    	}

    	destination := Destination{
    		Name:                d.Get("name").(string),
    		Description:         d.Get("description").(string),
    		Type:                "snowflake",
    		URL:                 snowflake["url"].(string),
    		Schema:              snowflake["schema"].(string),
    		Warehouse:           snowflake["warehouse"].(string),
    		Database:            snowflake["database"].(string),
    		CredentialsProvider: credentialsProvider,
    		Labels:              expandStringList(d.Get("labels").([]interface{})),
    		Attributes:          expandAttributes(d),
    	}
    	return &destination, nil
    }

	return nil, errors.New("Invalid destination type")
}

func composeGCSStagingArea(d map[string]interface{}) (*GCSStagingArea, error) {
	if temporary, _ := expandSingleMap(d["temporary_staging_area"]); temporary != nil {
		stagingArea := GCSStagingArea{
			Type:   "temporary",
			Bucket: temporary["bucket"].(string),
		}
		return &stagingArea, nil
	}

	if persistent, _ := expandSingleMap(d["persistent_staging_area"]); persistent != nil {
		stagingArea := GCSStagingArea{
			Type:   "persistent",
			Bucket: persistent["bucket"].(string),
			Path:   persistent["path"].(string),
		}
		return &stagingArea, nil
	}

	return nil, errors.New("Invalid staging area type")
}
