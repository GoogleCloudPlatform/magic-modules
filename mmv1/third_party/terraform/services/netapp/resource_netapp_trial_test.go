package netapp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetappTrial_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappTrialDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappTrial_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_netapp_trial.default", "name"),
					resource.TestCheckResourceAttr("google_netapp_trial.default", "location", "us-central1"),
				),
			},
			{
				ResourceName:      "google_netapp_trial.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetappTrial_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_trial" "default" {
  location = "us-central1"
}
`, context)
}

func testAccCheckNetappTrialDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_netapp_trial" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "projects/{{project}}/locations/{{location}}/trial")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("NetappTrial still exists at %s", url)
			}
		}

		return nil
	}
}
