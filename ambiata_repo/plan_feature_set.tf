

resource "anaml_feature_set" "plan_information" {
  name           = "plan_information"
  entity         = anaml_entity.phone_plan.id
  features       = [
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


   , anaml_feature.count_data_packets["7"].id
   , anaml_feature.count_data_packets["14"].id
   , anaml_feature.count_data_packets["28"].id
   , anaml_feature.count_data_packets["56"].id

   , anaml_feature.count_slow_data_usage["7"].id
   , anaml_feature.count_slow_data_usage["14"].id
   , anaml_feature.count_slow_data_usage["28"].id
   , anaml_feature.count_slow_data_usage["56"].id

   , anaml_feature.count_data_loss["7"].id
   , anaml_feature.count_data_loss["14"].id
   , anaml_feature.count_data_loss["28"].id
   , anaml_feature.count_data_loss["56"].id
   ]
}
