resource "google_ml_engine_model" "default" {
  name = "default"
  description = "My model"
  regions = ["us-central1"]
  labels  = {
    my_model = "foo"
  }
  online_prediction_logging = true
  online_prediction_console_logging = true
}
