resource "anaml_feature" "plan_customer" {
  name           = "plan_customer"
  description    = "Customer for a plan"
  table          = anaml_table.plans.id
  select         = "customer"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_size" {
  name           = "plan_size"
  description    = "Plan size"
  table          = anaml_table.plans.id
  select         = "size"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_roaming_enabled" {
  name           = "roaming_enabled"
  description    = "Whether roaming is enabled"
  table          = anaml_table.plans.id
  select         = "roaming"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_free_sport_enabled" {
  name           = "free_sport_enabled"
  description    = "Whether free sport is enabled"
  table          = anaml_table.plans.id
  select         = "free_sport"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_free_music_enabled" {
  name           = "free_music_enabled"
  description    = "Whether free music is enabled"
  table          = anaml_table.plans.id
  select         = "free_music"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_is_business" {
  name           = "is_business_plan"
  description    = "Whether the plan is a business plan"
  table          = anaml_table.plans.id
  select         = "is_business"
  aggregation    = "last"
  open           = true
}

resource "anaml_feature" "plan_age" {
  name           = "plan_age"
  description    = "Plan age"
  table          = anaml_table.plans.id
  select         = "datediff(feature_date(), start_date)"
  aggregation    = "last"
  open           = true
}
