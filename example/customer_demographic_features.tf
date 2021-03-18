
resource "anaml_feature" "customer_age" {
  name           = "age"
  description    = "Age of a customer"
  table          = anaml_table.customer.id
  select         = "age"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "customer_town" {
  name           = "town"
  description    = "Most recent town information for a customer"
  table          = anaml_table.customer.id
  select         = "town"
  aggregation    = "last"
  open           = true
}
