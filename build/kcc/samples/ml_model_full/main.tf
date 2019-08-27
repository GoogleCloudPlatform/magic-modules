resource "google_ml_engine_model" "default" {
  name = "default-${local.name_suffix}"
  description = "My model"
  regions = ["us-central1"]
  labels  = {
    my_model = "foo"
  }
  online_prediction_logging = true
  online_prediction_console_logging = true
}
