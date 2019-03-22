provider "google-beta" {
  region      = "${var.gcp_location}",
  project         = "${var.gcp_project_id}"
}

resource "google_compute_subnetwork" "subnetwork-iam" {
  name          = "test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  network       = "${google_compute_network.custom-test.self_link}"
  project         = "${var.gcp_project_id}"
}

resource "google_compute_network" "custom-test" {
  name                    = "test-network"
  auto_create_subnetworks = false
  project         = "${var.gcp_project_id}"
}

resource "google_compute_subnetwork_iam_binding" "subnet" {
  subnetwork = "${google_compute_subnetwork.subnetwork-iam.name}"
  role       = "roles/compute.networkUser"
	project         = "${var.gcp_project_id}"
  
  members = [
    "user:slevenick@google.com",
  ]
}