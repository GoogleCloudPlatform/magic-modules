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
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iterator"
)

func TestAccSecurityCenterV2ProjectBigQueryExportConfig_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	datasetID := "tf_test_" + randomSuffix
	orgID := envvar.GetTestOrgFromEnv(t)

	// Setup test context
	ctx := context.Background()
	projectID := envvar.GetTestProjectFromEnv()
	location := "US"

	// Run cleanup before the test starts
	if err := cleanupProjectBigQueryDatasets(ctx, projectID, "tf_test_", location); err != nil {
		t.Fatalf("Failed to clean up BigQuery datasets: %v", err)
	}

	context := map[string]interface{}{
		"org_id":              orgID,
		"random_suffix":       randomSuffix,
		"dataset_id":          datasetID,
		"big_query_export_id": "tf-test-export-" + randomSuffix,
		"name": fmt.Sprintf("projects/%s/locations/global/bigQueryExports/%s",
			envvar.GetTestProjectFromEnv(), "tf-test-export-"+randomSuffix),
		"project": envvar.GetTestProjectFromEnv(),
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
				Config:  testAccSecurityCenterV2ProjectBigQueryExportConfig_basic(context),
				Destroy: true,
			},
			{
				ResourceName:            "google_scc_v2_project_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "project"},
			},
			{
				Config:  testAccSecurityCenterV2ProjectBigQueryExportConfig_update(context),
				Destroy: true,
			},
			{
				ResourceName:            "google_scc_v2_project_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "project"},
			},
		},
	})
}

func cleanupProjectBigQueryDatasets(ctx context.Context, projectID, prefix, location string) error {

	service, err := bigquery.NewService(ctx)
	if err != nil {
		return err
	}

	listCall := service.Datasets.List(projectID)
	response, err := listCall.Do()

	if err != nil {
		return err
	}

	for _, dataset := range response.Datasets {
		datasetID := dataset.DatasetReference.DatasetId

		if strings.HasPrefix(datasetID, prefix) {
			log.Printf("Deleting dataset with ID: %s", datasetID)
			deleteCall := service.Datasets.Delete(projectID, datasetID).DeleteContents(true)
			if err := deleteCall.Context(ctx).Do(); err != nil {
				log.Printf("Failed to delete dataset %s: %v", datasetID, err)
			}
		}
	}

	return nil
}

func testAccSecurityCenterV2ProjectBigQueryExportConfig_basic(context map[string]interface{}) string {
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
	# need to wait for destruction due to 
	# 'still in use' error from api 
	destroy_duration = "1m"
}

resource "google_scc_v2_project_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  project      = "%{project}"
  dataset      = google_bigquery_dataset.default.id
  location     = "global"
  description  = "Cloud Security Command Center Findings Big Query Export Config"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_1_minute]
}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "3m"
	depends_on = [google_scc_v2_project_scc_big_query_export.default]
}

`, context)
}

func testAccSecurityCenterV2ProjectBigQueryExportConfig_update(context map[string]interface{}) string {
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
	# need to wait for destruction due to 
	# 'still in use' error from api 
	destroy_duration = "1m"
}

resource "google_scc_v2_project_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  project      = "%{project}"
  dataset      = google_bigquery_dataset.default.id
  location     = "global"
  description  = "SCC Findings Big Query Export Update"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_1_minute] 

}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "3m"
	depends_on = [google_scc_v2_project_scc_big_query_export.default]
}

`, context)
}
