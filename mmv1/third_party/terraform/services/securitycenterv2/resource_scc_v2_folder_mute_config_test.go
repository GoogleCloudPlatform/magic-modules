package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2FolderMuteConfig_basic(t *testing.T) {
	t.Parallel()

	contextBasic := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"folder_name":     fmt.Sprintf("test1-%s", acctest.RandString(t, 10)),
		"random_suffix":   acctest.RandString(t, 10),
		"location":        "global",
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
	}

	contextHighSeverity := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"folder_name":     fmt.Sprintf("test2-%s", acctest.RandString(t, 10)),
		"random_suffix":   acctest.RandString(t, 10),
		"location":        "global",
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2FolderMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test1", "description", "A test folder mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test1", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test1", "folder_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test1", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test1", "parent", fmt.Sprintf("folders/%s", contextBasic["folder_name"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_folder_mute_config.folder_mute_test1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "folder_mute_config_id"},
			},
			{
				Config: testAccSecurityCenterV2FolderMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test2", "description", "A test folder mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test2", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test2", "folder_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test2", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.folder_mute_test2", "parent", fmt.Sprintf("folders/%s", contextHighSeverity["folder_name"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_folder_mute_config.folder_mute_test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "folder_mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterV2FolderMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "test1" {
  display_name = "%{folder_name}"
  parent       = "organizations/%{org_id}"
}

resource "google_project_iam_binding" "folder_test1_binding" {
	project = "billingProject"
	role    = "roles/securitycenter.admin"
	members = [
	  "serviceAccount: %{service_account}",
	]
}

resource "google_scc_v2_folder_mute_config" "folder_mute_test1" {
  description          = "A test folder mute config"
  filter               = "severity = \"LOW\""
  mute_config_id       = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "${google_folder.test1.id}"
  type                 =  "STATIC"
}
`, context)
}

func testAccSecurityCenterV2FolderMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "test2" {
  display_name = "%{folder_name}"
  parent       = "organizations/%{org_id}"
}

resource "google_project_iam_binding" "folder_test2_binding" {
	project = "YOUR_PROJECT_ID"
	role    = "roles/securitycenter.admin"
	members = [
	  "serviceAccount:YOUR_SERVICE_ACCOUNT",
	]
}

resource "google_scc_v2_folder_mute_config" "folder_mute_test2" {
  description          = "A test folder mute config with high severity"
  filter               = "severity = \"HIGH\""
  mute_config_id       = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "${google_folder.test2.id}"
  type                 =  "STATIC"
}
`, context)
}
