package chromepolicy_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChromePolicies_basic(t *testing.T) {
	acctest.SkipIfVcr(t) // Chrome Policy API uses non-standard POST-based endpoints
	t.Parallel()

	custId := envvar.GetTestCustIdFromEnv(t)
	orgUnitId := os.Getenv("GOOGLE_CHROME_POLICY_ORG_UNIT_ID")
	if orgUnitId == "" {
		t.Skip("GOOGLE_CHROME_POLICY_ORG_UNIT_ID must be set for Chrome Policy acceptance tests")
	}

	context := map[string]interface{}{
		"cust_id":     custId,
		"org_unit_id": orgUnitId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChromePolicies_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_chrome_policies.test", "id"),
					resource.TestCheckResourceAttr("google_chrome_policies.test", "schema_filter", "chrome.users.MaxConnectionsPerProxy"),
				),
			},
		},
	})
}

func testAccChromePolicies_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  scopes = [
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
    "https://www.googleapis.com/auth/chrome.management.policy",
  ]
}

resource "google_chrome_policies" "test" {
  customer_id          = "%{cust_id}"
  org_unit_id          = "%{org_unit_id}"
  schema_filter = "chrome.users.MaxConnectionsPerProxy"

  policies = [
    {
      schema = "chrome.users.MaxConnectionsPerProxy"
      value = {
        maxConnectionsPerProxy = 32
      }
    },
  ]
}
`, context)
}
