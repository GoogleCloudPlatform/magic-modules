resource "google_vmwareengine_external_address" "main" {
  name        = "gg-asset-ext-addr-03361-811b"
  parent      = "projects/{{.Provider.project}}/locations/us-central1-a/privateClouds/gg-asset-pc-03361-811b"
  internal_ip = "10.100.0.10"
  description = "External address for testing"
}