
resource "anaml_feature" "count_line_items" {
  name           = "count_line_items_bought"
  description    = "The number of line items a customer has bought"
  table          = anaml_table.transactions.id
  select         = "1"
  aggregation    = "count"
  open           = true
}

resource "anaml_feature" "count_baskets" {
  name           = "count_baskets"
  description    = "The number of baskets a customer used"
  table          = anaml_table.transactions.id
  select         = "basket"
  aggregation    = "countdistinct"
  open           = true
}

resource "anaml_feature" "count_stores" {
  name           = "count_stores"
  description    = "The number of stores visited"
  table          = anaml_table.transactions.id
  select         = "store"
  aggregation    = "countdistinct"
  open           = true
}
