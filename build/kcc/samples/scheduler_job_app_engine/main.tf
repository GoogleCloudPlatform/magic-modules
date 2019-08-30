resource "google_cloud_scheduler_job" "job" {
  name     = "test-job-${local.name_suffix}"
  schedule = "*/4 * * * *"
  description = "test app engine job"
  time_zone = "Europe/London"

  app_engine_http_target {
    http_method = "POST"

    app_engine_routing {
      service = "web"
      version = "prod"
      instance = "my-instance-001"
    }

    relative_uri = "/ping"
  }
}
