resource "google_cloud_scheduler_job" "job" {
  name     = "test-job"
  description = "test http job"
  schedule = "*/8 * * * *"
  time_zone = "America/New_York"

  http_target {
    http_method = "POST"
    uri = "https://example.com/ping"
  }
}
