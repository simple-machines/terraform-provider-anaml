
provider "anaml" {
  host     = "http://127.0.0.1:8080/api"
  username = "admin"
  password = "test password"
  branch   = "development"
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
  name           = "household_normalised"
  description    = "A household level view"

  select = "these"
  window {
    rows = 3
  }
  table = anaml_table.household.id

  aggregation = "sum"

}
