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

resource "google_container_node_pool" "full_list_default_1" {
  name = "test-node-pool"

  # TODO: add beta fields.
  # node_config.sandbox_config
  # node_config.workload_metadata_config
  # node_locations
  # taint

  autoscaling {
    min_node_count = 42
    max_node_count = 1337
  }
  cluster = "test-cluster"
  initial_node_count = 42
  location = "us-central1"
  management {
    auto_repair = true
    auto_upgrade = true
  }
  max_pods_per_node = 42
  node_config {
    disk_size_gb = 42
    disk_type = "pd-standard"
    guest_accelerator {
      type  = "test-type1"
      count = 1
    }
    guest_accelerator {
      type  = "test-type2"
      count = 1
    }
    image_type = "test-image_type"
    labels = {
      test-label-key = "test-label-value"
    }
    local_ssd_count = 42
    machine_type = "test-machine_type"
    metadata = {
      test-metadata-key = "test-metadata-value"
    }
    min_cpu_platform = "test-min_cpu_platform"
    oauth_scopes = ["test-oauth_scopes", "storage-ro"]
    preemptible = true
    service_account = "test-service_account"
    tags = ["test-tags"]
  }
  node_count = 42
  version = "test-version"
}
