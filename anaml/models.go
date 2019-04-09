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
	Id            int            `json:"entityId"`
	Column        string         `json:"keyColumn"`
	TimestampInfo *TimestampInfo `json:"timestampInfo"`
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
	Id          int                 `json:"id,omitempty"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"adt_type"`
	DataType    DataType            `json:"dataType"`
	Table       int                 `json:"table"`
	Window      EventWindow         `json:"window"`
	Select      SQLExpression       `json:"select"`
	Aggregate   AggregateExpression `json:"aggregate"`
	// PostExpr    SQLExpression       `json:"postAggregateExpr"`
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
