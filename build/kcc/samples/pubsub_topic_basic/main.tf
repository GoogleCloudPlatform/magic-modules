resource "google_pubsub_topic" "example" {
  name = "example-topic"

  labels = {
    foo = "bar"
  }
}
