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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, "first", "y"),
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, "second", "z"),
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

func testAccBigqueryDataTransferConfig_scheduledQueryUpdate(random_suffix, schedule, letter string) string {
	return fmt.Sprintf(`
resource "google_bigquery_data_transfer_config" "query_config" {
  display_name = "my-query-%s"
  location = "asia-northeast1"
  data_source_id = "scheduled_query"
  schedule = "%s sunday of quarter 00:00"
  destination_dataset_id = "${google_bigquery_dataset.my-dataset.dataset_id}"
  params = {
    destination_table_name_template = "my-table"
    write_disposition = "WRITE_APPEND"
    query = "SELECT name FROM tabl WHERE x = '%s'"
  }
}

resource "google_bigquery_dataset" "my-dataset" {
  dataset_id = "my_dataset%s"
  friendly_name = "foo"
  description = "bar"
  location = "asia-northeast1"
}
`, random_suffix, schedule, letter, random_suffix)
}
