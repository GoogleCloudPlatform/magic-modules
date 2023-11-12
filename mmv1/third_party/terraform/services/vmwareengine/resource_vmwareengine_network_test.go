package vmwareengine_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVmwareengineNetwork_vmwareEngineNetworkUpdate(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
	}

	configTemplate := vmwareEngineNetworkConfigTemplate(context)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(configTemplate, "description1"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: fmt.Sprintf(configTemplate, "description2"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func vmwareEngineNetworkConfigTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  name        = "tf-test-network-%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "%s"
}
`, context)
}
