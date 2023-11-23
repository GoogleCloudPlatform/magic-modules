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
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context),
			},
			{
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
		},
	})
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture_deployment" "postureDeployment" {
	posture_deployment_id          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "a new posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = "organizations/%{org_id}/locations/global/postures/testPosture"
    posture_revision_id = "eb29beb8"
}
`, context)
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture_deployment" "postureDeployment" {
	posture_deployment_id          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "an updated posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = "organizations/%{org_id}/locations/global/postures/posture-foo-5"
    posture_revision_id = "48e17293"
}
`, context)
}
