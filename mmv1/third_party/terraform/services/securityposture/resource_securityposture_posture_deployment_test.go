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
		CheckDestroy:             testAccCheckSecurityPosturePostureDeploymentDestroyProducer(t),
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
	name          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "a new posture deployment"
    target_resource = "folders/123456"
    posture_id = "organizations/%{org_id}/locations/global/postures/posture1"
    posture_revision_id = "abcdef"
}
`, context)
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture_deployment" "postureDeployment" {
	name          = "posture_deployment_1"
	parent = "organizations/%{org_id}/locations/global"
    description = "an updated posture deployment"
    target_resource = "folders/123456"
    posture_id = "organizations/%{org_id}/locations/global/postures/posture2"
    posture_revision_id = "bcdefg"
}
`, context)
}
