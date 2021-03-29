package anaml

// Entity ..
type Entity struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultColumn string `json:"defaultColumn"`
}

// EntityMapping ..
type EntityMapping struct {
	ID      int `json:"id,omitempty"`
	From    int `json:"from"`
	To      int `json:"to"`
	Mapping int `json:"mapping"`
}

// TimestampInfo ..
type TimestampInfo struct {
	Column string `json:"timestampColumn"`
	Zone   string `json:"timezone,omitempty"`
}

// EventDescription ..
type EventDescription struct {
	Entities      map[string]string `json:"entities"`
	TimestampInfo *TimestampInfo    `json:"timestampInfo"`
}

// Table ...
type Table struct {
	ID            int               `json:"id,omitempty"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Type          string            `json:"adt_type"`
	Sources       []int             `json:"sources,omitempty"`
	Source        *SourceReference  `json:"source,omitempty"`
	Expression    string            `json:"expression,omitempty"`
	EventInfo     *EventDescription `json:"eventDescription,omitempty"`
	EntityMapping int               `json:"entityMapping,omitempty"`
	ExtraFeatures []int             `json:"extraFeatures,omitempty"`
}

// EventWindow ...
type EventWindow struct {
	Type string `json:"adt_type"`
	Days int    `json:"days,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

// SQLExpression ...
type SQLExpression struct {
	SQL string `json:"sql"`
}

// AggregateExpression ...
type AggregateExpression struct {
	Type string `json:"adt_type"`
}

// DataType ...
type DataType struct {
	Type string `json:"adt_type"`
}

// Feature ... again, completely normalised.
type Feature struct {
	ID          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"adt_type"`
	DataType    DataType             `json:"dataType"`
	Table       int                  `json:"table,omitempty"`
	Window      *EventWindow         `json:"window,omitempty"`
	Select      SQLExpression        `json:"select"`
	Filter      *SQLExpression       `json:"filter"`
	Aggregate   *AggregateExpression `json:"aggregate,omitempty"`
	PostExpr    *SQLExpression       `json:"postAggregateExpr,omitempty"`
	Over        []int                `json:"over,omitempty"`
	EntityID    int                  `json:"entityId,omitempty"`
	TemplateID  *int                 `json:"template,omitempty"`
}

// FeatureTemplate ... Again, completely normalised.
type FeatureTemplate struct {
	ID          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"adt_type"`
	DataType    DataType             `json:"dataType"`
	Table       int                  `json:"table"`
	Select      SQLExpression        `json:"select"`
	Filter      *SQLExpression       `json:"filter"`
	Aggregate   *AggregateExpression `json:"aggregate,omitempty"`
	PostExpr    SQLExpression        `json:"postAggregateExpr"`
}

// FeatureSet ...
type FeatureSet struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EntityID    int    `json:"entity,omitempty"`
	Features    []int  `json:"features"`
}

// FeatureStore ...
type FeatureStore struct {
	ID           int                    `json:"id,omitempty"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	FeatureSet   int                    `json:"featureSet"`
	Enabled      bool                   `json:"enabled"`
	Schedule     *Schedule              `json:"schedule"`
	Destinations []DestinationReference `json:"destinations"`
	Cluster      int                    `json:"cluster"`
}

type Schedule struct {
	Type           string       `json:"adt_type"`
	StartTimeOfDay *string      `json:"startTimeOfDay,omitempty"`
	CronString     string       `json:"cronString,omitempty"`
	RetryPolicy    *RetryPolicy `json:"retryPolicy,omitempty"`
}

type RetryPolicy struct {
	Type        string `json:"adt_type"`
	Backoff     string `json:"backoff,omitempty"`
	MaxAttempts int    `json:"maxAttempts,omitempty"`
}

// Source ...
type Source struct {
	ID                  int                             `json:"id,omitempty"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	Type                string                          `json:"adt_type"`
	Bucket              string                          `json:"bucket,omitempty"`
	Path                string                          `json:"path,omitempty"`
	FileFormat          *FileFormat                     `json:"fileFormat,omitempty"`
	Endpoint            string                          `json:"endpoint,omitempty"`
	AccessKey           string                          `json:"accessKey,omitempty"`
	SecretKey           string                          `json:"secretKey,omitempty"`
	URL                 string                          `json:"url,omitempty"`
	Schema              string                          `json:"schema,omitempty"`
	CredentialsProvider *LoginCredentialsProviderConfig `json:"credentialsProvider,omitempty"`
	Database            string                          `json:"database,omitempty"`
	BootstrapServers    string                          `json:"bootstrapServers,omitempty"`
	SchemaRegistryURL   string                          `json:"schemaRegistryUrl,omitempty"`
}

type FileFormat struct {
	Type          string `json:"adt_type"`
	IncludeHeader *bool  `json:"includeHeader,omitempty"`
}

// SourceReference ...
type SourceReference struct {
	SourceID  int    `json:"sourceId"`
	Folder    string `json:"folder,omitempty"`
	TableName string `json:"tableName,omitempty"`
}

// Destination ...
type Destination struct {
	ID                  int                             `json:"id,omitempty"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	Type                string                          `json:"adt_type"`
	Bucket              string                          `json:"bucket,omitempty"`
	Path                string                          `json:"path,omitempty"`
	FileFormat          *FileFormat                     `json:"fileFormat,omitempty"`
	Endpoint            string                          `json:"endpoint,omitempty"`
	AccessKey           string                          `json:"accessKey,omitempty"`
	SecretKey           string                          `json:"secretKey,omitempty"`
	URL                 string                          `json:"url,omitempty"`
	Schema              string                          `json:"schema,omitempty"`
	CredentialsProvider *LoginCredentialsProviderConfig `json:"credentialsProvider,omitempty"`
	Database            string                          `json:"database,omitempty"`
	BootstrapServers    string                          `json:"bootstrapServers,omitempty"`
	SchemaRegistryURL   string                          `json:"schemaRegistryUrl,omitempty"`
	StagingArea         *GCSStagingArea                 `json:"stagingArea,omitempty"`
}

// GCSStagingArea ...
type GCSStagingArea struct {
	Type   string `json:"adt_type"`
	Bucket string `json:"bucket"`
	Path   string `json:"path,omitempty"`
}

// DestinationReference ...
type DestinationReference struct {
	DestinationID int    `json:"destinationId"`
	Folder        string `json:"folder,omitempty"`
	TableName     string `json:"tableName,omitempty"`
}

// Cluster ...
type Cluster struct {
	ID                  int                     `json:"id,omitempty"`
	Name                string                  `json:"name"`
	Description         string                  `json:"description"`
	Type                string                  `json:"adt_type"`
	IsPreviewCluster    bool                    `json:"isPreviewCluster"`
	AnamlServerURL      string                  `json:"anamlServerUrl,omitempty"`
	SparkServerURL      string                  `json:"sparkServerUrl,omitempty"`
	CredentialsProvider *JWTTokenProviderConfig `json:"credentialsProvider,omitempty"`
	SparkConfig         *SparkConfig            `json:"sparkConfig,omitempty"`
}

// JWTTokenProviderConfig ...
type JWTTokenProviderConfig struct {
	Type                           string                          `json:"adt_type"`
	LoginServerURL                 string                          `json:"loginServerUrl"`
	LoginCredentialsProviderConfig *LoginCredentialsProviderConfig `json:"loginCredentialsProviderConfig"`
}

// LoginCredentialsProviderConfig  ...
type LoginCredentialsProviderConfig struct {
	Type                  string `json:"adt_type"`
	Username              string `json:"username"`
	Password              string `json:"password,omitempty"`
	PasswordSecretProject string `json:"passwordSecretProject,omitempty"`
	PasswordSecretId      string `json:"passwordSecretId,omitempty"`
}

// SparkConfig ...
type SparkConfig struct {
	EnableHiveSupport         bool              `json:"enableHiveSupport"`
	HiveMetastoreURL          string            `json:"hiveMetastoreUrl,omitempty"`
	AdditionalSparkProperties map[string]string `json:"additionalSparkProperties"`
}

// User ..
type User struct {
	ID        int      `json:"id,omitempty"`
	Name      string   `json:"name"`
	Email     *string  `json:"email,omitempty"`
	GivenName *string  `json:"givenName,omitempty"`
	Surname   *string  `json:"surname,omitempty"`
	Password  *string  `json:"password,omitempty"`
	Roles     []string `json:"roles"`
}

type ChangeOtherPasswordRequest struct {
	Password string `json:"password"`
}
