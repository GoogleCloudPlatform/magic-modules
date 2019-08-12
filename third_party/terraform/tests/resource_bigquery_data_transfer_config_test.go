package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBigqueryDataTransferConfig_scheduledQueryUpdate(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(10)

	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, projectOrg, projectBillingAccount, "first", "y"),
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, projectOrg, projectBillingAccount, "second", "z"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, org, billing, schedule, letter string) string {
	return fmt.Sprintf(`
resource "google_project" "update" {
  name = "terraform-%s"
  project_id = "terraform-%s"
  org_id = "%s"
  billing_account = "%s"
}

resource "google_project_services" "update" {
  project = google_project.update.project_id
  services = ["bigquerydatatransfer.googleapis.com", "bigquery-json.googleapis.com"]
}

resource "google_project_iam_member" "permissions" {
  depends_on = [google_project_services.update]
  project = google_project.update.project_id
  role = "roles/iam.serviceAccountShortTermTokenMinter"
  member = "serviceAccount:service-${google_project.update.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  project = google_project.update.project_id

  depends_on = [google_project_iam_member.permissions]

  display_name = "my-query-%s"
  location = "asia-northeast1"
  data_source_id = "scheduled_query"
  schedule = "%s sunday of quarter 00:00"
  destination_dataset_id = google_bigquery_dataset.my-dataset.dataset_id
  params = {
    destination_table_name_template = "my-table"
    write_disposition = "WRITE_APPEND"
    query = "SELECT name FROM tabl WHERE x = '%s'"
  }
}

resource "google_bigquery_dataset" "my-dataset" {
  project = google_project.update.project_id

  depends_on = [google_project_iam_member.permissions]

  dataset_id = "my_dataset%s"
  friendly_name = "foo"
  description = "bar"
  location = "asia-northeast1"
}
`, random_suffix, random_suffix, org, billing, random_suffix, schedule, letter, random_suffix)
}
