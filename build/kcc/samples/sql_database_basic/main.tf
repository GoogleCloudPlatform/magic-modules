resource "google_sql_database" "database" {
	name = "my-database"
	instance = "${google_sql_database_instance.instance.name}"
}

resource "google_sql_database_instance" "instance" {
	name = "my-database-instance"
	region = "us-central"
	settings {
		tier = "D0"
	}
}
