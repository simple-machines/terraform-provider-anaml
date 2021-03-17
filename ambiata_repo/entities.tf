resource "anaml_entity" "town" {
  name           = "town"
  description    = "A town plan"
  default_column = "town_name"
}

resource "anaml_entity" "tower" {
  name           = "tower"
  description    = "A tower plan"
  default_column = "tower"
}

resource "anaml_entity" "customer" {
  name           = "customer"
  description    = "A customer identified in the system"
  default_column = "customer"
}

resource "anaml_entity" "phone_plan" {
  name           = "phone_plan"
  description    = "A phone plan"
  default_column = "plan"
}

resource "anaml_entity" "household" {
  name           = "household"
  description    = "A household level view"
  default_column = "household"
}

resource "anaml_entity" "store" {
  name           = "store"
  description    = "A supermarket near you"
  default_column = "store"
}

resource "anaml_entity_mapping" "household_to_customer" {
  from     = anaml_entity.household.id
  to       = anaml_entity.customer.id
  mapping  = anaml_feature.plan_customer.id
}

resource "anaml_entity_mapping" "plan_to_customer" {
  from     = anaml_entity.phone_plan.id
  to       = anaml_entity.customer.id
  mapping  = anaml_feature.plan_customer.id
}
