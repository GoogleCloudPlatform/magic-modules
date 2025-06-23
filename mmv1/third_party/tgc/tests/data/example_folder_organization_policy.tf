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

resource "google_folder_organization_policy" "serial_port_policy" {
  folder     = "folders/{{.FolderID}}"
  constraint = "compute.disableSerialPortAccess"

  boolean_policy {
    enforced = true
  }
}

resource "google_folder_organization_policy" "services_policy_a" {
  folder     = "folders/{{.FolderID}}"
  constraint = "serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}

resource "google_folder_organization_policy" "services_policy_b" {
  folder     = "folders/{{.FolderID}}"
  constraint = "serviceuser.services"

  list_policy {
    suggested_value = "compute.googleapis.com"

    deny {
      values = ["cloudresourcemanager.googleapis.com"]
    }
  }
}

resource "google_folder_organization_policy" "services_policy_c" {
  folder     = "folders/{{.FolderID}}"
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}
