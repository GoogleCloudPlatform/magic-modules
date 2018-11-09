variable "project_name" {}
variable "zone" {}
variable "region" {}
variable "network" {
  type = "map"
}

variable "subnetwork" {
  type = "map"
}

provider "google" {
  project = "${var.project_name}"
  region = "${var.region}"
  zone = "${var.zone}"  
}


resource "google_service_account" "inspecaccount" {
  account_id = "inspec-account"
  display_name = "InSpec Service Account"
}

resource "google_service_account_key" "inspeckey" {
  service_account_id = "${google_service_account.inspecaccount.name}"
  public_key_type = "TYPE_X509_PEM_FILE"
}

resource "google_project_iam_member" "inspec-iam-member" {
  role = "roles/viewer"
  member = "serviceAccount:${google_service_account.inspecaccount.email}"
}

resource "local_file" "file" {
  content = "${base64decode(google_service_account_key.inspeckey.private_key)}"
  filename = "${path.module}/inspec.json"
}

# Network
resource "google_compute_network" "inspec-gcp-network" {
  name = "${var.network["name"]}"
  auto_create_subnetworks = "false"
  routing_mode = "${var.network["routing_mode"]}"
}

# Subnetwork
resource "google_compute_subnetwork" "inspec-gcp-subnetwork" {
  ip_cidr_range = "${var.subnetwork["ip_range"]}"
  name =  "${var.subnetwork["name"]}"
  network = "${google_compute_network.inspec-gcp-network.self_link}"
}