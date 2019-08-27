resource "google_logging_metric" "logging_metric" {
  name = "my-(custom)/metric-${local.name_suffix}"
  filter = "resource.type=gae_app AND severity>=ERROR"
  metric_descriptor {
    metric_kind = "DELTA"
    value_type = "INT64"
  }
}
