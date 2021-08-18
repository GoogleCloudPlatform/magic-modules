package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFirewallPolicyAssociation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", getTestOrgFromEnv(t)),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyAssociation_start(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_association.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeFirewallPolicyAssociation_start(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
	parent      = google_folder.folder.name
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_association" "default" {
	firewall_policy = google_compute_firewall_policy.default.name
  attachment_target = google_folder.folder.name
  name = "tf-test-association-%{random_suffix}"
}
`, context)
}
