package storageinsights_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccStorageInsightsReportConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageInsightsReportConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsReportConfig_full(context),
			},
			{
				ResourceName:            "google_storage_insights_report_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		        {
				Config: testAccStorageInsightsReportConfig_update(context),
			},
			{
				ResourceName:            "google_storage_insights_report_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},

		},
	})
}

func testAccStorageInsightsReportConfig_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_storage_insights_report_config" "config" {
  display_name = "Test Report Config"
  location = "us-central1"
  frequency_options {
    frequency = "WEEKLY"
    start_date {
      day = 15
      month = 3
      year = 2050
    }
    end_date {
      day = 15
      month = 4
      year = 2050
    }
  }
  csv_options {
    record_separator = "\n"
    delimiter = ","
    header_required = false
  }
  object_metadata_report_options {
    metadata_fields = ["bucket", "name", "project"]
    storage_filters {
      bucket = google_storage_bucket.report_bucket.name
    }
    storage_destination_options {
      bucket = google_storage_bucket.report_bucket.name
      destination_path = "test-insights-reports"
    }
  }
}

resource "google_storage_bucket" "report_bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.report_bucket.name
  role   = "roles/storage.admin"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-storageinsights.iam.gserviceaccount.com"
}
`, context)
}

func testAccStorageInsightsReportConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_storage_insights_report_config" "config" {
  display_name = "Test Report Config Updated"
  location = "us-central1"
  frequency_options {
    frequency = "MONTHLY"
    start_date {
      day = 14
      month = 3
      year = 2040
    }
    end_date {
      day = 14
      month = 4
      year = 2040
    }
  }
  csv_options {
    record_separator = ";"
    delimiter = "."
    header_required = true
  }
  object_metadata_report_options {
    metadata_fields = ["name", "project"]
    storage_filters {
	  bucket = google_storage_bucket.update_bucket.name
	}
    storage_destination_options {
      bucket = google_storage_bucket.update_bucket.name
      destination_path = "test-insights-reports-updated"
    }
  }
}

resource "google_storage_bucket" "update_bucket" {
	name                        = "tf-test-update-bucket%{random_suffix}"
	location                    = "us-central1"
	force_destroy               = true
	uniform_bucket_level_access = true
  }

resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.report_bucket.name
  role   = "roles/storage.admin"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-storageinsights.iam.gserviceaccount.com"
}
`, context)
}
