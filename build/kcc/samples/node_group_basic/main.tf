data "google_compute_node_types" "central1a" {
  zone = "us-central1-a"
}

resource "google_compute_node_template" "soletenant-tmpl" {
  name = "soletenant-tmpl"
  region = "us-central1"
  node_type = "${data.google_compute_node_types.central1a.names[0]}"
}

resource "google_compute_node_group" "nodes" {
  name = "soletenant-group"
  zone = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  size = 1
  node_template = "${google_compute_node_template.soletenant-tmpl.self_link}"
}
