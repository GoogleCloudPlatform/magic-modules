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

resource "google_bigquery_dataset" "example-dataset" {
  dataset_id                  = "test_dataset"
  location                    = "EU"
  project                     = "{{.Provider.project}}"
  default_table_expiration_ms = 3600000

  labels = {
    env = "dev"
  }

}

resource "google_bigquery_dataset_iam_policy" "dataset" {
  dataset_id  = google_bigquery_dataset.example-dataset.dataset_id
  policy_data = "{\"bindings\":[{\"members\":[\"serviceAccount:998476993360@cloudservices.gserviceaccount.com\",\"allAuthenticatedUsers\"],\"role\":\"roles/editor\"}]}"
}
