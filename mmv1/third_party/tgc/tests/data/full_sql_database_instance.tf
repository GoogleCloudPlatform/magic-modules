/**
 * Copyright 2019 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

resource "google_sql_database_instance" "main" {
  name             = "main-instance"
  database_version = "POSTGRES_9_6"
  region           = "us-central1"

  depends_on = [
    google_service_networking_connection.private_vpc_connection
  ]

  master_instance_name = "test-master_instance_name"
  replica_configuration {
    ca_certificate          = "test-ca_certificate"
    client_certificate      = "test-client_certificate"
    client_key              = "test-client_key"
    connect_retry_interval  = 42
    dump_file_path          = "test-dump_file_path"
    failover_target         = true
    master_heartbeat_period = 42
    password                = "test-password"
    # TODO.
    # Got: An argument named "sslCipher" is not expected here. Did you mean "ssl_cipher"?
    ssl_cipher                = "test-sslCipher"
    username                  = "test-username"
    verify_server_certificate = true
  }
  settings {
    activation_policy           = "test-activation_policy"
    availability_type           = "REGIONAL"
    backup_configuration {
      binary_log_enabled = true
      enabled            = true
      start_time         = "42:42"
      location           = "us"
    }
    database_flags {
      name  = "test-name1"
      value = "test-value1"
    }
    database_flags {
      name  = "test-name2"
      value = "test-value2"
    }
    disk_autoresize = true
    disk_size       = 42
    disk_type       = "test-disk_type"
    ip_configuration {
      authorized_networks {
        name            = "test-authorized_networks-name1"
        value           = "test-authorized_networks-value1"
        expiration_time = "test-expiration_time"
      }
      authorized_networks {
        name            = "test-authorized_networks-name2"
        value           = "test-authorized_networks-value2"
        expiration_time = "test-expiration_time"
      }
      ipv4_enabled    = true
      private_network = google_compute_network.private_network.self_link
      ssl_mode        = "TRUSTED_CLIENT_CERTIFICATE_REQUIRED"
    }
    location_preference {
      follow_gae_application = "test-follow_gae_application"
      zone                   = "us-central1-a"
      secondary_zone         = "us-central1-b"
    }
    maintenance_window {
      day          = 42
      hour         = 42
      update_track = "test-update_track"
    }
    pricing_plan     = "test-pricing_plan"
    tier             = "db-f1-micro"
    user_labels = {
      user_labels_foo = "user_labels_bar"
    }
  }
}


resource "google_compute_network" "private_network" {
  name = "private-network"
}

resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.self_link
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.private_network.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}
