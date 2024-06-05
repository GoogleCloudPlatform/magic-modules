package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccSecurityCenterV2FolderMuteConfig_basic(t *testing.T) {
	t.Parallel()

	contextBasic := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"folder_id":     acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	contextHighSeverity := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"folder_id":     acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "us_central",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterV2FolderMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterv2FolderMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "description", "A test folder mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "folder_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "parent", fmt.Sprintf("folders/%s", contextBasic["folder_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_folder_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "folder_mute_config_id"},
			},
			{
				Config: testAccSecurityCenterV2FolderMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "description", "A test folder mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "folder_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_folder_mute_config.default", "parent", fmt.Sprintf("folders/%s", contextHighSeverity["folder_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_folder_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "folder_mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterV2FolderMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_folder_mute_config" "default" {
  description          = "A test folder mute config"
  filter               = "severity = \"LOW\""
  folder_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "folders/%{folder_id}"
}
`, context)
}

func testAccSecurityCenterV2FolderMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_folder_mute_config" "default" {
  description          = "A test folder mute config with high severity"
  filter               = "severity = \"HIGH\""
  folder_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "folders/%{folder_id}"
}
`, context)
}
func testAccCheckSecurityCenterV2FolderMuteConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.Provider.Meta().(*transport_tpg.Config)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_v2_folder_mute_config" {
				continue
			}
			// Initialize Security Command Center Service
			sc, err := securitycenter.NewService(context.Background(), config.GoogleClientOptions...)
			if err != nil {
				return fmt.Errorf("Error creating Security Command Center client: %s", err)
			}
			// Get the folder mute config by name
			name := rs.Primary.ID
			_, err = sc.Folders.MuteConfigs.Get(name).Do()
			if err == nil {
				return fmt.Errorf("Folder mute config %s still exists", name)
			}
		}

		return nil
	}
}
