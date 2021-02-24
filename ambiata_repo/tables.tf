
resource "anaml_table" "towns" {
  name           = "towns"
  description    = "Town information"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.town.id) = "town_name"
    }
    timestamp_column = "timestamp"
  }
}

resource "anaml_table" "towers" {
  name           = "towers"
  description    = "Towers we provide for customers"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.tower.id) = "tower_id"
    }
    timestamp_column = "timestamp"
  }
}

resource "anaml_table" "plans" {
  name           = "plans"
  description    = "Information about phone plans"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.phone_plan.id) = "plan_id"
    }
    timestamp_column = "start_date"
  }
}

resource "anaml_table" "customer" {
  name           = "customers"
  description    = "Customer demographic information"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.customer.id) = "customer_id"
    }
    timestamp_column = "join_date"
  }
}

resource "anaml_table" "data_usage" {
  name           = "data_usage"
  description    = "Plan data usage"
  source         = data.anaml_source.minio.id

  event {
    entities = {
      (anaml_entity.phone_plan.id) = "plan_id"
    }
    timestamp_column = "usage_timestamp"
  }
}
