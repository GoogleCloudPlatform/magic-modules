resource "google_logging_metric" "logging_metric" {
  name = "my-(custom)/metric-${local.name_suffix}"
  filter = "resource.type=gae_app AND severity>=ERROR"
  metric_descriptor {
    metric_kind = "DELTA"
    value_type = "DISTRIBUTION"
    labels {
        key = "mass"
        value_type = "STRING"
        description = "amount of matter"
    }
  }
  value_extractor = "EXTRACT(jsonPayload.request)"
  label_extractors = { "mass": "EXTRACT(jsonPayload.request)" }
  bucket_options {
    linear_buckets {
      num_finite_buckets = 3
      width = 1
      offset = 1
    }
  }
}
