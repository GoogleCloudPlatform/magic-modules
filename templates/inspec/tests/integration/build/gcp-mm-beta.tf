provider "google-beta" {
  region  = "${var.gcp_location}",
  project = "${var.gcp_project_id}"
}

variable "subnetwork_policy_binding" {
  type = "map"
}

resource "google_compute_subnetwork" "inspec-gcp-subnetwork-iam" {
  name          = "${var.subnetwork_policy_binding["subnet_name"]}"
  ip_cidr_range = "10.2.0.0/16"
  network       = "${google_compute_network.inspec-gcp-network-iam.self_link}"
  project       = "${var.gcp_project_id}"
}

resource "google_compute_network" "inspec-gcp-network-iam" {
  project                 = "${var.gcp_project_id}"
  name                    = "test-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork_iam_binding" "inspec-gcp-subnet" {
  subnetwork = "${google_compute_subnetwork.inspec-gcp-subnetwork-iam.name}"
  role       = "${var.subnetwork_policy_binding["role"]}"
	project    = "${var.gcp_project_id}"
  
  members = [
    "${var.subnetwork_policy_binding["member"]}"
  ]
}