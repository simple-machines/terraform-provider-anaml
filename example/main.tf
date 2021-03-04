provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "official"
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

  source {
    source = data.anaml_destination.minio.id
    folder = "household"
  }

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

resource "anaml_feature_template" "household_count" {
  name           = "household_count"
  description    = "Count of household items"
  table          = anaml_table.household.id
  select         = "count"
  aggregation    = "sum"
}

resource "anaml_feature" "household_count" {
  for_each       = toset(["1", "2", "4"])
  days           = parseint(each.key, 10)

  name           = "household_count_${each.key}_days"
  description    = "Count of household items"
  table          = anaml_table.household.id
  select         = "count"
  aggregation    = "sum"
  template       = anaml_feature_template.household_count.id
}

resource "anaml_feature_set" "household" {
  name           = "household"
  entity         = anaml_entity.household.id
  features       = [
      anaml_feature.household_count["1"].id
    , anaml_feature.household_count["2"].id
    , anaml_feature.household_count["4"].id
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
    folder = "household_results"
  }
}
