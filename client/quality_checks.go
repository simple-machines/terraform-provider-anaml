package anaml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandColumnConstraints(info Bag) []ColumnConstraint {
	res := make([]ColumnConstraint, 0, 0)

	extractName := func(single Bag) *string {
		if fetched, ok := single["name"].(string); ok && fetched != "" {
			return &fetched
		}
		return nil
	}

	extractConstraint := func(raw interface{}, constraintType string) ColumnConstraint {
		built := ColumnConstraint{
			Type: constraintType,
		}
		if single, ok := raw.(Bag); ok {
			built.Name = extractName(single)
		}
		return built
	}

	for _, raw := range info["not_null"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_NOT_NULL)
		if grav, ok := raw.(Bag); ok {
			if fetched, ok := grav["threshold"].(float64); ok && fetched != 0.0 {
				built.Threshold = &fetched
			}
		}
		res = append(res, built)
	}

	for _, raw := range info["unique"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_UNIQUE)
		res = append(res, built)
	}

	for _, raw := range info["not_constant"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_NOT_CONSTANT)
		if fetched, ok := raw.(Bag)["enforce_in_partitions"].(bool); ok {
			built.PerPartition = &fetched
		}
		res = append(res, built)
	}

	for _, raw := range info["within_range"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_IN_RANGE)
		if fetched, ok := raw.(Bag)["minimum"].(string); ok && fetched != "" {
			built.Min = &fetched
		}
		if fetched, ok := raw.(Bag)["maximum"].(string); ok && fetched != "" {
			built.Max = &fetched
		}
		if fetched, ok := raw.(Bag)["threshold"].(float64); ok && fetched != 0.0 {
			built.Threshold = &fetched
		}
		res = append(res, built)
	}

	for _, raw := range info["aggregate_within_range"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_STATISTICS_IN_RANGE)
		if fetched, ok := raw.(Bag)["aggregation"].(string); ok && fetched != "" {
			built.Aggregation = &AggregateExpression{
				Type: fetched,
			}
		}
		if fetched, ok := raw.(Bag)["minimum"].(string); ok && fetched != "" {
			built.Min = &fetched
		}
		if fetched, ok := raw.(Bag)["maximum"].(string); ok && fetched != "" {
			built.Max = &fetched
		}
		res = append(res, built)
	}

	for _, raw := range info["row_check"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_ROW_CHECK)
		if fetched, ok := raw.(Bag)["expression"].(string); ok && fetched != "" {
			built.Expression = &SQLExpression{
				SQL: fetched,
			}
		}
		if fetched, ok := raw.(Bag)["threshold"].(float64); ok && fetched != 0.0 {
			built.Threshold = &fetched
		}
		res = append(res, built)
	}

	for _, raw := range info["aggregate_check"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_AGGREGATE_CHECK)
		if fetched, ok := raw.(Bag)["expression"].(string); ok && fetched != "" {
			built.Expression = &SQLExpression{
				SQL: fetched,
			}
		}
		res = append(res, built)
	}

	for _, raw := range info["accepted_values"].([]interface{}) {
		built := extractConstraint(raw, ColumnConstraint_ACCEPTED_VALUES)
		if fetched, ok := raw.(Bag)["values"].(*schema.Set); ok {
			built.Acceptable = expandStringList(fetched.List())
		}
		res = append(res, built)
	}

	return res
}

func flattenColumnConstraints(constraints []ColumnConstraint) Bag {
	notnulls := makeBags(1)
	uniques := makeBags(1)
	notconstants := makeBags(1)
	inranges := makeBags(1)
	agginranges := makeBags(1)
	rowchecks := makeBags(1)
	aggregatechecks := makeBags(1)
	acceptedvalues := makeBags(1)

	for _, constraint := range constraints {
		single := make(Bag)
		if constraint.Name != nil {
			single["name"] = *constraint.Name
		}
		if constraint.Type == ColumnConstraint_NOT_NULL {
			if constraint.Threshold != nil {
				single["threshold"] = constraint.Threshold
			}
			notnulls = append(notnulls, single)
		} else if constraint.Type == ColumnConstraint_UNIQUE {
			uniques = append(uniques, single)
		} else if constraint.Type == ColumnConstraint_NOT_CONSTANT {
			single["enforce_in_partitions"] = constraint.PerPartition
			notconstants = append(notconstants, single)
		} else if constraint.Type == ColumnConstraint_IN_RANGE {
			single["minimum"] = constraint.Min
			single["maximum"] = constraint.Max
			if constraint.Threshold != nil {
				single["threshold"] = constraint.Threshold
			}
			inranges = append(inranges, single)
		} else if constraint.Type == ColumnConstraint_STATISTICS_IN_RANGE {
			single["minimum"] = constraint.Min
			single["maximum"] = constraint.Max
			if constraint.Aggregation != nil {
				single["aggregation"] = constraint.Aggregation.Type
			}
			agginranges = append(agginranges, single)
		} else if constraint.Type == ColumnConstraint_ROW_CHECK {
			if constraint.Expression != nil {
				single["expression"] = constraint.Expression.SQL
			}
			if constraint.Threshold != nil {
				single["threshold"] = constraint.Threshold
			}
			rowchecks = append(rowchecks, single)
		} else if constraint.Type == ColumnConstraint_AGGREGATE_CHECK {
			if constraint.Expression != nil {
				single["expression"] = constraint.Expression.SQL
			}
			aggregatechecks = append(aggregatechecks, single)
		} else if constraint.Type == ColumnConstraint_ACCEPTED_VALUES {
			single["values"] = constraint.Acceptable
			acceptedvalues = append(acceptedvalues, single)
		}
	}

	return Bag{
		"not_null":               notnulls,
		"unique":                 uniques,
		"not_constant":           notconstants,
		"accepted_values":        acceptedvalues,
		"within_range":           inranges,
		"aggregate_within_range": agginranges,
		"row_check":              rowchecks,
		"aggregate_check":        aggregatechecks,
	}
}
