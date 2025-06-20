/**
 * Copyright 2021 Google LLC
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

resource "google_spanner_instance" "spanner-instance" {
  name         = "spanner-instance"
  display_name = "spanner-instance"
  config       = "regional-us-central1"
  num_nodes    = 1
}

resource "google_spanner_instance_iam_binding" "instance" {
  instance = google_spanner_instance.spanner-instance.name
  role     = "roles/compute.networkUser"

  members = [
    "user:jane@example.com",
  ]
}
