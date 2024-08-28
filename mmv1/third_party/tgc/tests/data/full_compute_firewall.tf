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

resource "google_compute_firewall" "full_list_default_1" {
  name    = "test-firewall1"
  network = google_compute_network.default.name

  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = [80, 8080, "1000-2000"]
  }
  description        = "test-description"
  destination_ranges = ["192.168.0.42/32", "192.168.0.43/32"]
  direction          = "INGRESS"
  disabled           = true
  # TODO: beta feature
  # Got: An argument named "enable_logging" is not expected here.
  # enable_logging = true
  priority = 42
  source_service_accounts = ["test-source_service_account1", "test-source_service_account2"]
}

resource "google_compute_firewall" "full_list_default_2" {
  name    = "test-firewall2"
  network = google_compute_network.default.name

  deny {
    protocol = "icmp"
  }
  deny {
    protocol = "tcp"
    ports    = [80, 8080, "1000-2000"]
  }
  source_ranges           = ["192.168.0.44/32", "192.168.0.45/32"]
  source_service_accounts = ["test-source_service_account1", "test-source_service_account2"]
  target_service_accounts = ["test-target_service_account1", "test-target_service_account2"]
}

resource "google_compute_firewall" "full_list_default_3" {
  name    = "test-firewall3"
  network = google_compute_network.default.name

  deny {
    protocol = "icmp"
  }
  source_tags = ["web"]
  target_tags = ["test-target_tag1", "test-target_tag2"]
}

resource "google_compute_network" "default" {
  name = "test-network"
}
