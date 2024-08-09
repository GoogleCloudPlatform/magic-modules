package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2OrganizationBigQueryExportConfig_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	dataset_id := "tf_test_" + randomSuffix

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": randomSuffix,
		"dataset_id": dataset_id, 
		"dataset": fmt.Sprintf("projects/%s/datasets/%s",
					envvar.GetTestProjectFromEnv(), dataset_id),
		"export_id": "tf-test-export-"+randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationBigQueryExportConfig_basic(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_scc_big_query_exports.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterV2OrganizationBigQueryExportConfig_update(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_scc_big_query_exports.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterV2OrganizationBigQueryExportConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  labels = {
    env = "default"
  }
}

resource "google_scc_v2_organization_scc_big_query_exports" "default" {
  export_id    = "%{export_id}"
  organization = "%{org_id}"
  dataset      = "%{dataset}"
  location     = "global"
  description  = "Cloud Security Command Center Findings Big Query Export Config"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""
}
`, context)
}

func testAccSecurityCenterV2OrganizationBigQueryExportConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  labels = {
    env = "default"
  }
}

resource "google_scc_v2_organization_scc_big_query_exports" "default" {
  export_id    = "%{export_id}"
  organization = "%{org_id}"
  dataset      = "%{dataset}"
  location     = "global"
  description  = "SCC Findings Big Query Export Update"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""
}
`, context)
}
