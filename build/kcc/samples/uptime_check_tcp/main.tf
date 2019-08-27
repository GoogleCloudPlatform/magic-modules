resource "google_monitoring_uptime_check_config" "tcp_group" {
  display_name = "tcp-uptime-check"
  timeout = "60s"

  tcp_check {
    port = 888
  }

  resource_group {
    resource_type = "INSTANCE"
    group_id = "${google_monitoring_group.check.name}"
  }
}


resource "google_monitoring_group" "check" {
  display_name = "uptime-check-group"
  filter = "resource.metadata.name=has_substring(\"foo\")"
}
