provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "development"
  version  = "0.3.4"
}

resource "anaml_entity" "household" {
  name           = "household"
  description    = "A household level view"
  default_column = "household"
}

resource "anaml_table" "household" {
  name           = "household"
  description    = "A household level view"

  event {
    entity = anaml_entity.household.id
    key_column = "household_id"
    timestamp_column = "timestamp"
  }
}

resource "anaml_table" "household_normalised" {
  name           = "household_normalised"
  description    = "A household level view"

  expression     = "SELECT * FROM household"
  sources        = [ anaml_table.household.id ]

  event {
    entity = anaml_entity.household.id
    key_column = "household"
    timestamp_column = "timestamp"
  }
}

resource "anaml_feature" "household" {
  name           = "household_count"
  description    = "Count of household items"
  table          = anaml_table.household.id
  select         = "count"
  aggregation    = "sum"
  days           = 4
}


resource "anaml_feature_set" "household" {
  name           = "household"
  entity         = anaml_entity.household.id
  features       = [ anaml_feature.household.id ]
}

resource "anaml_feature_store" "household" {
  name           = "household"
  description    = "Daily view of households"
  feature_set    = anaml_feature_set.household.id
  mode           = "daily"
  namespace      = "household"
}
