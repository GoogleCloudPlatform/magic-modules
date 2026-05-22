package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseGoldengateConnectionTypes_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseGoldengateConnectionTypesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_connection_types.my_connection_types", "goldengate_connection_types.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_connection_types.my_connection_types", "goldengate_connection_types.0.connection_type"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_goldengate_connection_types.my_connection_types", "goldengate_connection_types.0.technology_types.#"),
				),
			},
		},
	})
}

func testAccOracleDatabaseGoldengateConnectionTypesConfig() string {
	return fmt.Sprintf(`
data "google_oracle_database_goldengate_connection_types" "my_connection_types" {
	location = "us-east4"
}
`)
}
