resource "google_pubsub_topic" "example" {
  name = "example-topic-${local.name_suffix}"
}

resource "google_pubsub_subscription" "example" {
  name  = "example-subscription-${local.name_suffix}"
  topic = "${google_pubsub_topic.example.name}"

  labels = {
    foo = "bar"
  }

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages = true

  ack_deadline_seconds = 20

  expiration_policy {
    ttl = "300000.5s"
  }
}
