package biglake_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBiglakeTable_bigqueryBiglakeTable_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeTable_bigqueryBiglakeTableExample(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "database", "catalog"},
			},
			{
				Config: testAccBiglakeTable_bigqueryBiglakeTable_update(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "database", "catalog"},
			},
		},
	})
}

func testAccBiglakeTable_bigqueryBiglakeTable_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_biglake_catalog" "catalog" {
	name = "<%= ctx[:vars]['catalog'] %>"
	location = "US"
}
resource "google_storage_bucket" "bucket" {
	name                        = "<%= ctx[:vars]['bucket'] %>"
	location                    = "US"
	force_destroy               = true
	uniform_bucket_level_access = true
}
resource "google_storage_bucket_object" "metadata_folder" {
	name    = "metadata/"
	content = " "
	bucket  = google_storage_bucket.bucket.name
}
resource "google_storage_bucket_object" "data_folder" {
	name    = "data/"
	content = " "
	bucket  = google_storage_bucket.bucket.name
}
resource "google_biglake_database" "database" {
	name = "<%= ctx[:vars]['database'] %>"
	catalog = google_biglake_catalog.catalog.id
	location = google_biglake_catalog.catalog.location
	type = "HIVE"
	hive_options {
		location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.metadata_folder.name}"
		parameters = {
			"name" = "wrench"
		}
	}
}
resource "google_biglake_table" "table" {
    name = "tf_test_my_table%{random_suffix}"
    database = google_biglake_database.database.id
    location = google_biglake_catalog.catalog.location
    type = "HIVE"
    hive_options {
		table_type = "MANAGED_TABLE"
		storage_descriptor {
		  location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.data_folder.name}"
		  input_format = "org.apache.hadoop.mapred.SequenceFileInputFormat",
		  output_format =  "org.apache.hadoop.hive.ql.io.HiveSequenceFileOutputFormat"
		}
		# Some Example Parameters.
		parameters = {
		  "spark.sql.create.version" = "3.1.7"
		  "spark.sql.sources.schema.numParts" = "1"
		  "transient_lastDdlTime" = "1680895000"
		  "spark.sql.partitionProvider" = "catalog"
		  "owner" = "Jane Doe"
		  "spark.sql.sources.schema.part.0" = "{\"type\":\"struct\",\"fields\":[{\"name\":\"id\",\"type\":\"integer\",\"nullable\":true,\"metadata\":{}},{\"name\":\"name\",\"type\":\"string\",\"nullable\":true,\"metadata\":{}},{\"name\":\"age\",\"type\":\"integer\",\"nullable\":true,\"metadata\":{}}]}"
		  "spark.sql.sources.provider": "iceberg"
		  "provider" = "iceberg"
		}
	}
}
`, context)
}
