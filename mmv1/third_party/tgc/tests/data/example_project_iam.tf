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

resource "google_project_iam_member" "owner-a" {
  project = "{{.Provider.project}}"
  role    = "roles/owner"
  member  = "user:example-a@google.com"
}

resource "google_project_iam_member" "viewer-a" {
  project = "{{.Provider.project}}"
  role    = "roles/viewer"
  member  = "user:example-a@google.com"
}

resource "google_project_iam_member" "viewer-b" {
  project = "{{.Provider.project}}"
  role    = "roles/viewer"
  member  = "user:example-b@google.com"
}

resource "google_project_iam_binding" "editors" {
  project = "{{.Provider.project}}"
  role    = "roles/editor"
  members  = [
    "user:example-a@google.com",
    "user:example-b@google.com"
  ]
}

resource "google_project_iam_binding" "storage" {
  project = "{{.Provider.project}}"
  role    = "roles/storage.admin"
  members  = [
    "user:example-a@google.com",
    "user:example-b@google.com"
  ]
}
