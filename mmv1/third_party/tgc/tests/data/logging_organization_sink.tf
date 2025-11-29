resource "google_bigquery_dataset" "my_dataset" {
  project    = "{{.Provider.project}}"
  dataset_id = "my_dataset"
  location   = "US"
}

resource "google_logging_organization_sink" "my_sink" {
  name             = "gg-asset-88093-71a3-sink"
  org_id           = "{{.OrgID}}"
  destination      = "bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/${google_bigquery_dataset.my_dataset.dataset_id}"
  include_children = false
  exclusions {
    name        = "gg-asset-88093-71a3-exclusion"
    description = "Exclude all GCE instance logs"
    filter      = "resource.type = gce_instance"
    disabled    = true
  }
}

resource "google_logging_organization_sink" "my_sink_with_intercept" {
  name             = "gg-asset-88093-71a3-sink-with-intercept"
  org_id           = "{{.OrgID}}"
  destination      = "bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/${google_bigquery_dataset.my_dataset.dataset_id}"
  intercept_children = true
  exclusions {
    name        = "gg-asset-88093-71a3-exclusion"
    description = "Exclude all GCE instance logs"
    filter      = "resource.type = gce_instance"
    disabled    = true
  }
}

resource "google_bigquery_dataset_iam_member" "my_iam" {
  project    = "{{.Provider.project}}"
  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  role       = "roles/bigquery.dataEditor"
  member     = google_logging_organization_sink.my_sink.writer_identity
}

resource "google_logging_organization_sink" "my_sink_with_children" {
  name             = "gg-asset-88093-71a3-sink-with-children"
  org_id           = "{{.OrgID}}"
  destination      = "bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/${google_bigquery_dataset.my_dataset.dataset_id}"
  include_children = true
  exclusions {
    name        = "gg-asset-88093-71a3-exclusion"
    description = "Exclude all GCE instance logs"
    filter      = "resource.type = gce_instance"
    disabled    = true
  }
}
