package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrgPolicyCustomConstraint_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrgPolicyCustomConstraintDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyCustomConstraint_v1(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"display_name", "Disable GKE auto upgrade"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"description", "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"action_type", "ALLOW"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"condition", "resource.management.autoUpgrade == false"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"method_types.0", "CREATE"),
				),
			},
			{
				ResourceName:            "google_org_policy_custom_constraint.constraint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccOrgPolicyCustomConstraint_v2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"display_name", "Updated"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"description", "Updated"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"action_type", "DENY"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"condition", "resource.management.autoUpgrade == true"),
					resource.TestCheckResourceAttr("google_org_policy_custom_constraint.constraint",
						"method_types.0", "UPDATE"),
				),
			},
			{
				ResourceName:            "google_org_policy_custom_constraint.constraint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccOrgPolicyCustomConstraint_v1(context map[string]interface{}) string {
	return Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "organizations/%{org_id}/customConstraints/custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}
`, context)
}

func testAccOrgPolicyCustomConstraint_v2(context map[string]interface{}) string {
	return Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "organizations/%{org_id}/customConstraints/custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Updated"
  description  = "Updated"

  action_type    = "DENY"
  condition      = "resource.management.autoUpgrade == true"
  method_types   = ["UPDATE", "CREATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}
`, context)
}
