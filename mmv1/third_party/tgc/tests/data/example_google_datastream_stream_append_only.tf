terraform {
  required_providers {
    google = {
      source = "hashicorp/google-beta"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_bigquery_dataset" "postgres" {
  dataset_id    = "stpostgres"
  friendly_name = "stpostgres"
  description   = "Database of postgres"
  location      = "us-central1"
}

resource "google_datastream_stream" "default" {
  display_name  = "postgres to bigQuery"
  location      = "us-central1"
  stream_id     = "postgres-bigquery"

   source_config {
    source_connection_profile = google_datastream_connection_profile.source_connection_profile.id
    mysql_source_config {}
  }

  destination_config {
    destination_connection_profile = google_datastream_connection_profile.destination_connection_profile2.id
    bigquery_destination_config {
      data_freshness = "900s"
      single_target_dataset {
        dataset_id = google_bigquery_dataset.postgres.id
      }
      append_only {}
    }
  }

  backfill_all {
  }

}

resource "google_datastream_connection_profile" "destination_connection_profile2" {
    display_name          = "Connection profile"
    location              = "us-central1"
    connection_profile_id = "dest-st-profile"
    bigquery_profile {}
}

resource "google_sql_database_instance" "instance" {
    name             = "instance-stream-name"
    database_version = "MYSQL_8_0"
    region           = "us-central1"
    settings {
        tier = "db-f1-micro"
        backup_configuration {
            enabled            = true
            binary_log_enabled = true
        }

        ip_configuration {
            // Datastream IPs will vary by region.
            authorized_networks {
                value = "34.71.242.81"
            }

            authorized_networks {
                value = "34.72.28.29"
            }

            authorized_networks {
                value = "34.67.6.157"
            }

            authorized_networks {
                value = "34.67.234.134"
            }

            authorized_networks {
                value = "34.72.239.218"
            }
        }
    }

    deletion_protection  = false
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name     = "my-user"
    instance = google_sql_database_instance.instance.name
    host     = "%"
    password = random_password.pwd.result
}

resource "google_datastream_connection_profile" "source_connection_profile" {
    display_name          = "Source connection profile"
    location              = "us-central1"
    connection_profile_id = "source-profile-st"

    mysql_profile {
        hostname = google_sql_database_instance.instance.public_ip_address
        username = google_sql_user.user.name
        password = google_sql_user.user.password
    }
}