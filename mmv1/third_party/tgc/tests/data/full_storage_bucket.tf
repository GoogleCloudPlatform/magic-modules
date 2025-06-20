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

resource "google_storage_bucket" "full-list-default" {
  name     = "image-store-bucket"
  location = "EU"

  uniform_bucket_level_access = true
  cors {
    origin          = ["test-origin1", "test-origin2"]
    method          = ["test-method1", "test-method2"]
    response_header = ["test-response_header1", "test-response_header2"]
    max_age_seconds = 42
  }
  cors {
    origin          = ["test-origin1", "test-origin2"]
    method          = ["test-method1", "test-method2"]
    response_header = ["test-response_header1", "test-response_header2"]
    max_age_seconds = 42
  }
  encryption {
    default_kms_key_name = "test-default_kms_key_name"
  }
  force_destroy = true
  labels = {
    label_foo1 = "label-bar1"
  }
  lifecycle_rule {
    action {
      type          = "test-type"
      storage_class = "test-storage_class"
    }
    condition {
      age                   = 42
      created_before        = "test-created_before"
      matches_storage_class = ["test-matches_storage_class1", "matches_storage_class2"]
      num_newer_versions    = 42
      with_state            = "LIVE"
    }
  }
  lifecycle_rule {
    action {
      type          = "test-type"
      storage_class = "test-storage_class"
    }
    condition {
      age                   = 42
      created_before        = "test-created_before"
      matches_storage_class = ["test-matches_storage_class1", "matches_storage_class2"]
      num_newer_versions    = 42
      with_state            = "LIVE"
    }
  }
  logging {
    log_bucket        = "test-log_bucket"
    log_object_prefix = "test-log_object_prefix"
  }
  requester_pays = true
  retention_policy {
    is_locked        = true
    retention_period = 42
  }
  storage_class = "test-storage_class"
  versioning {
    enabled = true
  }
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
