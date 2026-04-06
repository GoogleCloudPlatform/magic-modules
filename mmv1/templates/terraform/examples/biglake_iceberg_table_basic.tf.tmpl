resource "google_storage_bucket" "bucket" {
  name          = "my-bucket-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_table" "my_iceberg_table" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "my_table_%{random_suffix}"
  location  = "gs://${google_storage_bucket.bucket.name}/${google_biglake_iceberg_namespace.namespace.namespace_id}/my_table_%{random_suffix}"
  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
      doc      = "The ID of the record"
    }
    fields {
      id       = 2
      name     = "name"
      type     = "string"
      required = false
    }
    identifier_field_ids = [1]
  }
  partition_spec {
    fields {
      name      = "id_partition"
      source_id = 1
      transform = "identity"
    }
  }
}
