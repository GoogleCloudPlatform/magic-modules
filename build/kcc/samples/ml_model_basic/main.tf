resource "google_ml_engine_model" "default" {
  name = "default"
  description = "My model"
  regions = ["us-central1"]
}
