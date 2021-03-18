
resource "anaml_feature" "latest_bill_amount" {
  name           = "latest_bill_amount"
  description    = "Total amount charged in latest bill"
  table          = anaml_table.bills.id
  select         = "amount"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "bill_shock" {
  name           = "bill_shock"
  description    = "Total amount charged in latest bill"
  table          = anaml_table.bills.id
  select         = "amount"
  aggregation    = "percentagechange"
  rows           = 2
}

resource "anaml_feature" "number_of_bills_received" {
  name           = "number_of_bills_received"
  description    = "Total amount of bills received"
  table          = anaml_table.bills.id
  select         = "amount"
  aggregation    = "count"
  open           = true
}

resource "anaml_feature" "average_bill_amount" {
  name           = "average_bill_amount"
  description    = "Average amount charged across bills"
  table          = anaml_table.bills.id
  select         = "amount"
  aggregation    = "avg"
  open           = true
}
