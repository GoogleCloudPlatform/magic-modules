data "google_tpu_tensorflow_versions" "available" { }

resource "google_tpu_node" "tpu" {
	name           = "test-tpu-${local.name_suffix}"
	zone           = "us-central1-b"

	accelerator_type   = "v3-8"
	tensorflow_version = "${data.google_tpu_tensorflow_versions.available.versions[0]}"
	cidr_block         = "10.2.0.0/29"
}
