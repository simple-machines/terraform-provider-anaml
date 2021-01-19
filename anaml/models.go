package anaml

type Entity struct {
	Id            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultColumn string `json:"defaultColumn"`
}

type TimestampInfo struct {
	Column string `json:"timestampColumn"`
	Zone   string `json:"timezone,omitempty"`
}

type EventDescription struct {
	Entities      map[string]string `json:"entities"`
	TimestampInfo *TimestampInfo    `json:"timestampInfo"`
}

// Go's support for ADTs is so bad we need to use a completely normalised form.
type Table struct {
	Id          int               `json:"id,omitempty"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"adt_type"`
	Sources     []int             `json:"sources"`
	Expression  string            `json:"expression"`
	EventInfo   *EventDescription `json:"eventDescription,omitempty"`
}

type EventWindow struct {
	Type string `json:"adt_type"`
	Days int    `json:"days,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

type SQLExpression struct {
	SQL string `json:"sql"`
}

type AggregateExpression struct {
	Type string `json:"adt_type"`
}

type DataType struct {
	Type string `json:"adt_type"`
}

// Again, completely normalised.
type Feature struct {
	Id          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"adt_type"`
	DataType    DataType             `json:"dataType"`
	Table       int                  `json:"table,omitempty"`
	Window      *EventWindow         `json:"window,omitempty"`
	Select      SQLExpression        `json:"select"`
	Aggregate   *AggregateExpression `json:"aggregate,omitempty"`
	PostExpr    *SQLExpression       `json:"postAggregateExpr,omitempty"`
	Over        []int                `json:"over,omitempty"`
	EntityId    int                  `json:"entityId,omitempty"`
}

// Again, completely normalised.
type FeatureTemplate struct {
	Id           int                   `json:"id,omitempty"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Type         string                `json:"adt_type"`
	DataType     DataType              `json:"dataType"`
	Table        int                   `json:"table"`
	Windows      []EventWindow         `json:"windows"`
	Select       SQLExpression         `json:"select"`
	Aggregations []AggregateExpression `json:"aggregations"`
	PostExpr     SQLExpression         `json:"postAggregateExpr"`
}

type FeatureSet struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EntityId    int    `json:"entity,omitempty"`
	Features    []int  `json:"features,omitempty"`
}

type FeatureStore struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FeatureSet  int    `json:"featureSet"`
	Namespace   string `json:"namespace"`
	TableName   string `json:"tableName,omitempty"`
	Mode        string `json:"mode"`
}
