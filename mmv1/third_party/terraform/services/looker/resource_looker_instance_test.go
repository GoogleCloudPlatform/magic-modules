package looker_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLookerInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLookerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLookerInstance_lookerInstanceBasicExample(context),
			},
			{
				ResourceName:            "google_looker_instance.looker-instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_config", "region"},
			},
			{
				Config: testAccLookerInstance_lookerInstanceFullExample(context),
			},
			{
				ResourceName:            "google_looker_instance.looker-instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_config", "region"},
			},
		},
	})
}

func TestAccLookerInstance_updateControlledEgress(t *testing.T) {
	t.Parallel()

	// Step 1: Create instance WITHOUT controlled egress
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Config A: Basic Instance
				Config: testAccLookerInstance_basic(context),
			},
			{
				// Config B: Update to ENABLE controlled egress
				// This triggers the PATCH with your update_mask logic
				Config: testAccLookerInstance_controlledEgress(context),
			},
		},
	})
}

func testAccLookerInstance_basic(context map[string]interface{}) string {
	return fmt.Sprintf(`
resource "google_looker_instance" "test" {
  name               = "tf-test-looker-%s"
  platform_edition   = "LOOKER_CORE_ENTERPRISE_ANNUAL"
  region             = "us-central1"
  public_ip_enabled  = true

  oauth_config {
    client_id     = "my-client-id"
    client_secret = "my-client-secret"
  }
}
`, context["random_suffix"])
}

func testAccLookerInstance_controlledEgress(context map[string]interface{}) string {
	return fmt.Sprintf(`
resource "google_looker_instance" "test" {
  name               = "tf-test-looker-%s"
  platform_edition   = "LOOKER_CORE_ENTERPRISE_ANNUAL"
  region             = "us-central1"
  public_ip_enabled  = true

  controlled_egress_enabled = true

  controlled_egress_config {
    marketplace_enabled = true
    egress_fqdns        = ["google.com", "github.com"]
  }

  oauth_config {
    client_id     = "my-client-id"
    client_secret = "my-client-secret"
  }
}
`, context["random_suffix"])
}
