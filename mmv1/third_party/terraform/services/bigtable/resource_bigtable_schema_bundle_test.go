package bigtable_test

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBigtableSchemaBundle_update(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	sbName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableSchemaBundle_update(instanceName, tableName, sbName, "proto_schema_bundle"),
			},
			{
				ResourceName:            "google_bigtable_schema_bundle.schema_bundle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableSchemaBundle_update(instanceName, tableName, sbName, "updated_proto_schema_bundle"),

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigtable_schema_bundle.schema_bundle", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_bigtable_schema_bundle.schema_bundle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableSchemaBundle_update(instanceName, tableName, sbName, "proto_schema_bundle"),
			},
			{
				ResourceName:            "google_bigtable_schema_bundle.schema_bundle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func testAccBigtableSchemaBundle_update(instanceName, tableName, sbName, fileName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s-c"
    zone       = "us-east1-b"
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

resource "google_bigtable_schema_bundle" "schema_bundle" {
  schema_bundle_id = "%s"
  instance         = google_bigtable_instance.instance.name
  table            = google_bigtable_table.table.name

  proto_schema {
    proto_descriptors = filebase64("test-fixtures/%s.pb")
  }

  ignore_warnings = true
}
`, instanceName, instanceName, tableName, sbName, fileName)
}
