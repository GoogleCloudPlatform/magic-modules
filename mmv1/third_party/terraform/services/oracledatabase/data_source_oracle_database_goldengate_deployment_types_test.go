package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/oracledatabase"
)

func TestAccOracleDatabaseGoldengateDeploymentTypes_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseGoldengateDeploymentTypesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_types.my_deployment_types", "goldengate_deployment_types.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_deployment_types.my_deployment_types", "goldengate_deployment_types.0.deployment_type"),
				),
			},
		},
	})
}

func testAccOracleDatabaseGoldengateDeploymentTypesConfig() string {
	return fmt.Sprintf(`
data "google_oracle_database_goldengate_deployment_types" "my_deployment_types" {
	location = "us-east4"
	project  = "oci-terraform-testing-prod"
}
`)
}
