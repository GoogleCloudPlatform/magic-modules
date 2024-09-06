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

resource "google_project_organization_policy" "serial_port_policy" {
  project    = "{{.Provider.project}}"
  constraint = "compute.disableSerialPortAccess"

  boolean_policy {
    enforced = true
  }
}

resource "google_project_organization_policy" "services_policy" {
  project    = "{{.Provider.project}}"
  constraint = "serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
resource "google_project_organization_policy" "service_account_policy" {
  project    = "{{.Provider.project}}"
  constraint = "iam.disableServiceAccountCreation"

  restore_policy {
    default = true
  }
}
