/**
 * Copyright 2022 Google LLC
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

resource "random_string" "suffix" {
  length  = 4
  upper   = false
  special = false
}


data "google_iam_policy" "owner" {
  binding {
    role = "roles/bigquery.dataOwner"

    members = [
      "${random_string.suffix.result}:jane@example.com",
    ]
  }
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "example_dataset"
}

resource "google_bigquery_dataset_iam_policy" "dataset" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
  policy_data = data.google_iam_policy.owner.policy_data
}
