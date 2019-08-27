data "google_compute_node_types" "central1a" {
  zone = "us-central1-a"
}

resource "google_compute_node_template" "template" {
  name = "soletenant-tmpl-${local.name_suffix}"
  region = "us-central1"
  node_type = "${data.google_compute_node_types.central1a.names[0]}"
}
