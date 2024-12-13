package securitycenterv2_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	bigquery "google.golang.org/api/bigquery/v2"
)

func TestAccSecurityCenterV2OrganizationBigQueryExportsConfig_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	datasetID := "tf_test_" + randomSuffix
	orgID := envvar.GetTestOrgFromEnv(t)

	// Run cleanup before the test starts
	ctx := context.Background()
	projectID := envvar.GetTestProjectFromEnv()
	err := cleanupOrganizationBigQueryExportsDatasets(ctx, "tf_test_", projectID)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	context := map[string]interface{}{
		"org_id":              orgID,
		"random_suffix":       randomSuffix,
		"dataset_id":          datasetID,
		"big_query_export_id": "tf-test-export-" + randomSuffix,
		"name": fmt.Sprintf("organizations/%s/locations/global/bigQueryExports/%s",
			orgID, "tf-test-export-"+randomSuffix),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccSecurityCenterV2OrganizationBigQueryExportsConfig_basic(context),
				Destroy: true,
			},
			{
				ResourceName:            "google_scc_v2_organization_scc_big_query_exports.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
			{
				Config:  testAccSecurityCenterV2OrganizationBigQueryExportsConfig_update(context),
				Destroy: true,
			},
			{
				ResourceName:            "google_scc_v2_organization_scc_big_query_exports.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
		},
	})
}

func cleanupOrganizationBigQueryExportsDatasets(ctx context.Context, prefix string, projectID string) error {

	service, err := bigquery.NewService(ctx)

	if err != nil {
		return fmt.Errorf("failed to create BigQuery service: %v", err)
	}

	datasetsService := bigquery.NewDatasetsService(service)
	datasetsListCall := datasetsService.List(projectID)
	datasets, err := datasetsListCall.Do()

	if err != nil {
		return fmt.Errorf("failed to list datasets: %v", err)
	}

	for _, dataset := range datasets.Datasets {

		if strings.HasPrefix(dataset.Id, prefix) {

			log.Printf("Deleting dataset with ID: %s", dataset.Id)

			err := datasetsService.Delete(projectID, dataset.Id).DeleteContents(true).Do()

			if err != nil {
				return fmt.Errorf("failed to delete dataset %s: %v", dataset.Id, err)
			}
		}
	}
	return nil
}

func testAccSecurityCenterV2OrganizationBigQueryExportsConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  delete_contents_on_destroy  = true

  labels = {
    env = "default"
  }

  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}

resource "time_sleep" "wait_1_minute" {
	depends_on = [google_bigquery_dataset.default]
	create_duration = "6m"
}

resource "google_scc_v2_organization_scc_big_query_exports" "default" {
  name		   = "%{name}"
  big_query_export_id    = "%{big_query_export_id}"
  organization = "%{org_id}"
  dataset      = google_bigquery_dataset.default.id
  location     = "global"
  description  = "Cloud Security Command Center Findings Big Query Export Config"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_1_minute]
}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "3m"
	depends_on = [google_scc_v2_organization_scc_big_query_exports.default]
}
`, context)
}

func testAccSecurityCenterV2OrganizationBigQueryExportsConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  delete_contents_on_destroy  = true

  labels = {
    env = "default"
  }

  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}

resource "time_sleep" "wait_1_minute" {
	depends_on = [google_bigquery_dataset.default]
	create_duration = "6m"
}

resource "google_scc_v2_organization_scc_big_query_exports" "default" {
  name		   = "%{name}"
  big_query_export_id    = "%{big_query_export_id}"
  organization = "%{org_id}"
  dataset      = google_bigquery_dataset.default.id
  location     = "global"
  description  = "SCC Findings Big Query Export Update"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_1_minute]  
}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "3m"
	depends_on = [google_scc_v2_organization_scc_big_query_exports.default]
}
`, context)
}
