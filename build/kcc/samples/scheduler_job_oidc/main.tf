data "google_compute_default_service_account" "default" { }

resource "google_cloud_scheduler_job" "job" {
  name     = "test-job-${local.name_suffix}"
  description = "test http job"
  schedule = "*/8 * * * *"
  time_zone = "America/New_York"

  http_target {
    http_method = "GET"
    uri = "https://example.com/ping"

    oidc_token {
      service_account_email = "${data.google_compute_default_service_account.default.email}"
    }
  }
}
