package bigtable_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBigtableMaterializedView_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	mvName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableMaterializedView(instanceName, tableName, mvName),
			},
			{
				ResourceName:      "google_bigtable_materialized_view.materialized_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableMaterializedView(instanceName, tableName, mvName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s-c"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id

  column_family {
	family = "CF"
  }
}

resource "google_bigtable_materialized_view" "materialized_view" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  deletion_protection = false
  query = <<EOT
SELECT _key, CF 
FROM %s
EOT  
}

  depends_on = [
    google_bigtable_table.table
  ]
`, instanceName, instanceName, tableName, mvName, tableName)
}
