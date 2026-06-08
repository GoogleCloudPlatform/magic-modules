package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/oracledatabase"
)

func TestAccOracleDatabaseGoldengateDeploymentEnvironments_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseGoldengateDeploymentEnvironmentsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_environments.my_deployment_environments", "goldengate_deployment_environments.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_environments.my_deployment_environments", "goldengate_deployment_environments.0.environment_type"),
				),
			},
		},
	})
}

func testAccOracleDatabaseGoldengateDeploymentEnvironmentsConfig() string {
	return fmt.Sprintf(`
data "google_oracle_database_goldengate_deployment_environments" "my_deployment_environments" {
	location = "us-east4"
	project  = "oci-terraform-testing-prod"
}
`)
}
