package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudIdentityGroup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersOiCS,
		CheckDestroy: testAccCheckCloudIdentityGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroup_basic(context),
			},
			{
				Config: testAccCloudIdentityGroup_update(context),
			},
		},
	})
}

func testAccCloudIdentityGroup_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_identity_group" "cloud_identity_group_basic" {
  provider = google-beta
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}
`, context)
}

func testAccCloudIdentityGroup_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_identity_group" "cloud_identity_group_basic" {
  provider = google-beta
  display_name = "tf-test-my-identity-group%{random_suffix}-update"
  description  = "my-description"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}
`, context)
}
