resource "google_project" "example-project" {
  billing_account = "example-account"
  folder_id       = "456"

  labels = {
    project-label-key-a = "project-label-val-a"
  }

  name       = "My Project"
  project_id = "example-project"
}
