package securitycenterv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2ProjectMuteConfig_basic(t *testing.T) {
	t.Parallel()

	contextBasic := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"folder_name":   acctest.RandString(t, 10),
		"project_id":    fmt.Sprintf("tf-test-project-%s", acctest.RandString(t, 10)),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
		"parent_org":    os.Getenv("TEST_ORG_ID"),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
	}

	contextHighSeverity := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"folder_name":   acctest.RandString(t, 10),
		"project_id":    fmt.Sprintf("tf-test-project-%s", acctest.RandString(t, 10)),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
		"parent_org":    os.Getenv("TEST_ORG_ID"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2ProjectMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test1", "description", "A test project mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test1", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test1", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test1", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test1", "project", fmt.Sprintf("projects/%s", contextBasic["project_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_project_mute_config.project_mute_test1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "location"},
			},
			{
				Config: testAccSecurityCenterV2ProjectMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test2", "description", "A test project mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test2", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test2", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test2", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.project_mute_test2", "project", fmt.Sprintf("projects/%s", contextHighSeverity["project_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_project_mute_config.project_mute_test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "location"},
			},
		},
	})
}

func testAccSecurityCenterV2ProjectMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name       = "Test project"
  project_id = "%{project_id}"
  org_id     = "%{parent_org}"
}

resource "google_project_iam_member" "project1_iam_member" {
	project = google_project.project.project_id
	role    = "roles/securitycenter.admin"
	member = "serviceAccount:%{service_account}"
}  

resource "google_scc_v2_project_mute_config" "project_mute_test1" {
  description          = "A test project mute config"
  filter               = "severity = \"LOW\""
  mute_config_id       = "tf-test-my-config-%{random_suffix}"
  location             = "%{location}"
  project              = "${google_project.project.project_id}"
  type                 =  "STATIC"
  depends_on   = [google_project_iam_member.project1_iam_member]
}
`, context)
}

func testAccSecurityCenterV2ProjectMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_project" "project" {
  name       = "Test project"
  project_id = "%{project_id}"
  org_id     = "%{parent_org}"
}

resource "google_project_iam_member" "project2_iam_member" {
	project = google_project.project.project_id
	role    = "roles/securitycenter.admin" 
	member = "serviceAccount:%{service_account}"
}

resource "google_scc_v2_project_mute_config" "project_mute_test2" {
  description          = "A test project mute config with high severity"
  filter               = "severity = \"HIGH\""
  mute_config_id       = "tf-test-my-config-%{random_suffix}"
  location             = "%{location}"
  project              = "${google_project.project.project_id}"
  type                 =  "STATIC"
}
`, context)
}
