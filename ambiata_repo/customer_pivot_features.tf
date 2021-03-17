
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
