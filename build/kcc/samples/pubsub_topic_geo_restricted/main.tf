resource "google_pubsub_topic" "example" {
  name = "example-topic-${local.name_suffix}"

  message_storage_policy {
    allowed_persistence_regions = [
      "europe-west3",
    ]
  }

}
