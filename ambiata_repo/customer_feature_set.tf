
resource "anaml_feature_set" "customer_information" {
  name           = "customer_information"
  entity         = anaml_entity.customer.id
  features       = [
     anaml_feature.customer_age.id
   , anaml_feature.customer_town.id
  ]
}

resource "anaml_feature_store" "customer_information" {
  name           = "customer_information"
  description    = "Daily customer information runs"
  enabled        = true
  mode           = "daily"
  feature_set    = anaml_feature_set.customer_information.id
  cluster        = data.anaml_cluster.default_local.id

  destination {
    destination = data.anaml_destination.minio.id
    folder = "somewhere"
  }
}
