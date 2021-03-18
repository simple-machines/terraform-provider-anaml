terraform {
  required_providers {
    anaml = {
      source  = "registry.anaml.io/anaml/anaml"
      version = "1.0.0"
    }
    anaml-operations = {
      source  = "registry.anaml.io/anaml/anaml-operations"
      version = "1.0.0"
    }
  }
}

provider "anaml-operations" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
}

provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "official"
}
