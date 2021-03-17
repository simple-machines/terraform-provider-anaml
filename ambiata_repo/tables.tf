
resource "anaml_table" "towns" {
  name           = "towns"
  description    = "Town information"
  source {
    source = data.anaml_destination.minio.id
    folder = "towns"
  }

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
  source {
    source = data.anaml_destination.minio.id
    folder = "towers"
  }
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
  source {
    source = data.anaml_destination.minio.id
    folder = "plan"
  }
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
  source {
    source = data.anaml_destination.minio.id
    folder = "customers"
  }
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
  source {
    source = data.anaml_destination.minio.id
    folder = "data-usage"
  }
  event {
    entities = {
      (anaml_entity.phone_plan.id) = "plan"
    }
    timestamp_column = "usage_time"
  }
}
