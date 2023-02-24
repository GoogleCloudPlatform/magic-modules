package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceCloudIdentityGroupMemberships_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    acctest.GetTestOrgDomainFromEnv(t),
		"cust_id":       acctest.GetTestCustIdFromEnv(t),
		"identity_user": acctest.GetTestIdentityUserFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	memberId := Nprintf("%{identity_user}@%{org_domain}", context)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.TestAccPreCheck(t) },
		Providers: acctest.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembershipConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.#", "1"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.0.roles.#", "2"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.0.preferred_member_key.0.id", memberId),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupMembershipConfig(context map[string]interface{}) string {
	return testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExample(context) + Nprintf(`

data "google_cloud_identity_group_memberships" "members" {
  group = google_cloud_identity_group_membership.cloud_identity_group_membership_basic.group
}
`, context)
}
