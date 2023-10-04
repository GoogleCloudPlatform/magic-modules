package cloudidentity_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccDataSourceCloudIdentityGroupLookup_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupLookupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_cloud_identity_group_lookup.lookup",
						"name", regexp.MustCompile("^groups/.*$")),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupLookupConfig(context map[string]interface{}) string {
	return testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample(context) + acctest.Nprintf(`
data "google_cloud_identity_group_lookup" "lookup" {
  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }
}
`, context)
}
