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
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeploymentBasic(context),
			},
			{
				ResourceName:            "google_security_posture_posture_deployment.postureDeployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context),
			},
			{
				ResourceName:            "google_security_posture_posture_deployment.postureDeployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
		},
	})
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeploymentBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture_deployment" "postureDeployment" {
	posture_deployment_id          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "a new posture deployment"
    target_resource = "projects/190507214861"
    posture_id = "organizations/%{org_id}/locations/global/postures/gcloud-test-posture"
    posture_revision_id = "1fe5ff7a"
}
`, context)
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture_deployment" "postureDeployment" {
	posture_deployment_id          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "an updated posture deployment"
    target_resource = "projects/190507214861"
    posture_id = "organizations/%{org_id}/locations/global/postures/gcloud-test-list-posture"
    posture_revision_id = "101b17f1"
}
`, context)
}
