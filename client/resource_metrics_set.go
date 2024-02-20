package anaml

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const metricsSetDescription = `# Metrics Sets

A Metrics Set is collection of metrics and dimensions that are used to perform business level aggregations.
`

func ResourceMetricsSet() *schema.Resource {
	return &schema.Resource{
		Description: metricsSetDescription,
		Create:      resourceMetricsSetCreate,
		Read:        resourceMetricsSetRead,
		Update:      resourceMetricsSetUpdate,
		Delete:      resourceMetricsSetDelete,
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
			"features_source": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     featuresSourceSchema(),
			},
			"tables_source": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				Elem:         tablesSourceSchema(),
				ExactlyOneOf: []string{"tables_source", "features_source"},
			},
			"time_dimension": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     timeDimensionSchema(),
				MaxItems: 1,
			},
			"dimension": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     dimensionSchema(),
			},
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     metricSchema(),
			},
		},
	}
}

func featuresSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"feature_set": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
		},
	}
}

func tablesSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:         schema.TypeString,
				Description:  "The root tables containing measures",
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
			},
			"joins": {
				Type:        schema.TypeList,
				Description: "Dimensions tables to join to.",
				Optional:    true,

				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
		},
	}
}

func dimensionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"expression": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func timeDimensionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"granularity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"daily", "weekly", "monthly", "quarterly",
				}, false),
			},
			"edge": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"from_today", "to_ending",
				}, false),
				Default: "to_ending",
			},
			"back": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func metricSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"select": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "An SQL expression for the column to aggregate.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "An SQL column expression to filter with.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"aggregation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation to perform.",
				ValidateFunc: validation.StringInSlice([]string{
					"sum", "count", "countdistinct", "avg", "std", "min", "max", "minby", "maxby",
					"basketsum", "basketmax", "basketmin", "collectlist", "collectset",
				}, false),
			},
			"post_aggregation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An SQL expression to apply to the result of the feature aggregation.",
			},
		},
	}
}

func resourceMetricsSetRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsSetID := d.Id()

	MetricsSet, err := c.GetMetricsSet(MetricsSetID)
	if err != nil {
		return err
	}
	if MetricsSet == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", MetricsSet.Name); err != nil {
		return err
	}
	if err := d.Set("description", MetricsSet.Description); err != nil {
		return err
	}
	if err := d.Set("labels", MetricsSet.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(MetricsSet.Attributes)); err != nil {
		return err
	}
	if err := readSource(d, MetricsSet.Source); err != nil {
		return err
	}
	if err := readDimensions(d, MetricsSet.Dimensions); err != nil {
		return err
	}
	if err := readMetrics(d, MetricsSet.Metrics); err != nil {
		return err
	}

	return err
}

func resourceMetricsSetCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	MetricsSet, err := buildMetricsSet(d)
	if err != nil {
		return err
	}
	e, err := c.CreateMetricsSet(*MetricsSet)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceMetricsSetUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsSetID := d.Id()

	MetricsSet, err := buildMetricsSet(d)
	if err != nil {
		return err
	}
	err = c.UpdateMetricsSet(MetricsSetID, *MetricsSet)
	if err != nil {
		return err
	}

	return nil
}

func resourceMetricsSetDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	MetricsSetID := d.Id()

	err := c.DeleteMetricsSet(MetricsSetID)
	if err != nil {
		return err
	}

	return nil
}

func buildMetricsSet(d *schema.ResourceData) (*MetricsSet, error) {
	Source, err := buildSource(d)
	if err != nil {
		return nil, err
	}
	Dimensions, err := buildDimensions(d)
	if err != nil {
		return nil, err
	}
	Metrics, err := buildMetrics(d)
	if err != nil {
		return nil, err
	}
	MetricsSet := MetricsSet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      expandLabels(d),
		Attributes:  expandAttributes(d),
		Source:      *Source,
		Dimensions:  Dimensions,
		Metrics:     Metrics,
	}

	return &MetricsSet, nil
}

func buildSource(d *schema.ResourceData) (*MetricsSource, error) {
	features_array := d.Get("features_source").([]interface{})
	tables_array := d.Get("tables_source").([]interface{})

	if len(features_array) > 0 {
		val := features_array[0].(map[string]interface{})
		fs, _ := strconv.Atoi(val["feature_set"].(string))
		res := MetricsSource{
			Type:       "features",
			FeatureSet: &fs,
		}
		return &res, nil
	}

	if len(tables_array) > 0 {
		val := tables_array[0].(map[string]interface{})
		tb, _ := strconv.Atoi(val["table"].(string))
		res := MetricsSource{
			Type:  "table",
			Table: &tb,
			Joins: expandIdentifierList(val["joins"].([]interface{})),
		}
		return &res, nil
	}

	return nil, errors.New("Cluster is null")
}

func buildDimensions(d *schema.ResourceData) ([]Dimension, error) {
	time_array := d.Get("time_dimension").([]interface{})
	user_array := d.Get("dimension").([]interface{})
	res := make([]Dimension, 0, len(time_array)+len(user_array))

	for _, time_dim := range time_array {
		value := time_dim.(map[string]interface{})
		granularity := TypeTag{
			Type: value["granularity"].(string),
		}
		edgeRaw, found := value["edge"].(string)
		edgeTag := "toending"

		if found && edgeRaw == "from_today" {
			edgeTag = "fromtoday"
		}
		edge := TypeTag{
			Type: edgeTag,
		}
		backRaw, found := value["back"].(int)
		var back *int
		if found && edgeRaw == "from_today" {
			back = &backRaw
		}
		Dimension := Dimension{
			Type:        "time",
			Granularity: &granularity,
			Edge:        &edge,
			Back:        back,
		}

		res = append(res, Dimension)
	}

	for _, user_dim := range user_array {
		value := user_dim.(map[string]interface{})
		name := value["name"].(string)
		expression := value["expression"].(string)
		var filter *string
		if filterRaw, ok := value["filter"].(string); ok {
			filter = &filterRaw
		}
		Dimension := Dimension{
			Type:       "user",
			Name:       &name,
			Expression: &expression,
			Filter:     filter,
		}

		res = append(res, Dimension)
	}

	return res, nil
}

func buildMetrics(d *schema.ResourceData) ([]Metric, error) {
	metric_array := d.Get("metric").([]interface{})
	res := make([]Metric, 0, len(metric_array))

	for _, raw := range metric_array {
		value := raw.(map[string]interface{})
		name := value["name"].(string)
		expression := value["select"].(string)
		aggregation := value["aggregation"].(string)
		var filter *SQLExpression
		if filterRaw, ok := value["filter"].(string); ok {
			filter = &SQLExpression{
				SQL: filterRaw,
			}
		}
		var postAgg *SQLExpression
		if postAggRaw, ok := value["post_aggregation"].(string); ok {
			postAgg = &SQLExpression{
				SQL: postAggRaw,
			}
		}
		Metric := Metric{
			Name: &name,
			Select: SQLExpression{
				SQL: expression,
			},
			Aggregate: &AggregateExpression{
				Type: aggregation,
			},
			Filter:      filter,
			PostAggExpr: postAgg,
		}

		res = append(res, Metric)
	}

	return res, nil
}

func readSource(d *schema.ResourceData, source MetricsSource) error {
	empty := make([]interface{}, 0, 0)
	if source.Type == "features" {
		res := make([]interface{}, 0, 1)
		single := make(map[string]interface{})
		single["feature_set"] = strconv.Itoa(*source.FeatureSet)
		res = append(res, single)

		if err := d.Set("features_source", res); err != nil {
			return err
		}
		if err := d.Set("tables_source", empty); err != nil {
			return err
		}

		return nil
	}

	if source.Type == "table" {
		res := make([]interface{}, 0, 1)
		single := make(map[string]interface{})
		single["table"] = strconv.Itoa(*source.Table)
		single["joins"] = identifierList(source.Joins)
		res = append(res, single)

		if err := d.Set("features_source", empty); err != nil {
			return err
		}
		if err := d.Set("tables_source", res); err != nil {
			return err
		}

		return nil
	}

	return errors.New("Unrecognised Metrics Source tag")
}

func readDimensions(d *schema.ResourceData, dimensions []Dimension) error {
	user := make([]interface{}, 0, len(dimensions))
	time := make([]interface{}, 0, len(dimensions))

	for _, dimension := range dimensions {
		if dimension.Type == "user" {
			single := make(map[string]interface{})
			single["name"] = dimension.Name
			single["expression"] = dimension.Expression
			single["filter"] = dimension.Filter
			user = append(user, single)
		} else if dimension.Type == "time" {
			single := make(map[string]interface{})
			single["granularity"] = dimension.Granularity.Type
			if dimension.Edge.Type == "fromtoday" {
				single["edge"] = "from_today"
			} else {
				single["edge"] = "to_ending"
			}
			single["back"] = dimension.Back
			time = append(time, single)
		} else {
			return errors.New("Unrecognised Dimension tag")
		}
	}

	if err := d.Set("time_dimension", time); err != nil {
		return err
	}
	if err := d.Set("dimension", user); err != nil {
		return err
	}

	return nil
}

func readMetrics(d *schema.ResourceData, input []Metric) error {
	metrics := make([]interface{}, 0, len(input))

	for _, metric := range input {
		single := make(map[string]interface{})
		single["name"] = metric.Name
		single["select"] = metric.Select.SQL
		single["aggregation"] = metric.Aggregate.Type
		single["filter"] = nil
		if metric.Filter != nil {
			single["filter"] = metric.Filter.SQL
		}
		single["post_aggregation"] = nil
		if metric.PostAggExpr != nil {
			single["post_aggregation"] = metric.PostAggExpr.SQL
		}

		metrics = append(metrics, single)
	}

	if err := d.Set("metric", metrics); err != nil {
		return err
	}

	return nil
}
