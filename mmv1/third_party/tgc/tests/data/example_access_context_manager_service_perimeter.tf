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

resource "google_access_context_manager_service_perimeter" "service-perimeter" {
  parent = "accessPolicies/987654"
  name   = "accessPolicies/987654/servicePerimeters/restrict_storage"
  title  = "restrict_storage"

  status {
    restricted_services = ["storage.googleapis.com", "bigquery.googleapis.com"]
    resources = ["projects/54321", "projects/4321"]

    ingress_policies {
      ingress_from {
        sources {
          access_level = "accessPolicies/987654/accessLevels/restrict_storage"
        }
        identity_type = "ANY_IDENTITY"
      }

      ingress_to {
        resources = ["*"]
        operations {
          service_name = "storage.googleapis.com"
          method_selectors {
            method = "google.storage.objects.create"
          }
        }
      }
    }

    egress_policies {
      egress_from {
        sources {
          access_level = "accessPolicies/987654/accessLevels/restrict_storage"
        }
	source_restriction = "SOURCE_RESTRICTION_ENABLED"
        identity_type = "ANY_USER_ACCOUNT"
      }
    }
  }
}
