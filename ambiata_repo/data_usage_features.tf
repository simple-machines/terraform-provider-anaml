
# Domains
# = Other
# | Social
# | Music
# | Sport
# | News
# | Video

resource "anaml_feature" "total_data_usage" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "total_plan_data_usage_${each.key}_days"
  description    = "Total data usage over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  aggregation    = "sum"
  days           = parseint(each.key, 10)
}

resource "anaml_feature" "total_data_usage_video" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "total_plan_video_data_usage_${each.key}_days"
  description    = "Total data usage for video over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  filter         = "domain = 'video'"
  aggregation    = "sum"
  days           = parseint(each.key, 10)
}

resource "anaml_feature" "total_data_usage_music" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "total_plan_music_data_usage_${each.key}_days"
  description    = "Total data usage for music over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  filter         = "domain = 'music'"
  aggregation    = "sum"
  days           = parseint(each.key, 10)
}

resource "anaml_feature" "total_data_usage_sport" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "total_plan_sport_data_usage_${each.key}_days"
  description    = "Total data usage for sport over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "megabytes"
  filter         = "domain = 'sport'"
  aggregation    = "sum"
  days           = parseint(each.key, 10)
}


resource "anaml_feature" "count_slow_data_usage" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "count_slow_data_usage_${each.key}_days"
  description    = "Count of data usage issues over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "1"
  filter         = "issue rlike 'slow'"
  aggregation    = "count"
  days           = parseint(each.key, 10)
}

resource "anaml_feature" "count_data_loss" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "count_data_loss_${each.key}_days"
  description    = "Count of data lost packets past ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "1"
  filter         = "issue rlike 'data_loss'"
  aggregation    = "count"
  days           = parseint(each.key, 10)
}



resource "anaml_feature" "count_data_packets" {
  for_each       = toset( ["7", "14", "28", "56"] )
  name           = "count_data_packets_${each.key}_days"
  description    = "Count of data usage issues over the last ${each.key} days"
  table          = anaml_table.data_usage.id
  select         = "1"
  aggregation    = "count"
  days           = parseint(each.key, 10)
}

