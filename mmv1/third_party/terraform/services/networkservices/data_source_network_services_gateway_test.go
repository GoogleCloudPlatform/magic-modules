package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/networkservices"
)

func TestAccNetworkServicesGatewayDatasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGatewayDatasource_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_network_services_gateway.default",
						"google_network_services_gateway.default",
						[]string{
							// Client-side only virtual field; not returned by the API GET.
							"delete_swg_autogen_router_on_destroy",
						},
					),
				),
			},
		},
	})
}

func testAccNetworkServicesGatewayDatasource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_gateway" "default" {
  name        = "tf-test-gateway-%{random_suffix}"
  scope       = "default-scope-basic"
  type        = "OPEN_MESH"
  ports       = [443]
  description = "my description"
}

data "google_network_services_gateway" "default" {
  name     = google_network_services_gateway.default.name
  location = google_network_services_gateway.default.location
}
`, context)
}
