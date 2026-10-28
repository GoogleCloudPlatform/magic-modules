package workstations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/workstations"
)

func TestAccWorkstationsWorkstation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstation_basic(context),
			},
			{
				ResourceName:            "google_workstations_workstation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "terraform_labels", "workstation_cluster_id", "workstation_config_id", "workstation_id"},
			},
			{
				Config: testAccWorkstationsWorkstation_modified(context),
			},
			{
				ResourceName:            "google_workstations_workstation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "terraform_labels", "workstation_cluster_id", "workstation_config_id", "workstation_id"},
			},
		},
	})
}

func testAccWorkstationsWorkstation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  prefer_regional_endpoints = true
}

resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-east4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-east4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation" "default" {
  workstation_id         = "tf-test-workstation%{random_suffix}"
  workstation_config_id  = google_workstations_workstation_config.default.workstation_config_id
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east4"

  labels = {
    foo = "bar"
  }

  annotations = {
    label-one = "value-one"
  }

  env = {
    name = "bar"
  }
}
`, context)
}

func testAccWorkstationsWorkstation_modified(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  prefer_regional_endpoints = true
}

resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-east4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-east4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation" "default" {
  workstation_id         = "tf-test-workstation%{random_suffix}"
  workstation_config_id  = google_workstations_workstation_config.default.workstation_config_id
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east4"
  display_name           = "workstation%{random_suffix}"

  labels = {
    foo = "baz"
  }

  annotations = {
    label-one = "value-two"
  }

  env = {
    name = "test"
  }
}
`, context)
}
