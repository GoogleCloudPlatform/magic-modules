resource "google_monitoring_notification_channel" "basic" {
  display_name = "Test Notification Channel-${local.name_suffix}"
  type = "email"
  labels = {
    email_address = "fake_email@blahblah.com"
  }
}
