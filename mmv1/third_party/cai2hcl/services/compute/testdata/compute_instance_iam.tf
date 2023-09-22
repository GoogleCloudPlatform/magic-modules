resource "google_compute_instance_iam_policy" "example_instance_iam_policy" {
  instance_name = "example_instance"
  policy_data   = "{\"bindings\":[{\"role\":\"roles/compute.osLogin\",\"members\":[\"user:jane@example.com\"]}]}"
  project       = "test-project"
  zone          = "example_zone"
}
