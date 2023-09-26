resource "google_project_iam_policy" "example-project_iam_policy" {
  policy_data = "{\"bindings\":[{\"role\":\"roles/editor\",\"members\":[\"user:example-a@google.com\",\"user:example-b@google.com\"]},{\"role\":\"roles/storage.admin\",\"members\":[\"user:example-a@google.com\",\"user:example-b@google.com\"]},{\"role\":\"roles/owner\",\"members\":[\"user:example-a@google.com\"]},{\"role\":\"roles/viewer\",\"members\":[\"user:example-a@google.com\",\"user:example-b@google.com\"]}]}"
  project     = "example-project"
}
