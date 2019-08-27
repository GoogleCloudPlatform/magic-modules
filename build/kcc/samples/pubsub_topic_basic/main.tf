resource "google_pubsub_topic" "example" {
  name = "example-topic-${local.name_suffix}"

  labels = {
    foo = "bar"
  }
}
