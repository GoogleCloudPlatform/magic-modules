# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

terraform {
required_providers {
google = {
source = "hashicorp/google"
version = ">= 4.34.0"
}
}
}

provider "google" {
project = "{{.Provider.project}}"
}

resource "google_kms_key_ring" "gg_asset_key_ring_43576_f7a1" {
name = "gg-asset-key-ring-43576-f7a1"
location = "global"
project = "{{.Provider.project}}"
}

resource "google_kms_crypto_key" "gg_asset_crypto_key_43576_f7a1" {
name = "gg-asset-crypto-key-43576-f7a1"
key_ring = google_kms_key_ring.gg_asset_key_ring_43576_f7a1.id
}

resource "google_kms_crypto_key_version" "gg_asset_crypto_key_version_43576_f7a1" {
crypto_key = google_kms_crypto_key.gg_asset_crypto_key_43576_f7a1.id
state = "ENABLED"
}
