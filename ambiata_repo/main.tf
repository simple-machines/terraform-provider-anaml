provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "development"
  version  = "0.3.4"
}

data "anaml_source" "minio" {
  name = "Minio S3 Source"
}

data "anaml_cluster" "default_local" {
  name = "Default Local"
}

data "anaml_destination" "minio" {
  name = "Minio S3 Destination"
}
