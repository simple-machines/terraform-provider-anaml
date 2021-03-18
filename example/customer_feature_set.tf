
resource "anaml_feature_set" "customer_information" {
  name           = "customer_information"
  entity         = anaml_entity.customer.id
  features       = [
     anaml_feature.customer_age.id
   , anaml_feature.customer_town.id

   , anaml_feature.num_plans.id
   , anaml_feature.count_roaming_enabled.id

   , anaml_feature.customer_total_data_usage["7"].id
   , anaml_feature.customer_total_data_usage["14"].id
   , anaml_feature.customer_total_data_usage["28"].id
   , anaml_feature.customer_total_data_usage["56"].id

   , anaml_feature.customer_count_slow_data_usage["7"].id
   , anaml_feature.customer_count_slow_data_usage["14"].id
   , anaml_feature.customer_count_slow_data_usage["28"].id
   , anaml_feature.customer_count_slow_data_usage["56"].id

   , anaml_feature.customer_count_data_loss["7"].id
   , anaml_feature.customer_count_data_loss["14"].id
   , anaml_feature.customer_count_data_loss["28"].id
   , anaml_feature.customer_count_data_loss["56"].id

   , anaml_feature.count_line_items.id
   , anaml_feature.count_baskets.id
   , anaml_feature.count_stores.id

   , anaml_feature.count_line_items_last_n_days["7"].id
   , anaml_feature.count_line_items_last_n_days["14"].id
   , anaml_feature.count_line_items_last_n_days["28"].id
   , anaml_feature.count_line_items_last_n_days["56"].id
   , anaml_feature.count_line_items_last_n_days["84"].id

   , anaml_feature.count_baskets_last_n_days["7"].id
   , anaml_feature.count_baskets_last_n_days["14"].id
   , anaml_feature.count_baskets_last_n_days["28"].id
   , anaml_feature.count_baskets_last_n_days["56"].id
   , anaml_feature.count_baskets_last_n_days["84"].id

   , anaml_feature.transaction_spend_last_n_days["7"].id
   , anaml_feature.transaction_spend_last_n_days["14"].id
   , anaml_feature.transaction_spend_last_n_days["28"].id
   , anaml_feature.transaction_spend_last_n_days["56"].id
   , anaml_feature.transaction_spend_last_n_days["84"].id
  ]
}


resource "anaml_feature_store" "customer_information" {
  name           = "customer_information"
  description    = "Daily customer information runs"
  enabled        = true
  mode           = "daily"
  feature_set    = anaml_feature_set.customer_information.id
  cluster        = data.anaml_cluster.local.id

  destination {
    destination = data.anaml_destination.s3a.id
    folder = "somewhere"
  }
}
