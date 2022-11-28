package anaml

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const eventStoreDescription = `# Event Stores`

func ResourceEventStore() *schema.Resource {
	return &schema.Resource{
		Description: eventStoreDescription,
		Create:      resourceEventStoreCreate,
		Read:        resourceEventStoreRead,
		Update:      resourceEventStoreUpdate,
		Delete:      resourceEventStoreDelete,
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
				Optional: true,
			},
			"bootstrap_servers": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"schema_registry_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"property": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     sensitiveAttributeSchema(),
			},
			"ingestion": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     ingestionSchema(),
			},
			"connect_base_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"batch_ingest_base_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				AtLeastOneOf: []string{"connect_base_uri", "batch_ingest_base_uri"},
			},
			"scatter_base_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"glacier_base_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"daily_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          dailyScheduleSchema(),
				ConflictsWith: []string{"cron_schedule"},
			},
			"cron_schedule": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Elem:          cronScheduleSchema(),
				ConflictsWith: []string{"daily_schedule"},
			},
			"cluster": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"access_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Access rules to attach to the object",
				Elem:        accessRuleSchema(),
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Labels to attach to the object",
				Elem:        labelSchema(),
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

func ingestionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic": {
				Type:     schema.TypeString,
				Required: true,
			},
			"streaming": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"entity_column": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timestamp_column": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEventStoreRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entityID := d.Id()

	entity, err := c.GetEventStore(entityID)
	if err != nil {
		return err
	}

	if entity == nil {
		d.SetId("")
		return nil
	}

	sensitives := make([]map[string]interface{}, len(entity.KafkaProperties))
	for i, v := range entity.KafkaProperties {
		sa, err := parseSensitiveAttribute(&v)
		if err != nil {
			return err
		}

		sensitives[i] = sa
	}

	ingestions := make([]map[string]interface{}, 0)
	for k, v := range entity.Ingestions {
		ingestion := make(map[string]interface{})
		ingestion["topic"] = k
		ingestion["streaming"] = v.HasStreaming
		ingestion["entity_column"] = v.Entity
		ingestion["timestamp_column"] = v.TimestampInfo.Column
		ingestion["timezone"] = v.TimestampInfo.Zone
		ingestions = append(ingestions, ingestion)
	}

	if err := d.Set("name", entity.Name); err != nil {
		return err
	}
	if err := d.Set("description", entity.Description); err != nil {
		return err
	}
	if err := d.Set("labels", entity.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(entity.Attributes)); err != nil {
		return err
	}
	if err := d.Set("bootstrap_servers", entity.BootstrapServers); err != nil {
		return err
	}
	if err := d.Set("schema_registry_url", entity.SchemaRegistryURL); err != nil {
		return err
	}
	if err := d.Set("property", sensitives); err != nil {
		return err
	}
	if err := d.Set("ingestion", ingestions); err != nil {
		return err
	}
	if err := d.Set("connect_base_uri", entity.ConnectBaseURI); err != nil {
		return err
	}
	if err := d.Set("batch_ingest_base_uri", entity.BatchIngestBaseURI); err != nil {
		return err
	}
	if err := d.Set("scatter_base_uri", entity.ScatterBaseURI); err != nil {
		return err
	}
	if err := d.Set("glacier_base_uri", entity.GlacierBaseURI); err != nil {
		return err
	}
	if err := d.Set("cluster", strconv.Itoa(entity.Cluster)); err != nil {
		return err
	}
	if err := d.Set("access_rules", flattenAccessRules(entity.AccessRules)); err != nil {
		return err
	}
	if entity.Schedule.Type == "daily" {
		dailySchedules, err := parseDailySchedule(entity.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("daily_schedule", dailySchedules); err != nil {
			return err
		}
		if err := d.Set("cron_schedule", nil); err != nil {
			return err
		}
	} else if entity.Schedule.Type == "cron" {
		cronSchedules, err := parseCronSchedule(entity.Schedule)
		if err != nil {
			return err
		}
		if err := d.Set("cron_schedule", cronSchedules); err != nil {
			return err
		}
		if err := d.Set("daily_schedule", nil); err != nil {
			return err
		}
	} else {
		if err := d.Set("cron_schedule", nil); err != nil {
			return err
		}
		if err := d.Set("daily_schedule", nil); err != nil {
			return err
		}
	}
	return err
}

func resourceEventStoreCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	eventStore, err := buildEventStore(d)
	if err != nil {
		return err
	}
	e, err := c.CreateEventStore(*eventStore)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceEventStoreUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	eventStoreID := d.Id()
	eventStore, err := buildEventStore(d)
	if err != nil {
		return err
	}

	err = c.UpdateEventStore(eventStoreID, *eventStore)
	if err != nil {
		return err
	}

	return nil
}

func resourceEventStoreDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	entityID := d.Id()

	err := c.DeleteEventStore(entityID)
	if err != nil {
		return err
	}

	return nil
}

func buildEventStore(d *schema.ResourceData) (*EventStore, error) {
	accessRules, err := expandAccessRules(d.Get("access_rules").([]interface{}))
	if err != nil {
		return nil, err
	}

	array, ok := d.Get("property").([]interface{})
	if !ok {
		return nil, fmt.Errorf("Kafka Properties Value is not an array.")
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

	ingestions := make(map[string]EventStoreTopicColumns)
	rawIngest := d.Get("ingestion").([]interface{})
	for _, v := range rawIngest {
		vv := v.(map[string]interface{})
		hasStreaming := vv["streaming"].(bool)

		ingestions[vv["topic"].(string)] = EventStoreTopicColumns{
			Entity: vv["entity_column"].(string),
			TimestampInfo: &TimestampInfo{
				Column: vv["timestamp_column"].(string),
				Zone:   vv["timezone"].(string),
			},
			HasStreaming: hasStreaming,
		}
	}
	cluster, err := strconv.Atoi(d.Get("cluster").(string))
	if err != nil {
		return nil, err
	}
	schedule := composeNeverSchedule()
	if dailySchedule, _ := expandSingleMap(d.Get("daily_schedule")); dailySchedule != nil {
		schedule, err = composeDailySchedule(dailySchedule)
		if err != nil {
			return nil, err
		}
	}
	if cronSchedule, _ := expandSingleMap(d.Get("cron_schedule")); cronSchedule != nil {
		schedule, err = composeCronSchedule(cronSchedule)
		if err != nil {
			return nil, err
		}
	}

	entity := EventStore{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		BootstrapServers:   d.Get("bootstrap_servers").(string),
		SchemaRegistryURL:  d.Get("schema_registry_url").(string),
		KafkaProperties:    sensitives,
		Ingestions:         ingestions,
		ConnectBaseURI:     getNullableString(d, "connect_base_uri"),
		BatchIngestBaseURI: getNullableString(d, "batch_ingest_base_uri"),
		ScatterBaseURI:     d.Get("scatter_base_uri").(string),
		GlacierBaseURI:     d.Get("glacier_base_uri").(string),
		Labels:             expandLabels(d),
		Attributes:         expandAttributes(d),
		Cluster:            cluster,
		Schedule:           schedule,
		AccessRules:        accessRules,
	}
	return &entity, err
}
