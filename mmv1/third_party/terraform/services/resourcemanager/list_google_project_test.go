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
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

// TestAccProjectListResource_queryIdentity lists projects via the
// provider list resource API and asserts a known identity appears in the query
// results (Terraform 1.14+).
func TestAccProjectListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: testAccProjectListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_project.all_in_project", map[string]knownvalue.Check{
						"project_id": knownvalue.StringExact(project),
					}),
					querycheck.ExpectLengthAtLeast("google_project.all_in_project", 1),
				},
			},
		},
	})
}

func testAccProjectListQuery(project string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_project" "all_in_project" {
  provider = google

  config {
    filter = "id:%s"
  }
}
`, project)
}
