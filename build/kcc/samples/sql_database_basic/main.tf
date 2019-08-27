resource "google_sql_database" "database" {
	name = "my-database-${local.name_suffix}"
	instance = "${google_sql_database_instance.instance.name}"
}

resource "google_sql_database_instance" "instance" {
	name = "my-database-instance-${local.name_suffix}"
	region = "us-central"
	settings {
		tier = "D0"
	}
}
