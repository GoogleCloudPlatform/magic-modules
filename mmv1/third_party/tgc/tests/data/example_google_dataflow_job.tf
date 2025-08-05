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

resource "google_dataflow_job" "word_count_job" {
  name = "my-word-count-job"
  template_gcs_path = "gs://dataflow-templates/latest/Word_Count"
  parameters = {
    inputFile = "gs://dataflow_job1/data.txt"
    output = "gs://dataflow_job1/wordcount_output"
  }
  temp_gcs_location = "gs://dataflow_job1/tmp"
}