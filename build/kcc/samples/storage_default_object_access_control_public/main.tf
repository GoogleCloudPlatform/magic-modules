resource "google_storage_default_object_access_control" "public_rule" {
  bucket = "${google_storage_bucket.bucket.name}"
  role   = "READER"
  entity = "allUsers"
}

resource "google_storage_bucket" "bucket" {
	name = "static-content-bucket"
}
