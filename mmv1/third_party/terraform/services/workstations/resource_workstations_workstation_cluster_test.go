package workstations_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkstationsWorkstationCluster_update(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":                  envvar.GetTestProjectFromEnv(),
		"location":                 "us-central1",
		"random_suffix":            randString,
		"workstation_cluster_name": "tf-test-workstation-cluster" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationCluster_workstationClusterBasicExample(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationCluster_update(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccWorkstationsWorkstationCluster_Private_update(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":                  envvar.GetTestProjectFromEnv(),
		"location":                 "us-central1",
		"random_suffix":            randString,
		"workstation_cluster_name": "tf-test-workstation-cluster-private" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationCluster_workstationClusterPrivateExample(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationCluster_private_update(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkstationsWorkstationCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_workstations_workstation_cluster" "default" {
  provider   		     = google-beta
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location   		     = "us-central1"

  labels = {
    foo = "bar"
  }
}

resource "google_compute_network" "default" {
  provider                = google-beta
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  provider      = google-beta
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.default.name
}
`, context)
}

func testAccWorkstationsWorkstationCluster_private_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_workstations_workstation_cluster" "default" {
  provider   		     = google-beta
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location   		     = "us-central1"

  private_cluster_config {
    allowed_projects        = ["${data.google_project.project.project_id}"]
    enable_private_endpoint = true
  }

  labels = {
	foo = "bar"
  }
}

resource "google_compute_network" "default" {
  provider                = google-beta
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  provider      = google-beta
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.default.name
}
`, context)
}
