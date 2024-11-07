package accesscontextmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceAccessContextManagerServicePerimeter_basicTest(t *testing.T) {

	org := envvar.GetTestOrgFromEnv(t)
	policyTitle := "my title"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeterDataSource_basic(org, policyTitle),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_access_context_manager_access_policy.policy", "google_access_context_manager_access_policy.policy"),
				),
			},
		},
	})
}

func testAccAccessContextManagerServicePerimeterDataSource_basic(org, policyTitle string) string {
	return acctest.Nprintf(`
resource "google_access_context_manager_access_policy" "policy" {
  parent = "organizations/%{org}"
  title  = "%{policyTitle}"
}

data "google_access_context_manager_access_policy" "policy" {
  parent = "organizations/%{org}"
  depends_on = [ google_access_context_manager_access_policy.policy ]
}
`, map[string]interface{}{"org": org, "policyTitle": policyTitle})
}
