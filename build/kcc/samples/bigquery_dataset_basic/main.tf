resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "example_dataset-${local.name_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
  default_table_expiration_ms = 3600000

  labels = {
    env = "default"
  }

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    role   = "READER"
    domain = "example.com"
  }
}
