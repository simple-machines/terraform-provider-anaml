data "anaml_source" "s3a" {
  name = anaml-operations_source.s3a.name
}

data "anaml_cluster" "local" {
  name = anaml-operations_cluster.local.name
}

data "anaml_destination" "s3a" {
  name = anaml-operations_destination.s3a.name
}


resource "anaml-operations_cluster" "local" {
  name               = "Terraform Local Cluster"
  description        = "A local cluster created by Terraform"
  is_preview_cluster = true

  local {
    anaml_server_url = "http://localhost:8080"
    jwt_token_provider {
      login_server_url = "http://localhost:8080"
      basic {
        username = "admin"
        password = "test password"
      }
    }
  }

  spark_config {
    enable_hive_support = true
  }
}

resource "anaml-operations_cluster" "spark_server" {
  name               = "Terraform Spark Server Cluster"
  description        = "A Spark server cluster created by Terraform"
  is_preview_cluster = false

  spark_server {
    spark_server_url = "http://localhost:8080"
  }

  spark_config {
    enable_hive_support = true
  }
}

resource "anaml-operations_source" "s3" {
  name        = "Terraform S3 Source"
  description = "An S3 source created by Terraform"

  s3 {
    bucket         = "my-bucket"
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = true
  }
}

resource "anaml-operations_source" "s3a" {
  name        = "Terraform S3A Source"
  description = "An S3A source created by Terraform"

  s3a {
    bucket      = "spark-warehouse"
    path        = "/"
    endpoint    = "http://localhost:9000"
    file_format = "orc"
    access_key  = "accesskey"
    secret_key  = "secretkey"
  }
}

resource "anaml-operations_source" "hive" {
  name        = "Terraform Hive Source"
  description = "An Hive source created by Terraform"

  hive {
    database = "my_database"
  }
}

resource "anaml-operations_source" "jdbc" {
  name        = "Terraform JDBC Source"
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
}

resource "anaml-operations_source" "big_query" {
  name        = "Terraform BigQuery Source"
  description = "An BigQuery source created by Terraform"

  big_query {
    path = "/path/to/file"
  }
}

resource "anaml-operations_source" "gcs" {
  name        = "Terraform GCS Source"
  description = "An GCS source created by Terraform"

  gcs {
    bucket         = "my-bucket"
    path           = "/path/to/file"
    file_format    = "parquet"
  }
}

resource "anaml-operations_source" "local" {
  name        = "Terraform Local Source"
  description = "An Local source created by Terraform"

  local {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }
}

resource "anaml-operations_source" "hdfs" {
  name        = "Terraform HDFS Source"
  description = "An HDFS source created by Terraform"

  hdfs {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }
}

resource "anaml-operations_source" "kafka" {
  name        = "Terraform Kafka Source"
  description = "An Kafka source created by Terraform"

  kafka {
    bootstrap_servers = "http://bootstrap"
    schema_registry_url = "http://schema-registry"
  }
}

resource "anaml-operations_destination" "s3" {
  name        = "Terraform S3 Destination"
  description = "An S3 destination created by Terraform"

  s3 {
    bucket         = "my-bucket"
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = true
  }
}

resource "anaml-operations_destination" "s3a" {
  name        = "Terraform S3A Destination"
  description = "An S3A destination created by Terraform"

  s3a {
    bucket      = "my-bucket"
    path        = "/path/to/file"
    endpoint    = "http://example.com"
    file_format = "orc"
    access_key  = "access"
    secret_key  = "secret"
  }
}

resource "anaml-operations_destination" "hive" {
  name        = "Terraform Hive Destination"
  description = "An Hive destination created by Terraform"

  hive {
    database = "my_database"
  }
}

resource "anaml-operations_destination" "jdbc" {
  name        = "Terraform JDBC Destination"
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
}

resource "anaml-operations_destination" "big_query_temporary" {
  name        = "Terraform BigQuery Destination with Temporary Staging Area"
  description = "An BigQuery destination created by Terraform"

  big_query {
    path = "/path/to/file"
    temporary_staging_area {
      bucket = "my-bucket"
    }
  }
}

resource "anaml-operations_destination" "big_query_persistent" {
  name        = "Terraform BigQuery Destination with Persistent Staging Area"
  description = "An BigQuery destination created by Terraform"

  big_query {
    path = "/path/to/file"
    persistent_staging_area {
      bucket = "my-bucket"
      path = "/path/to/file"
    }
  }
}

resource "anaml-operations_destination" "gcs" {
  name        = "Terraform GCS Destination"
  description = "An GCS destination created by Terraform"

  gcs {
    bucket         = "my-bucket"
    path           = "/path/to/file"
    file_format    = "parquet"
  }
}

resource "anaml-operations_destination" "local" {
  name        = "Terraform Local Destination"
  description = "An Local destination created by Terraform"

  local {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }
}

resource "anaml-operations_destination" "hdfs" {
  name        = "Terraform HDFS Destination"
  description = "An HDFS destination created by Terraform"

  hdfs {
    path           = "/path/to/file"
    file_format    = "csv"
    include_header = false
  }
}

resource "anaml-operations_destination" "kafka" {
  name        = "Terraform Kafka Destination"
  description = "An Kafka destination created by Terraform"

  kafka {
    bootstrap_servers = "http://bootstrap"
    schema_registry_url = "http://schema-registry"
  }
}

resource "anaml-operations_user" "jane" {
  name       = "Jane"
  email      = "jane@example.com"
  given_name = "Jane"
  surname    = "Doe"
  password   = "hunter2"
  roles      = ["viewer", "operator", "author"]
}

resource "anaml-operations_user" "john" {
  name       = "John"
  email      = "john@example.com"
  given_name = "John"
  surname    = "Doe"
  password   = "hunter2"
  roles      = ["super_user"]
}
