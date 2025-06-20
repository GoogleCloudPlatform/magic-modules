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

resource "google_storage_bucket" "default" {
  name     = "fake-bucket-123456"
  location = "EU"
  uniform_bucket_level_access = true
  project = "{{.Provider.project}}"
}

resource "google_storage_bucket_iam_policy" "project" {
  bucket = google_storage_bucket.default.name
  policy_data = "{\"bindings\":[{\"members\":[\"user:jane@example.com\"],\"role\":\"roles/storage.admin\"}]}"
}
