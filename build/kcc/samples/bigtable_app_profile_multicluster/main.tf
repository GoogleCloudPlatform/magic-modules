resource "google_bigtable_instance" "instance" {
	name = "tf-test-instance--${local.name_suffix}"
	cluster {
		cluster_id   = "tf-test-instance--${local.name_suffix}"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
}

resource "google_bigtable_app_profile" "ap" {
	instance = google_bigtable_instance.instance.name
	app_profile_id = "tf-test-profile--${local.name_suffix}"

	multi_cluster_routing_use_any = true
	ignore_warnings = true
}
