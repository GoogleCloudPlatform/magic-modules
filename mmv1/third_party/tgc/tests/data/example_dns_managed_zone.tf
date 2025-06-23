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

resource "google_dns_managed_zone" "zone1" {
  name        = "publiczone"
  dns_name    = "publiczone.gsecurity.net."

  force_destroy = true
  visibility = "public"

  dnssec_config {
    state         = "on"
    kind          = "dns#managedZoneDnsSecConfig"
    non_existence = "nsec3"
    default_key_specs {
      //      key_type = "keySigning"
      key_type   = "zoneSigning" //      ZONE_SIGNING / RSASHA1 / 1024
      key_length = 1024
      kind       = "dns#dnsKeySpec"
      algorithm  = "rsasha1"
      //      ecdsap256sha256 ecdsap384sha384 rsasha1 rsasha256 rsasha512   //algorithm allowed values
    }
  }
}
