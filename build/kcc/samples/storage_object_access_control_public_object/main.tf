resource "google_storage_object_access_control" "public_rule" {
  object = "${google_storage_bucket_object.object.output_name}"
  bucket = "${google_storage_bucket.bucket.name}"
  role   = "READER"
  entity = "allUsers"
}

resource "google_storage_bucket" "bucket" {
	name = "static-content-bucket-${local.name_suffix}"
}

 resource "google_storage_bucket_object" "object" {
	name   = "public-object-${local.name_suffix}"
	bucket = "${google_storage_bucket.bucket.name}"
	source = "../static/header-logo.png"
}
