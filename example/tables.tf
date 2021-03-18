
resource "anaml_table" "towns" {
  name           = "towns"
  description    = "Town information"
  source {
    source = data.anaml_source.s3a.id
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
    source = data.anaml_source.s3a.id
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
    source = data.anaml_source.s3a.id
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
    source = data.anaml_source.s3a.id
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
    source = data.anaml_source.s3a.id
    folder = "data-usage"
  }
  event {
    entities = {
      (anaml_entity.phone_plan.id) = "plan"
    }
    timestamp_column = "usage_time"
  }
}

resource "anaml_table" "bills" {
  name           = "bills"
  description    = "Plan billing usage"
  source {
    source = data.anaml_source.s3a.id
    folder = "bills"
  }
  event {
    entities = {
      (anaml_entity.customer.id) = "customer"
      (anaml_entity.phone_plan.id) = "plan"
    }
    timestamp_column = "end_billing_period"
  }
}

resource "anaml_table" "transactions" {
  name           = "transactions"
  description    = "Supermarket transactions"
  source {
    source = data.anaml_source.s3a.id
    folder = "transactions"
  }
  event {
    entities = {
      (anaml_entity.customer.id) = "customer"
      (anaml_entity.store.id) = "store"
    }
    timestamp_column = "time"
  }
}
