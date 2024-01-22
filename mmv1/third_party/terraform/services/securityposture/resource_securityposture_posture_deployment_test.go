package securityposture_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":         envvar.GetTestOrgFromEnv(t),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_basic(context),
			},
			{
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment_one",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context),
			},
			{
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment_one",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
		},
	})
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_one" {
    posture_id          = "posture_one"
    parent = "organizations/%{org_id}/locations/global"
    state = "ACTIVE"
    description = "a new posture"
    policy_sets {
        policy_set_id = "org_policy_set"
        description = "set of org policies"
        policies {
            policy_id = "policy_1"
            constraint {
                org_policy_constraint {
                    canned_constraint_id = "storage.uniformBucketLevelAccess"
                    policy_rules {
                        enforce = true
                    }
                }
            }
        }
        policies {
    		policy_id = "policy_2"
    		constraint {
    			org_policy_constraint_custom {
    				custom_constraint {
    					name         = "organizations/%{org_id}/customConstraints/custom.disableGkeAutoUpgrade"
					  	display_name = "Disable GKE auto upgrade"
					  	description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

					  	action_type    = "ALLOW"
					  	condition      = "resource.management.autoUpgrade == false"
					  	method_types   = ["CREATE", "UPDATE"]
					  	resource_types = ["container.googleapis.com/NodePool"]
    				}
    				policy_rules {
    					enforce = true
    				}
    			}
    		}
		}
    }
}

resource "google_securityposture_posture_deployment" "postureDeployment_one" {
	posture_deployment_id          = "posture_deployment_one"
	parent = "organizations/%{org_id}"
	location = "global"
    description = "a new posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = google_securityposture_posture.posture_one.name
    posture_revision_id = google_securityposture_posture.posture_one.revision_id
}
`, context)
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_one" {
    posture_id          = "posture_one"
    parent = "organizations/%{org_id}/locations/global"
    state = "ACTIVE"
    description = "a new posture"
    policy_sets {
        policy_set_id = "org_policy_set"
        description = "set of org policies"
        policies {
            policy_id = "policy_1"
            constraint {
                org_policy_constraint {
                    canned_constraint_id = "storage.publicAccessPrevention"
                    policy_rules {
                        enforce = true
                    }
                }
            }
        }
    }
}

resource "google_securityposture_posture_deployment" "postureDeployment_one" {
	posture_deployment_id          = "posture_deployment_one"
	parent = "organizations/%{org_id}"
	location = "global"
    description = "an updated posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = google_securityposture_posture.posture_one.name
    posture_revision_id = google_securityposture_posture.posture_one.revision_id
}
`, context)
}
