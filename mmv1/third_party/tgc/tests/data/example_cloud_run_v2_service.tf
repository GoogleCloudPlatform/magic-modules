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

resource "google_cloud_run_v2_service" "default" {
  name     = "cloudrunv2-to-get-cai"
  location = "us-central1"
  project  = "{{.Provider.project}}"

  annotations = {
    "generated-by" = "magic-modules"
  }

  template {
    max_instance_request_concurrency = 10
    timeout       = "600s"

    containers {
      image = "gcr.io/cloudrun/hello"
      args  = ["arrgs"]
      ports {
        container_port = 8080
      }
    }
  }

  traffic {
    percent         = 100
    type            = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }

  lifecycle {
    ignore_changes = [
      annotations,
    ]
  }
}
