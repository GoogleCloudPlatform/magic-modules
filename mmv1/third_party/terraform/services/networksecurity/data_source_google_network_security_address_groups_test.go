package networksecurity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceNetworkSecurityAddressGroups_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      "us-central1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkSecurityAddressGroupsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_network_security_address_groups.all", "address_groups.#"),
					resource.TestCheckResourceAttr("data.google_network_security_address_groups.all", "address_groups.0.location", context["location"].(string)),
				),
			},
		},
	})
}

func testAccDataSourceNetworkSecurityAddressGroupsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_security_address_group" "basic" {
  name        = "tf-test-ag-%{random_suffix}"
  parent      = "projects/%{project}"
  location    = "%{location}"
  type        = "IPV4"
  capacity    = 100
  items       = ["208.80.154.224/32"]
}

data "google_network_security_address_groups" "all" {
  project    = "%{project}"
  location   = "%{location}"
  depends_on = [google_network_security_address_group.basic]
}
`, context)
}

func TestAccDataSourceNetworkSecurityAddressGroups_organization(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"location":      "us-central1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkSecurityAddressGroupsOrganizationConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_network_security_address_groups.org_all", "address_groups.#"),
					resource.TestCheckResourceAttr("data.google_network_security_address_groups.org_all", "address_groups.0.location", context["location"].(string)),
				),
			},
		},
	})
}

func testAccDataSourceNetworkSecurityAddressGroupsOrganizationConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_security_address_group" "org_basic" {
  name        = "tf-test-ag-org-%{random_suffix}"
  parent      = "organizations/%{org_id}"
  location    = "%{location}"
  type        = "IPV4"
  capacity    = 100
  items       = ["208.80.154.224/32"]
}

data "google_network_security_address_groups" "org_all" {
  parent     = "organizations/%{org_id}"
  location   = "%{location}"
  depends_on = [google_network_security_address_group.org_basic]
}
`, context)
}
