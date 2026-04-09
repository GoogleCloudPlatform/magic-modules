package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseOdbSubnet_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseOdbSubnet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_subnet.my-subnet", "name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_subnet.my-subnet", "cidr_range"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_subnet.my-subnet", "create_time"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_subnet.my-subnet", "purpose"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_subnet.my-subnet", "state"),
				),
			},
		},
	})
}

func testAccOracleDatabaseOdbSubnet_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_odb_subnet" "my-subnet" {
  odb_subnet_id = "tf-test-permanent-client-odbsubnet"
  odbnetwork = "tf-test-permanent-odbnetwork"
  location = "europe-west2"
  project = "oci-terraform-testing-prod"
}
`)
}
