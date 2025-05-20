package bigquery_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleBigqueryTable_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	expectedID := fmt.Sprintf("projects/%s/datasets/%s/tables/%s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("tf_test_ds_%s", context["random_suffix"]), fmt.Sprintf("tf_test_table_%s", context["random_suffix"]))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryTable_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "table_id", fmt.Sprintf("tf_test_table_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "dataset_id", fmt.Sprintf("tf_test_ds_%s", context["random_suffix"])),
					resource.TestCheckResourceAttrSet("data.google_bigquery_table.example", "schema.#"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.0.name", "name"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.0.type", "STRING"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.0.mode", "NULLABLE"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.2.name", "address"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.2.fields.1.name", "zip"),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "id", expectedID),
					resource.TestCheckResourceAttr("data.google_bigquery_table.example", "schema.0.fields.4.name", "policy_tag_test"),
					resource.TestMatchResourceAttr("data.google_bigquery_table.example", "schema.0.fields.4.policy_tags.0.names.0", regexp.MustCompile("^projects/[^/]+/locations/us-central1/taxonomies/[^/]+/policyTags/[^/]+$")),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryTable_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_data_catalog_policy_tag" "test" {
    taxonomy     = google_data_catalog_taxonomy.test.id
    display_name = "Low security"
    description  = "A policy tag normally associated with low security items"
  }
  
  resource "google_data_catalog_taxonomy" "test" {
    region                 = "us-central1"
    display_name           = "taxonomy_%{random_suffix}"
    description            = "A collection of policy tags"
    activated_policy_types = ["FINE_GRAINED_ACCESS_CONTROL"]
  }

  resource "google_bigquery_dataset" "test" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "us-central1"
    default_table_expiration_ms = 3600000
  }

  resource "google_bigquery_table" "test" {
    dataset_id          = google_bigquery_dataset.test.dataset_id
    table_id            = "tf_test_table_%{random_suffix}"
    deletion_protection = false
    depends_on          = [google_data_catalog_policy_tag.test]
    schema              = <<EOF
    [
      {
        "name": "name",
        "type": "STRING",
        "mode": "NULLABLE"
      },
	  {
		"name": "age",
		"type": "INTEGER",
		"mode": "NULLABLE",
		"description": "Age of the person"
	  },
	  {
		"name": "address",
		"type": "RECORD",
		"mode": "NULLABLE",
		"fields": [
		  {
			"name": "street",
			"type": "STRING",
			"mode": "NULLABLE"
		  },
		  {
			"name": "zip",
			"type": "STRING",
			"mode": "NULLABLE"
		  },
		  {
			"name": "city",
			"type": "STRING",
			"mode": "NULLABLE"
		  }
		],
		"description": "Address of the person"
      },
	  {
		"name": "phone_numbers",
		"type": "RECORD",
		"mode": "REPEATED",
		"fields": [
		  {
			"name": "type",
			"type": "STRING",
			"mode": "NULLABLE"
		  },
		  {
			"name": "number",
			"type": "STRING",
			"mode": "NULLABLE"
		  }
		],
		"description": "Phone numbers of the person"
	  },
	  {
		"name": "policy_tag_test",
		"type": "STRING",
		"mode": "NULLABLE",
		"description": "A test field with policy tags",
		"policyTags": {
		  "names": [
			"${google_data_catalog_policy_tag.test.id}"
          ]
		}
      }
    ]
    EOF
  }

  data "google_bigquery_table" "example" {
    dataset_id = google_bigquery_table.test.dataset_id
	table_id   = google_bigquery_table.test.table_id
  }
`, context)
}
