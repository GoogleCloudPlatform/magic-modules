package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/config"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/securitycenter/v1"
)

func TestAccSecurityCenterv2FolderMuteConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"folder_id":     acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterv2FolderMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterv2FolderMuteConfig_basic(context),
			},
			{
				ResourceName:            "google_scc_v2_folder_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterv2FolderMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_folder_mute_config" "default" {
  mute_config_id = "tf-test-my-config%{random_suffix}"
  parent         = "folders/%{folder_id}"
  description    = "A test folder mute config"
  filter         = "severity = \"LOW\""
}
`, context)
}

func testAccCheckSecurityCenterv2FolderMuteConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.Provider.Meta().(*config.Config)

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
			if !isGoogleAPINotFoundError(err) {
				return fmt.Errorf("Error while retrieving mute config: %s", err)
			}
		}

		return nil
	}
}

// isGoogleAPINotFoundError checks if an error is a Google API Not Found error
func isGoogleAPINotFoundError(err error) bool {
	apiErr, ok := err.(*googleapi.Error)
	return ok && apiErr.Code == 404
}
