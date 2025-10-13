resource "google_logging_project_sink" "my-sink" {
  name   = "my-sink"
  project = "my-project"
  destination = "bigquery.googleapis.com/projects/my-project/datasets/my_dataset"
  filter = "severity >= ERROR"
  description = "A sink for errors"
  disabled = false

  exclusions {
    name        = "exclude-debug"
    description = "Exclude debug logs"
    filter      = "severity < INFO"
    disabled    = false
  }

  bigquery_options {
    use_partitioned_tables = true
  }
}
