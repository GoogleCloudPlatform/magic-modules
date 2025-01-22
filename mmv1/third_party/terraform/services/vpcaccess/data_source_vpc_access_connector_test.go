package vpcaccess_test

import (
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVPCAccessConnectorDatasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vpc-access-connector"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnectorDatasourceConfig_basic(context),
			},
			{
				ResourceName:            "google_vpc_access_connector.connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "self_link"},
			},
		},
	})
}

func TestAccVPCAccessConnectorDatasource_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vpc-access-connector"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnectorDatasourceConfig_basic(context),
			},
			{
				ResourceName:            "google_vpc_access_connector.connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "self_link"},
			},
			{
				Config: testAccVPCAccessConnectorDatasourceConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vpc_access_connector.connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_vpc_access_connector.connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "self_link"},
			},
		},
	})
}

func testAccVPCAccessConnectorDatasourceConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          	= "tf-test-%{random_suffix}"
  ip_cidr_range 	= "10.8.0.32/28"
  network 			= data.google_compute_network.default.id
  region        	= "us-central1"
  machine_type      = "e2-micro"
  min_instances     = 2
  max_instances     = 3
}

data "google_vpc_access_connector" "connector" {
  name = google_vpc_access_connector.connector.name
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccVPCAccessConnectorDatasourceConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          	= "tf-test-%{random_suffix}"
  ip_cidr_range 	= "10.8.0.32/28"
  network 			= data.google_compute_network.default.id
  region        	= "us-central1"
  machine_type      = "f1-micro"
  min_instances     = 3
  max_instances     = 5
}

data "google_vpc_access_connector" "connector" {
  name = google_vpc_access_connector.connector.name
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}
