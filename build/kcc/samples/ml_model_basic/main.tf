resource "google_ml_engine_model" "default" {
  name = "default-${local.name_suffix}"
  description = "My model"
  regions = ["us-central1"]
}
