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
      source = "hashicorp/google"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_pubsub_subscription" "example" {
  name  = "example-subscription"
  topic = "example-pubsub-topic"

  ack_deadline_seconds = 20

  labels = {
    test-label1 = "test-value1"
  }

  push_config {
    push_endpoint = "https://example.com/push"

    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription_iam_policy" "editor" {
  subscription = google_pubsub_subscription.example.name
  policy_data = jsonencode(
    {
      bindings = [
        {
          members = [
            "user:jane@example.com",
          ]
          role = "roles/editor"
        }
      ]
    }
  )
}

