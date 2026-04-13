package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccServiceAccountListResource_queryIdentity lists service accounts via the
// provider list resource API and asserts a known identity appears in the query
// results (Terraform 1.14+).
func TestAccServiceAccountListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountBasic(accountId, "Terraform List Test", "list resource query test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_service_account.acceptance", "email", expectedEmail),
					resource.TestCheckResourceAttr("google_service_account.acceptance", "project", project),
				),
			},
			{
				Query:  true,
				Config: testAccServiceAccountListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_service_account.all_in_project", map[string]knownvalue.Check{
						"email":   knownvalue.StringExact(expectedEmail),
						"project": knownvalue.StringExact(project),
					}),
					querycheck.ExpectLengthAtLeast("google_service_account.all_in_project", 1),
				},
			},
		},
	})
}

func testAccServiceAccountListQuery(project string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_service_account" "all_in_project" {
  provider = google

  config {
    project = %q
  }
}
`, project)
}
