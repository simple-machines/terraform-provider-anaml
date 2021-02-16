provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "official"
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

resource "anaml_entity" "household" {
  name           = "household"
  description    = "A household level view"
  default_column = "household"
}

resource "anaml_table" "household" {
  name           = "household"
  description    = "A household level view"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.household.id) = "household_id"
    }
    timestamp_column = "timestamp"
  }
}

resource "anaml_table" "household_normalised" {
  name           = "household_normalised"
  description    = "A household level view"

  expression     = "SELECT * FROM household"
  sources        = [ anaml_table.household.id ]

  event {
    entities = {
      (anaml_entity.household.id) = "household"
    }
    timestamp_column = "timestamp"
  }
}

resource "anaml_feature_template" "household" {
  name           = "household_count"
  description    = "Count of household items"
  table          = anaml_table.household.id
  select         = "count"
  aggregations   = [ "sum" ]
  days           = [ 1, 2, 4 ]
}

data "anaml_feature" "household" {
  for_each       = toset(["1", "2", "4"])
  template       = anaml_feature_template.household.id
  days           = parseint(each.key, 10)
}

resource "anaml_feature_set" "household" {
  name           = "household"
  entity         = anaml_entity.household.id
  features       = [
      data.anaml_feature.household["1"].id
    , data.anaml_feature.household["2"].id
    , data.anaml_feature.household["4"].id
    ]
}

resource "anaml_feature_store" "household" {
  name           = "household"
  description    = "Daily view of households"
  feature_set    = anaml_feature_set.household.id
  mode           = "daily"
  enabled        = true
  cluster        = data.anaml_cluster.default_local.id
  destination {
    destination = data.anaml_destination.minio.id
    folder = "household"
  }
}
