package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChronicleDataTableRows_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"instance_id":   envvar.GetTestChronicleInstanceIdFromEnv(t),
		"data_table_id": "tf_test_terraform_test_rows_" + randomSuffix,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Step 1: Create initial 2 rows
			{
				Config: testAccChronicleDataTableRows_initial(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_chronicle_data_table_rows.example_rows", "rows.#", "2"),
				),
			},
			// Step 2: Import Check
			{
				ResourceName:            "google_chronicle_data_table_rows.example_rows",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_table_id", "instance", "location"},
			},
			// Step 3: Update to 3 rows
			{
				Config: testAccChronicleDataTableRows_updated(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_chronicle_data_table_rows.example_rows", "rows.#", "3"),
				),
			},
			// Step 4: Import Check after update
			{
				ResourceName:            "google_chronicle_data_table_rows.example_rows",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_table_id", "instance", "location"},
			},
		},
	})
}

func testAccChronicleDataTableRows_initial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_data_table" "example_dt" {
  location        = "us"
  instance        = "%{instance_id}"
  data_table_id   = "%{data_table_id}"
  description     = "Sample DataTable for DataTableRows test"
  deletion_policy = "FORCE"
  column_info {
    column_index    = 0
    original_column = "username"
    column_type     = "STRING"
  }
  column_info {
    column_index    = 1
    original_column = "ip_address"
    column_type     = "CIDR"
  }
}

resource "google_chronicle_data_table_rows" "example_rows" {
  location      = "us"
  instance      = "%{instance_id}"
  data_table_id = google_chronicle_data_table.example_dt.data_table_id
  rows {
    values           = ["user1", "192.168.1.1/32"]
    row_time_to_live = "48h"
  }
  rows {
    values           = ["user2", "10.0.0.1/32"]
    row_time_to_live = "48h"
  }
}
`, context)
}

func testAccChronicleDataTableRows_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_data_table" "example_dt" {
  location        = "us"
  instance        = "%{instance_id}"
  data_table_id   = "%{data_table_id}"
  description     = "Sample DataTable for DataTableRows test"
  deletion_policy = "FORCE"
  column_info {
    column_index    = 0
    original_column = "username"
    column_type     = "STRING"
  }
  column_info {
    column_index    = 1
    original_column = "ip_address"
    column_type     = "CIDR"
  }
}

resource "google_chronicle_data_table_rows" "example_rows" {
  location      = "us"
  instance      = "%{instance_id}"
  data_table_id = google_chronicle_data_table.example_dt.data_table_id
  rows {
    values           = ["user1", "192.168.1.1/32"]
    row_time_to_live = "48h"
  }
  rows {
    values           = ["user2", "10.0.0.1/32"]
    row_time_to_live = "48h"
  }
  rows {
    values           = ["user3", "172.16.0.1/32"]
    row_time_to_live = "48h"
  }
}
`, context)
}
