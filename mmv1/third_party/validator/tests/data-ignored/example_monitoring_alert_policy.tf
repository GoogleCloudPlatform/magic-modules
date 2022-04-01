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
      source  = "hashicorp/google"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials}}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_monitoring_alert_policy" "alert_policy" {
  project = "{{.Provider.project}}"
  display_name = "My Alert Policy"
  combiner = "OR"
  conditions {
    display_name = "test condition"
    condition_threshold {
      filter = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration = "60s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }

  notification_channels = [
    "projects/{{.Provider.project}}/notificationChannels/channel-1",
    "projects/{{.Provider.project}}/notificationChannels/channel-2",
    "projects/{{.Provider.project}}/notificationChannels/channel-3"
  ]
}