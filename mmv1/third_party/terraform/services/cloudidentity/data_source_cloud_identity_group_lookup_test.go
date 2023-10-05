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
					resource.TestCheckResourceAttrPair("data.google_cloud_identity_group_lookup.lookup", "name",
						"google_cloud_identity_group.cloud_identity_group_basic", "name"),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupLookupConfig(context map[string]interface{}) string {
	// reused function below creates a group resource `google_cloud_identity_group.cloud_identity_group_basic`
	return testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample(context) + acctest.Nprintf(`
data "google_cloud_identity_group_lookup" "lookup" {
  group_key {
    id = google_cloud_identity_group.cloud_identity_group_basic.group_key[0].id
  }
}
`, context)
}
