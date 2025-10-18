package cloudidentity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleCloudIdentityPolicy(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudIdentityPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "name", "policies/C0123abc:123456789"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "customer", "customers/C0123abc"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "policy_query.0.query", "entity.org_units.exists(org_unit, org_unit.org_unit_id == '123456789')"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "policy_query.0.group", ""),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "policy_query.0.org_unit", "123456789"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "policy_query.0.sort_order", "0"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "setting", `{"type":"some.setting.type","value":{"boolValue":true}}`),
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "type", "SYSTEM"),
				),
			},
		},
	})
}

const testAccDataSourceGoogleCloudIdentityPolicyConfig = `
data "google_cloud_identity_policy" "test" {
  name = "policies/C0123abc:123456789"
}
`
