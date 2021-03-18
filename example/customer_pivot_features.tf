
resource "anaml_table" "customer_plan_information" {
  name           = "customer_plan_information"
  description    = "Customer plans with data usage information."
  entity_mapping = anaml_entity_mapping.plan_to_customer.id

  extra_features = [
     anaml_feature.plan_size.id
   , anaml_feature.plan_roaming_enabled.id
   , anaml_feature.plan_free_sport_enabled.id
   , anaml_feature.plan_free_music_enabled.id
   , anaml_feature.plan_is_business.id
   , anaml_feature.plan_age.id

   , anaml_feature.total_data_usage["7"].id
   , anaml_feature.total_data_usage["14"].id
   , anaml_feature.total_data_usage["28"].id
   , anaml_feature.total_data_usage["56"].id

   , anaml_feature.count_slow_data_usage["7"].id
   , anaml_feature.count_slow_data_usage["14"].id
   , anaml_feature.count_slow_data_usage["28"].id
   , anaml_feature.count_slow_data_usage["56"].id

   , anaml_feature.count_data_loss["7"].id
   , anaml_feature.count_data_loss["14"].id
   , anaml_feature.count_data_loss["28"].id
   , anaml_feature.count_data_loss["56"].id


   , anaml_feature.total_data_usage_video["7"].id
   , anaml_feature.total_data_usage_video["14"].id
   , anaml_feature.total_data_usage_video["28"].id
   , anaml_feature.total_data_usage_video["56"].id

   , anaml_feature.total_data_usage_music["7"].id
   , anaml_feature.total_data_usage_music["14"].id
   , anaml_feature.total_data_usage_music["28"].id
   , anaml_feature.total_data_usage_music["56"].id

   , anaml_feature.total_data_usage_sport["7"].id
   , anaml_feature.total_data_usage_sport["14"].id
   , anaml_feature.total_data_usage_sport["28"].id
   , anaml_feature.total_data_usage_sport["56"].id
   ]
}

resource "anaml_feature" "num_plans" {
  name           = "num_plans"
  description    = "How many plans a customer has"
  table          = anaml_table.customer_plan_information.id
  select         = "1"
  aggregation    = "count"
  open           = true
}

resource "anaml_feature" "count_roaming_enabled" {
  name           = "count_roaming_enabled"
  description    = "Whether roaming is enabled"
  table          = anaml_table.customer_plan_information.id
  select         = "roaming_enabled"
  filter         = "roaming_enabled"
  aggregation    = "count"
  open           = true
}

resource "anaml_feature" "count_free_sport_enabled" {
  name           = "count_free_sport_enabled"
  description    = "Whether free sport is enabled"
  table          = anaml_table.customer_plan_information.id
  select         = "free_sport_enabled"
  filter         = "free_sport_enabled"
  aggregation    = "count"
  open           = true
}

resource "anaml_feature_template" "customer_total_data_usage" {
  name           = "customer_total_plan_data_usage_n_days"
  description    = "Total data usage for all of a customer's plans over the last n days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  aggregation    = "sum"
}

resource "anaml_feature" "customer_total_data_usage" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "customer_total_plan_data_usage_${each.key}_days"
  description    = "Total data usage for all of a customer's plans over the last ${each.key} days"
  table          = anaml_table.customer_plan_information.id
  select         = "total_plan_data_usage_${each.key}_days"
  aggregation    = "sum"
  open           = true
  template       = anaml_feature_template.customer_total_data_usage.id
}

resource "anaml_feature_template" "customer_count_slow_data_usage" {
  name           = "customer_count_slow_data_usage"
  description    = "Total data usage for all of a customer's plans over the last n days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  aggregation    = "sum"
}

resource "anaml_feature" "customer_count_slow_data_usage" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "customer_count_slow_data_usage_${each.key}_days"
  description    = "Total number of slow data usage events for all of a customer's plans over the last ${each.key} days"
  table          = anaml_table.customer_plan_information.id
  select         = "count_slow_data_usage_${each.key}_days"
  aggregation    = "sum"
  open           = true
  template       = anaml_feature_template.customer_count_slow_data_usage.id
}

resource "anaml_feature_template" "customer_count_data_loss" {
  name           = "customer_count_data_loss"
  description    = "Total data usage for all of a customer's plans over the last n days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  aggregation    = "sum"
}

resource "anaml_feature" "customer_count_data_loss" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "customer_count_data_loss_${each.key}_days"
  description    = "Total number of slow data usage events for all of a customer's plans over the last ${each.key} days"
  table          = anaml_table.customer_plan_information.id
  select         = "count_data_loss_${each.key}_days"
  aggregation    = "sum"
  open           = true
  template       = anaml_feature_template.customer_count_data_loss.id
}
