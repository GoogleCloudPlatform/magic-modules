package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApikeysKey_basic(t *testing.T) {
	// DCL currently fails due to transport modification
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"project":       getTestProjectFromEnv(),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: funcAccTestApikeysKeyCheckDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_basic(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_apikeys_key.key",
			},
		},
	})
}

func testAccApikeysKey_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_apikeys_key" "key" {
	display_name = "key%{random_suffix}"
	project = "%{project}"
}
`, context)
}

func funcAccTestApikeysKeyCheckDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_trigger" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ApikeysBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("ApikeysKey still exists at %s", url)
			}
		}

		return nil
	}
}
