package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/oracledatabase"
)

func TestAccDataSourceOracleDatabaseOdbDnsZone_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOracleDatabaseOdbDnsZone_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_dns_zone.my-zone", "name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_dns_zone.my-zone", "domain"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_dns_zone.my-zone", "create_time"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_odb_dns_zone.my-zone", "state"),
				),
			},
		},
	})
}

func testAccDataSourceOracleDatabaseOdbDnsZone_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_odb_dns_zone" "my-zone" {
  odb_dns_zone_id = "tf-test-permanent-odb-dns-zone"
  odbnetwork = "tf-test-permanent-odbnetwork"
  location = "europe-west2"
  project = "oasis-terraform-testing-prod"
}
`)
}
