variable "ssl_policy" {
  type = "map"
}

resource "google_compute_ssl_policy" "custom-ssl-policy" {
  name            = "${var.ssl_policy["name"]}"
  min_tls_version = "${var.ssl_policy["min_tls_version"]}"
  profile         = "${var.ssl_policy["profile"]}"
  custom_features = ["${var.ssl_policy["custom_feature"]}", "${var.ssl_policy["custom_feature2"]}"]
  project = "${var.gcp_project_id}"
}

resource "google_service_account" "inspecaccount" {
  account_id = "inspec-account"
  display_name = "InSpec Service Account"
  project = "${var.gcp_project_id}"
}

resource "google_service_account_key" "inspeckey" {
  service_account_id = "${google_service_account.inspecaccount.name}"
  public_key_type = "TYPE_X509_PEM_FILE"
  project = "${var.gcp_project_id}"
}

resource "google_project_iam_member" "inspec-iam-member" {
  role = "roles/viewer"
  member = "serviceAccount:${google_service_account.inspecaccount.email}"
  project = "${var.gcp_project_id}"
}

resource "local_file" "file" {
  content = "${base64decode(google_service_account_key.inspeckey.private_key)}"
  filename = "${path.module}/inspec.json"
}
