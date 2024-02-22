terraform {
  required_version = ">= 0.14"
  required_providers {
    anaml = {
      source = "registry.anaml.io/anaml/anaml"
    }
    anaml-operations = {
      source = "registry.anaml.io/anaml/anaml-operations"
    }
  }
}

provider "anaml" {
  host     = "http://localhost:8080/api"
  username = "03d147fe-0fa8-4aef-bce6-e6fbcd1cd000"
  password = "test secret"
  branch   = "official"
}

provider "anaml-operations" {
  host     = "http://127.0.0.1:8080/api"
  username = "03d147fe-0fa8-4aef-bce6-e6fbcd1cd000"
  password = "test secret"
}


data "anaml-operations_source" "s3a" {
  name = anaml-operations_source.s3a.name
}

data "anaml-operations_cluster" "local" {
  name = anaml-operations_cluster.local.name
}

data "anaml-operations_destination" "s3a" {
  name = anaml-operations_destination.s3a.name
}

data "anaml-operations_destination" "kafka" {
  name = anaml-operations_destination.kafka.name
}

data "anaml-operations_destination" "online" {
  name = anaml-operations_destination.online.name
}

resource "anaml-operations_attribute_restriction" "country" {
  key = "terraform_country"
  description = "Applicable country for terraformed resources"
  enum {
    choice {
      value          = "australia"
      display_emoji  = "🇦🇺"
      display_colour = "#00008B"
    }
    choice { value = "uk" }
    choice { value = "america" }
  }
  applies_to = ["cluster", "destination", "entity", "feature", "feature_set", "feature_store", "feature_template", "source", "table"]
}

resource "anaml-operations_label_restriction" "terraform" {
  text   = "terraformed resource"
  emoji  = "🛠️"
  colour = "#B0B0B0"
}

resource "anaml-operations_label_restriction" "important" {
  text = "important"
}

resource "anaml_entity" "household" {
  name           = "household"
  description    = "A household level view"
  default_column = "household"

  labels = [ anaml-operations_label_restriction.terraform.text, anaml-operations_label_restriction.important.text ]
  attribute {
    key   = anaml-operations_attribute_restriction.country.key
    value = "australia"
  }
}

resource "anaml_table" "household" {
  name        = "household"
  description = "A household level view"

  source {
    source = data.anaml-operations_source.s3a.id
    folder = "household"
  }

  event {
    entities = {
      (anaml_entity.household.id) = "household_id"
    }
    timestamp_column = "timestamp"
  }

  labels = [ anaml-operations_label_restriction.terraform.text, anaml-operations_label_restriction.important.text ]
}

resource "anaml_entity_population" "adults" {
  name        = "adults"
  description = "Adults in the household"

  entity     = anaml_entity.household.id
  labels     = [ anaml-operations_label_restriction.terraform.text ]
  expression = "SELECT customer, daily() FROM household WHERE AGE > 18"
  sources    = [anaml_table.household.id]
}

resource "anaml_table" "household_normalised" {
  name        = "household_normalised"
  description = "A household level view"

  view {
    sources    = [anaml_table.household.id]
    expression = "SELECT * FROM household"
  }

  event {
    entities = {
      (anaml_entity.household.id) = "household"
    }
    timestamp_column = "timestamp"
    timezone = "Australia/Brisbane"
  }

  labels = [ anaml-operations_label_restriction.terraform.text, anaml-operations_label_restriction.important.text ]
}

resource "anaml_feature_template" "household_count" {
  name                = "household_count"
  description         = "Count of household items"
  table               = anaml_table.household.id
  select              = "count"
  aggregation         = "sum"
  entity_restrictions = [anaml_entity.household.id]

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml_feature" "household_count" {
  for_each = toset(["1", "2", "4"])
  days     = parseint(each.key, 10)

  name                = "household_count_${each.key}_days"
  description         = "Count of household items"
  table               = anaml_table.household.id
  select              = "count"
  aggregation         = "sum"
  template            = anaml_feature_template.household_count.id
  entity_restrictions = anaml_feature_template.household_count.entity_restrictions

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml_feature" "household_count_without_entity_restrictions" {
  for_each = toset(["1", "2", "4"])
  days     = parseint(each.key, 10)

  name        = "household_count_without_er_${each.key}_days"
  description = "Count of household items"
  table       = anaml_table.household.id
  select      = "count"
  aggregation = "sum"

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml_feature_set" "household" {
  name   = "household"
  entity = anaml_entity.household.id
  features = [
    anaml_feature.household_count["1"].id
    , anaml_feature.household_count["2"].id
    , anaml_feature.household_count["4"].id
  ]

  labels = [ anaml-operations_label_restriction.terraform.text, anaml-operations_label_restriction.important.text ]
}

resource "anaml-operations_event_store" "basic" {
  name           = "household"
  description    = "A household level view"
  bootstrap_servers   = "http://bootstrap"
  schema_registry_url = "http://schema-registry"
  property {
    key = "jamf"
    gcp {
      secret_project = "example"
      secret_id      = "sid"
    }
  }
  ingestion {
    topic = "topic"
    entity_column = "what"
    timestamp_column = "when"
    timezone = "Australia/Sydney"
  }

  connect_base_uri = "connect"
  scatter_base_uri = "scatter"
  glacier_base_uri = "glacier"

  labels = [ anaml-operations_label_restriction.terraform.text ]
  attribute {
    key   = anaml-operations_attribute_restriction.country.key
    value = "australia"
  }
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
}

resource "anaml-operations_feature_store" "household_daily" {
  name        = "household_daily"
  description = "Daily view of households"
  start_date  = "2020-01-01"
  end_date    = "2021-01-01"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.s3a.id
    folder {
      path = "household_results"
      partitioning_enabled = true
      save_mode = "overwrite"
    }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_daily_table_dest" {
  name        = "household_daily_table_dest"
  description = "Daily view of households"
  start_date  = "2020-01-01"
  end_date    = "2021-01-01"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.online.id
    table {
      name = "household_results"
    }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_daily_table_spark_prop" {
  name        = "household_daily_table_spark_properties"
  description = "Daily view of households"
  start_date  = "2020-01-01"
  end_date    = "2021-01-01"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  additional_spark_properties = {
    "spark.driver.extraClassPath" : "/opt/docker/lib/*"
  }
  destination {
    destination = data.anaml-operations_destination.online.id
    table {
      name = "household_results"
    }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_daily_topic_dest" {
  name        = "household_daily_topic_dest"
  description = "Daily view of households"
  start_date  = "2020-01-01"
  end_date    = "2021-01-01"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.kafka.id
    topic {
      name = "household_results"
      format = "json"
    }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_cron" {
  name              = "household_cron"
  description       = "Daily view of households"
  feature_set       = anaml_feature_set.household.id
  enabled           = true
  cluster           = data.anaml-operations_cluster.local.id
  entity_population = anaml_entity_population.adults.id
  destination {
    destination                 = data.anaml-operations_destination.s3a.id
    folder {
      path = "household_results"
      partitioning_enabled = true
      save_mode = "append"
    }
  }
  cron_schedule {
    cron_string = "* * * * *"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_never" {
  name        = "household_never"
  description = "Manually scheduled view of households"
  start_date  = "2020-01-01"
  end_date    = "2021-01-01"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.s3a.id
    folder {
      path = "household_results"
      partitioning_enabled = true
      save_mode = "ignore"
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_daily_retry" {
  name        = "household_daily_retry"
  description = "Daily view of households"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.s3a.id
    folder {
      path = "household_results"
      partitioning_enabled = true
      save_mode = "errorifexists"
    }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"

    fixed_retry_policy {
      backoff      = "PT1H30M"
      max_attempts = 3
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_feature_store" "household_cron_retry" {
  name        = "household_cron_retry"
  description = "Daily view of households"
  feature_set = anaml_feature_set.household.id
  enabled     = true
  cluster     = data.anaml-operations_cluster.local.id
  destination {
    destination                 = data.anaml-operations_destination.s3a.id
    folder {
      path = "household_results"
      partitioning_enabled = true
      save_mode = "overwrite"
    }
  }
  cron_schedule {
    cron_string = "* * * * *"

    fixed_retry_policy {
      backoff      = "PT1H30M"
      max_attempts = 3
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_cluster" "local" {
  name               = "terraform_local_cluster"
  description        = "A local cluster created by Terraform"
  is_preview_cluster = true

  local {
    anaml_server_url = "http://localhost:8080"
    basic {
      username = "admin"
      password = "test password"
    }
  }

  spark_config {
    enable_hive_support = true
  }

  property_set {
    name = "small"
    additional_spark_properties = {"spark.dynamicAllocation.maxExecutors": "2"}
  }

  property_set {
      name = "medium"
      additional_spark_properties = {"spark.dynamicAllocation.maxExecutors": "4"}
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_cluster" "spark_server" {
  name               = "terraform_spark_server_cluster"
  description        = "A Spark server cluster created by Terraform"
  is_preview_cluster = false

  spark_server {
    spark_server_url = "http://localhost:8080"
  }

  spark_config {
    enable_hive_support = true
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "s3" {
  name        = "terraform_s3_source"
  description = "An S3 source created by Terraform"

  s3 {
    bucket = "my-bucket"
    path   = "/path/to/file"

    file_format                = "csv"
    compression                = "gzip"
    include_header             = true
    field_separator            = ","
    quote_all                  = true
    date_format                = "yyyy-MM-dd"
    timestamp_format           = "yyyy-MM-dd HH:MM:SS"
    ignore_leading_whitespace  = false
    ignore_trailing_whitespace = false
  }

  access_rule {
    resource = "customers"

    principals {
      user_group {
        id = anaml-operations_user_group.engineering.id
      }
    }

    masking_rule {
      filter {
        expression = "id % 2 = 0"
      }
    }
    masking_rule {
      mask {
        column     = "email"
        expression = "x -> NULL"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "s3a" {
  name        = "terraform_s3a_source"
  description = "An S3A source created by Terraform"

  s3a {
    bucket      = "my-bucket"
    path        = "/path/to/file"
    endpoint    = "http://example.com"
    file_format = "orc"
    access_key  = "access"
    secret_key  = "secret"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "hive" {
  name        = "terraform_hive_source"
  description = "An Hive source created by Terraform"

  hive {
    database = "my_database"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "jdbc" {
  name        = "terraform_jdbc_source"
  description = "An JDBC source created by Terraform"

  jdbc {
    url    = "jdbc://my/database"
    schema = "my_schema"

    credentials_provider {
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "big_query" {
  name        = "terraform_bigquery_source"
  description = "An BigQuery source created by Terraform"

  big_query {
    path = "/path/to/file"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "gcs" {
  name        = "terraform_gcs_source"
  description = "An GCS source created by Terraform"

  gcs {
    bucket      = "my-bucket"
    path        = "/path/to/file"
    file_format = "parquet"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "local" {
  name        = "terraform_local_source"
  description = "An Local source created by Terraform"

  local {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "hdfs" {
  name        = "terraform_hdfs_source"
  description = "An HDFS source created by Terraform"

  hdfs {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "kafka" {
  name        = "terraform_kafka_source"
  description = "An Kafka source created by Terraform"

  kafka {
    bootstrap_servers   = "http://bootstrap"
    schema_registry_url = "http://schema-registry"
    property {
      key = "clock"
      value = "time"
    }
    property {
      key = "quest"
      aws {
        secret_id      = "sid"
      }
    }
    property {
      key = "jamf"
      gcp {
        secret_project = "example"
        secret_id      = "sid"
      }
    }
    property {
      key = "exhibit"
      file {
        filepath = "example"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_source" "snowflake" {
  name        = "terraform_snowflake_source"
  description = "An Snowflake source created by Terraform"

  snowflake {
    url    = "snowflake://my/database"
    schema = "my_schema"
    database = "my_database"
    warehouse = "my_warehouse"

    credentials_provider {
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "s3" {
  name        = "terraform_s3_destination"
  description = "An S3 destination created by Terraform"

  s3 {
    bucket         = "my-bucket"
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = true
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "s3a" {
  name        = "terraform_s3a_destination"
  description = "An S3A destination created by Terraform"

  s3a {
    bucket      = "my-bucket"
    path        = "/path/to/file"
    endpoint    = "http://example.com"
    file_format = "orc"
    access_key  = "access"
    secret_key  = "secret"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "hive" {
  name        = "terraform_hive_destination"
  description = "An Hive destination created by Terraform"

  hive {
    database = "my_database"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "jdbc" {
  name        = "terraform_jdbc_destination"
  description = "An JDBC destination created by Terraform"

  jdbc {
    url    = "jdbc://my/database"
    schema = "my_schema"

    credentials_provider {
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "big_query_temporary" {
  name        = "terraform_bigquery_destination_with_temporary_staging_area"
  description = "An BigQuery destination created by Terraform"

  big_query {
    path = "path/to/file"
    temporary_staging_area {
      bucket = "my-bucket"
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "big_query_persistent" {
  name        = "terraform_bigquery_destination_with_persistent_staging_area"
  description = "An BigQuery destination created by Terraform"

  big_query {
    path = "path/to/file"
    persistent_staging_area {
      bucket = "my-bucket"
      path   = "/path/to/file"
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "gcs" {
  name        = "terraform_gcs_destination"
  description = "An GCS destination created by Terraform"

  gcs {
    bucket      = "my-bucket"
    path        = "/path/to/file"
    file_format = "parquet"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "local" {
  name        = "terraform_local_destination"
  description = "An Local destination created by Terraform"

  local {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "hdfs" {
  name        = "terraform_hdfs_destination"
  description = "An HDFS destination created by Terraform"

  hdfs {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "kafka" {
  name        = "terraform_kafka_destination"
  description = "An Kafka destination created by Terraform"

  kafka {
    bootstrap_servers   = "http://bootstrap"
    schema_registry_url = "http://schema-registry"
    property {
      key   = "username"
      value = "fred"
    }
    property {
      key = "password"
      aws {
        secret_id = "secret_number_3"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "online" {
  name        = "terraform_online_feature_store_destination"
  description = "An Online feature store destination created by Terraform"

  online {
    url    = "jdbc://localhost:5454"
    schema = "my_schema"

    credentials_provider {
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "snowflake" {
  name        = "terraform_snowflake_destination"
  description = "An Snowflake destination created by Terraform"

  snowflake {
    url    = "snowflake://my/database"
    schema = "my_schema"
    database = "my_database"
    warehouse = "my_warehouse"

    credentials_provider {
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_destination" "bigtable" {
  name        = "terraform_bigtable_destination"
  description = "A Bigtable destination created by Terraform"

  bigtable {
    project = "my_project"
    instance = "my_instance"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_user" "jane" {
  name       = "Jane"
  email      = "jane@example.com"
  given_name = "Jane"
  surname    = "Doe"
  password   = "hunter23"
  roles      = ["author"]
}

resource "anaml-operations_user" "john" {
  name       = "John"
  email      = "john@example.com"
  given_name = "John"
  surname    = "Doe"
  password   = "hunter23"
  roles      = ["super_user"]
}

resource "anaml-operations_access_token" "john_unpriviledged" {
  owner       = anaml-operations_user.john.id
  description = "Unprivileged Token"
  roles       = []
}

resource "anaml-operations_access_token" "john_runner" {
  owner       = anaml-operations_user.john.id
  description = "Runner Token"
  roles       = [ "run_featuregen", "run_monitoring", "run_caching" ]
}

resource "anaml-operations_caching" "caching" {
  name        = "household_caching"
  description = "Caching of tables for households"
  prefix_url  = "file:///tmp/anaml/caching"
  include {
    spec {
      table  = anaml_table.household.id
      entity = anaml_entity.household.id
    }
  }
  retainment = "PT48H"
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
}

resource "anaml-operations_caching" "caching_with_principal" {
  name        = "household_caching_with_principal"
  description = "Caching of tables for households"
  prefix_url  = "file:///tmp/anaml/caching"
  include {
    spec {
      table  = anaml_table.household.id
      entity = anaml_entity.household.id
    }
  }
  retainment = "PT48H"
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
  principal = anaml-operations_user.jane.id
}

resource "anaml-operations_caching" "caching_two" {
  name        = "household_caching_auto"
  description = "Caching of tables for households"
  prefix_url  = "file:///tmp/anaml/caching"
  auto {
    exclude {
      table  = anaml_table.household.id
      entity = anaml_entity.household.id
    }
  }
  retainment = "PT48H"
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
}

resource "anaml-operations_monitoring" "monitoring" {
  name        = "household_monitoring"
  description = "Monitoring of tables for households"
  enabled     = true
  include {
    tables = [
      anaml_table.household.id
    ]
  }
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
}

resource "anaml-operations_monitoring" "monitoring_with_principal" {
  name        = "household_monitoring_with_principal"
  description = "Monitoring of tables for households"
  enabled     = true
  include {
    tables = [
      anaml_table.household.id
    ]
  }
  cluster = data.anaml-operations_cluster.local.id
  daily_schedule {
    start_time_of_day = "00:00:00"
  }
  principal = anaml-operations_user.jane.id
}

resource "anaml-operations_user_group" "engineering" {
  name        = "Engineering"
  description = "A user group with engineering members."
  members {
    user_id = anaml-operations_user.jane.id
  }
  members {
    user_id = anaml-operations_user.john.id
  }
  roles = [
    "run_monitoring"
  ]
}


resource "anaml-operations_branch_protection" "official" {
  protection_pattern = "official"
  merge_approval_rules {
    restricted {
      num_required_approvals = 1
      approvers {
        user_group {
          id = anaml-operations_user_group.engineering.id
        }
      }
    }
  }
  merge_approval_rules {
    open {
      num_required_approvals = 2
    }
  }
  push_whitelist {
    user {
      id = anaml-operations_user.john.id
    }
  }
  apply_to_admins       = false
  allow_branch_deletion = false
}

resource "anaml-operations_webhook" "merge_hook" {
  name        = "merge_request_hook"
  description = "A hook for new merge requests"
  url         = "http://localhost:8095/hook"
  merge_requests {}
}


resource "anaml-operations_view_materialisation_job" "view_materialisation_batch_no_metadata" {
  name        = "view_materialisation_batch_no_metadata"
  description = "Materialise household normalised"
  cluster     = data.anaml-operations_cluster.local.id
  usagettl    = "PT48H"
  view {
    table                       = anaml_table.household_normalised.id
    destination {
        destination             = data.anaml-operations_destination.s3a.id
        folder {
          path = "household_normalised_view_results"
          partitioning_enabled = true
          save_mode = "overwrite"
        }
      }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}

resource "anaml-operations_view_materialisation_job" "view_materialisation_batch" {
  name        = "view_materialisation_batch"
  description = "Materialise household normalised"
  cluster     = data.anaml-operations_cluster.local.id
  usagettl    = "PT48H"
  include_metadata = true
  view {
    table                       = anaml_table.household_normalised.id
    destination {
        destination             = data.anaml-operations_destination.s3a.id
        folder {
          path = "household_normalised_view_results"
          partitioning_enabled = true
          save_mode = "overwrite"
        }
      }
  }
  daily_schedule {
    start_time_of_day = "00:00:00"
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}


resource "anaml-operations_view_materialisation_job" "view_materialisation_streaming" {
  name        = "view_materialisation_streaming"
  description = "Materialise household normalised"
  cluster     = data.anaml-operations_cluster.local.id
  usagettl    = "PT48H"

  view {
    table                       = anaml_table.household_normalised.id
    destination {
        destination             = data.anaml-operations_destination.s3a.id
        folder {
          path = "household_normalised_view_results"
          partitioning_enabled = true
          save_mode = "overwrite"
        }
      }
  }

  labels = [ anaml-operations_label_restriction.terraform.text ]
}