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
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google-beta" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_org_policy_custom_constraint" "gke_auto_upgrade_constraint" {
 
  provider       = google-beta
  name           = "custom.disableGkeAutoUpgrade"
  parent         = "organizations/{{.OrgID}}"
  display_name   = "Disable GKE auto upgrade"
  description    = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."
  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE", "UPDATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}

resource "google_org_policy_custom_constraint" "dataprocAmPrimaryOnlyEnforced" {
  provider      = google-beta
  name          = "custom.dataprocAmPrimaryOnlyEnforced"
  parent        = "organizations/{{.OrgID}}"
  display_name  = "Application master cannot run on preemptible workers"
  description   = "Property \"dataproc:am.primary_only\" must be \"true\"."
  action_type    = "ALLOW"
  condition      = "(\"dataproc:am.primary_only\" in resource.config.softwareConfig.properties) && (resource.config.softwareConfig.properties[\"dataproc:am.primary_only\"]==\"true\")"
  method_types   = ["CREATE"]
  resource_types = ["dataproc.googleapis.com/Cluster"]
}