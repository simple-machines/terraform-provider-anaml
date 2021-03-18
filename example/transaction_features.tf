
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
  description    = "The number of supermarket trips a customer has had"
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


resource "anaml_feature_template" "count_line_items_last_n_days" {
  name           = "count_line_items_bought"
  description    = "The number of line items a customer has bought"
  table          = anaml_table.transactions.id
  select         = "1"
  aggregation    = "count"
}

resource "anaml_feature" "count_line_items_last_n_days" {
  for_each       = toset( ["7", "14", "28", "56", "84"] )
  name           = "count_line_items_bought_last_${each.key}_days"
  description    = "The number of line items a customer has bought"
  table          = anaml_table.transactions.id
  select         = "1"
  aggregation    = "count"
  days           = parseint(each.key, 10)
  template       = anaml_feature_template.count_line_items_last_n_days.id
}

resource "anaml_feature_template" "count_baskets_last_n_days" {
  name           = "count_baskets_last_n_days"
  description    = "The number of supermarket trips a customer has had in the past n days"
  table          = anaml_table.transactions.id
  select         = "basket"
  aggregation    = "countdistinct"
}

resource "anaml_feature" "count_baskets_last_n_days" {
  for_each       = toset( ["7", "14", "28", "56", "84"] )
  name           = "count_baskets_last_${each.key}_days"
  description    = "The number of supermarket trips a customer has had in the past ${each.key} days"
  table          = anaml_table.transactions.id
  select         = "basket"
  aggregation    = "countdistinct"
  days           = parseint(each.key, 10)
  template       = anaml_feature_template.count_baskets_last_n_days.id
}

resource "anaml_feature_template" "transaction_spend_last_n_days" {
  name           = "transaction_spend_last_n_days"
  description    = "The amount spent in supermarket trips in the past n days"
  table          = anaml_table.transactions.id
  select         = "cost"
  aggregation    = "sum"
}

resource "anaml_feature" "transaction_spend_last_n_days" {
  for_each       = toset( ["7", "14", "28", "56", "84"] )
  name           = "transaction_spend_last_${each.key}_days"
  description    = "The amount spent in supermarket trips in the past ${each.key} days"
  table          = anaml_table.transactions.id
  select         = "cost"
  aggregation    = "sum"
  days           = parseint(each.key, 10)
  template       = anaml_feature_template.transaction_spend_last_n_days.id
}

