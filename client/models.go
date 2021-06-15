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
// Note
// Go is a bad language, We can't use omitempty for over, because both [] and 'nil' are empty.
// Empty list is appropriate, especially for templates. But unfortunately, we will be sending
// a really dumb `null` where it doesn't make sense to do so.
type Feature struct {
	ID          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"adt_type"`
	Table       int                  `json:"table,omitempty"`
	Window      *EventWindow         `json:"window,omitempty"`
	Select      SQLExpression        `json:"select"`
	Filter      *SQLExpression       `json:"filter"`
	Aggregate   *AggregateExpression `json:"aggregate,omitempty"`
	PostExpr    *SQLExpression       `json:"postAggregateExpr,omitempty"`
	Over        []int                `json:"over"`
	EntityID    int                  `json:"entityId,omitempty"`
	TemplateID  *int                 `json:"template,omitempty"`
}

// FeatureTemplate ... again, completely normalised.
type FeatureTemplate struct {
	ID          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"adt_type"`
	Table       int                  `json:"table"`
	Window      *EventWindow         `json:"window,omitempty"`
	Select      SQLExpression        `json:"select"`
	Filter      *SQLExpression       `json:"filter"`
	Aggregate   *AggregateExpression `json:"aggregate,omitempty"`
	PostExpr    SQLExpression        `json:"postAggregateExpr"`
	Over        []int                `json:"over"`
	EntityID    int                  `json:"entityId,omitempty"`
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
	StartDate    *string                `json:"startDate,omitempty"`
	EndDate      *string                `json:"endDate,omitempty"`
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
	ID                  int                             `json:"id,omitempty"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	Type                string                          `json:"adt_type"`
	IsPreviewCluster    bool                            `json:"isPreviewCluster"`
	AnamlServerURL      string                          `json:"anamlServerUrl,omitempty"`
	SparkServerURL      string                          `json:"sparkServerUrl,omitempty"`
	CredentialsProvider *LoginCredentialsProviderConfig `json:"credentialsProvider,omitempty"`
	SparkConfig         *SparkConfig                    `json:"sparkConfig,omitempty"`
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

// UserGroup ..
type UserGroup struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Members     []int   `json:"members"`
}

// BranchProtection
type BranchProtection struct {
	ID                  int            `json:"id,omitempty"`
	ProtectionPattern   string         `json:"protectionPattern"`
	MergeApprovalRules  []ApprovalRule `json:"mergeApprovalRules"`
	PushWhitelist       []PrincipalId  `json:"pushWhitelist"`
	ApplyToAdmins       bool           `json:"applyToAdmins"`
	AllowBranchDeletion bool           `json:"allowBranchDeletion"`
}

// ApprovalRule
type ApprovalRule struct {
	Approvers            []PrincipalId `json:"approvers,omitempty"`
	NumRequiredApprovals int           `json:"numRequiredApprovals"`
	Type                 string        `json:"adt_type"`
}

// PrincipalId
type PrincipalId struct {
	ID   int    `json:"id"`
	Type string `json:"adt_type"`
}

// TableMonitoring ...
type TableMonitoring struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tables      []int     `json:"tables"`
	Schedule    *Schedule `json:"schedule"`
	Cluster     int       `json:"cluster"`
	Enabled     bool      `json:"enabled"`
}

// TableCaching ...
type TableCaching struct {
	ID          int                `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Specs       []TableCachingSpec `json:"specs"`
  PrefixURI   string             `json:"prefixURI"`
	Schedule    *Schedule          `json:"schedule"`
	Cluster     int                `json:"cluster"`
}

type TableCachingSpec struct {
	Table   int  `json:"table"`
	Entity  int  `json:"entity"`
}
