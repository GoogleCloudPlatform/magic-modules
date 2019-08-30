resource "google_spanner_instance" "main" {
  config       = "regional-europe-west1"
  display_name = "main-instance"
}

resource "google_spanner_database" "database" {
  instance  = "${google_spanner_instance.main.name}"
  name      = "my-database-${local.name_suffix}"
  ddl       =  [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)"
  ]
}
