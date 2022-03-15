package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testOrgResourceSettingName = "iam-serviceAccountKeyExpiry"
)

func TestAccOrganizationResourceSetting_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrganizationResourceSettingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "google_organization_resource_setting" "mysetting" {
  organization_id = "%s"
  setting_name = "%s"
  local_value {
     string_value = "1hours"
  }
}
`, getTestOrgFromEnv(t), testOrgResourceSettingName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_organization_resource_setting.mysetting", "setting_name", testOrgResourceSettingName),
					resource.TestCheckResourceAttr(
						"google_organization_resource_setting.mysetting", "local_value.0.string_value", "1hours"),
				),
			},
		},
	})
}

func testAccCheckOrganizationResourceSettingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_organization_resource_setting" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			id := "organizations/" + getTestOrgFromEnv(t) + "/settings/" + testOrgResourceSettingName
			setting, err := config.NewResourceSettingsClient(config.userAgent).Organizations.Settings.Get(id).Do()
			if err != nil {
				return fmt.Errorf("Unable to get ResourceSetting (it should still exist): %s", id)
			}

			if setting.LocalValue != nil {
				return fmt.Errorf("Expected localValue to be cleared, but it still exists: %+v", setting.LocalValue)
			}
		}

		return nil
	}
}

func Test_resourceSettingFullName(t *testing.T) {
	cases := []struct {
		parentType  string
		parentID    string
		settingName string

		expectedOutput string
	}{
		{
			parentType:  "folder",
			parentID:    "abc123",
			settingName: "my-coolSetting",

			expectedOutput: "folders/abc123/settings/my-coolSetting",
		},
	}

	for _, c := range cases {
		t.Run(c.expectedOutput, func(t *testing.T) {
			out := resourceSettingFullName(c.parentType, c.parentID, c.settingName)
			if out != c.expectedOutput {
				t.Fatalf("expected: %s, got: %s", c.expectedOutput, out)
			}
		})
	}
}

func Test_resourceSettingShortName(t *testing.T) {
	cases := []struct {
		fullName       string
		expectedOutput string
	}{
		{
			fullName: "folders/abc123/settings/my-coolSetting",

			expectedOutput: "my-coolSetting",
		},
	}

	for _, c := range cases {
		t.Run(c.expectedOutput, func(t *testing.T) {
			out := resourceSettingShortName(c.fullName)
			if out != c.expectedOutput {
				t.Fatalf("expected: %s, got: %s", c.expectedOutput, out)
			}
		})
	}
}
