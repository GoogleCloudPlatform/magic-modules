package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRouterRoutePolicy_routerRoutePolicyExportExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterRoutePolicy_routerRoutePolicyExportExampleBasic(context),
			},
			{
				ResourceName:            "google_compute_router_route_policy.rp-export",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "router"},
			},
			{
				Config: testAccComputeRouterRoutePolicy_routerRoutePolicyExportExampleUpdate(context),
			},
			{
				ResourceName:            "google_compute_router_route_policy.rp-export",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "router"},
			},
		},
	})
}

func testAccComputeRouterRoutePolicy_routerRoutePolicyExportExampleBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  network       = google_compute_network.net.id
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "tf-test-my-router%{random_suffix}"
  region  = google_compute_subnetwork.subnet.region
  network = google_compute_network.net.id
}

resource "google_compute_router_route_policy" "rp-export" {
  router = google_compute_router.router.name
  region = google_compute_router.router.region
  name = "tf-test-my-rp1%{random_suffix}"
  type = "ROUTE_POLICY_TYPE_EXPORT"
  terms {
    priority = 1
    match {
      expression = "destination == '10.0.0.0/12'"
	}
    actions {
      expression = "accept()"
    }
  }
}
`, context)
}

func testAccComputeRouterRoutePolicy_routerRoutePolicyExportExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  network       = google_compute_network.net.id
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "tf-test-my-router%{random_suffix}"
  region  = google_compute_subnetwork.subnet.region
  network = google_compute_network.net.id
}

resource "google_compute_router_route_policy" "rp-export" {
  router = google_compute_router.router.name
  region = google_compute_router.router.region
  name = "tf-test-my-rp1%{random_suffix}"
  type = "ROUTE_POLICY_TYPE_EXPORT"
  terms {
    priority = 2
    match {
      expression = "destination == '10.0.0.1/12'"
	}
    actions {
      expression = "accept()"
    }
  }
}
`, context)
}

func TestAccComputeRouterRoutePolicy_routerRoutePolicyImportExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterRoutePolicy_routerRoutePolicyImportExampleBasic(context),
			},
			{
				ResourceName:            "google_compute_router_route_policy.rp-import",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "router"},
			},
			{
				Config: testAccComputeRouterRoutePolicy_routerRoutePolicyImportExampleUpdate(context),
			},
			{
				ResourceName:            "google_compute_router_route_policy.rp-import",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "router"},
			},
		},
	})
}

func testAccComputeRouterRoutePolicy_routerRoutePolicyImportExampleBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  network       = google_compute_network.net.id
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "tf-test-my-router%{random_suffix}"
  region  = google_compute_subnetwork.subnet.region
  network = google_compute_network.net.id
}

resource "google_compute_router_route_policy" "rp-import" {
  name = "tf-test-my-rp2%{random_suffix}"
  router = google_compute_router.router.name
  region = google_compute_router.router.region
  type = "ROUTE_POLICY_TYPE_IMPORT"
  terms {
    priority = 2
    match {
      expression = "destination == '10.0.0.0/12'"
	}
    actions {
      expression = "accept()"
    }
  }
}
`, context)
}

func testAccComputeRouterRoutePolicy_routerRoutePolicyImportExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  network       = google_compute_network.net.id
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "tf-test-my-router%{random_suffix}"
  region  = google_compute_subnetwork.subnet.region
  network = google_compute_network.net.id
}

resource "google_compute_router_route_policy" "rp-import" {
  name = "tf-test-my-rp2%{random_suffix}"
  router = google_compute_router.router.name
  region = google_compute_router.router.region
  type = "ROUTE_POLICY_TYPE_IMPORT"
  terms {
    priority = 3
    match {
      expression = "destination == '10.0.0.1/12'"
	}
    actions {
      expression = "accept()"
    }
  }
}
`, context)
}
