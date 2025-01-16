package vmwareengine_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineNetworkPolicy_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":               "me-west1", // region with allocated quota
		"random_suffix":        acctest.RandString(t, 10),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"billing_account":      envvar.GetTestBillingAccountFromEnv(t),
		"vmwareengine_project": os.Getenv("GOOGLE_VMWAREENGINE_PROJECT"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description1", "192.168.0.0/26", false, false),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_network_policy.ds", "google_vmwareengine_network_policy.vmw-engine-network-policy", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description2", "192.168.1.0/26", true, true),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
		},
	})
}

func testAccVmwareengineNetworkPolicy_config(context map[string]interface{}, description string, edgeServicesCidr string, internetAccess bool, externalIp bool) string {
	context["internet_access"] = internetAccess
	context["external_ip"] = externalIp
	context["edge_services_cidr"] = edgeServicesCidr
	context["description"] = description

	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-policy-nw" {
  project           = "%{vmwareengine_project}"
  name              = "tf-test-sample-nw%{random_suffix}"
  location          = "global" 
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
  project = "%{vmwareengine_project}"
  location = "%{region}"
  name = "tf-test-sample-network-policy%{random_suffix}"
  description = "%{description}" 
  internet_access {
    enabled = "%{internet_access}"
  }
  external_ip {
    enabled = "%{external_ip}"
  }
  edge_services_cidr = "%{edge_services_cidr}"
  vmware_engine_network = google_vmwareengine_network.network-policy-nw.id
}

data "google_vmwareengine_network_policy" "ds" {
  project = "%{vmwareengine_project}"
  name = google_vmwareengine_network_policy.vmw-engine-network-policy.name
  location = "%{region}"
}
`, context)
}
