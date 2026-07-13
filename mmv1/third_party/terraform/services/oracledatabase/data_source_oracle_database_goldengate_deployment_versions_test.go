package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseGoldengateDeploymentVersions_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseGoldengateDeploymentVersionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_versions.my_deployment_versions", "goldengate_deployment_versions.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_versions.my_deployment_versions", "goldengate_deployment_versions.0.properties.0.ogg_version"),
				),
			},
		},
	})
}

func testAccOracleDatabaseGoldengateDeploymentVersionsConfig() string {
	return fmt.Sprintf(`
data "google_oracle_database_goldengate_deployment_versions" "my_deployment_versions" {
	location = "us-east4"
	project  = "oci-terraform-testing-prod"
}
`)
}
