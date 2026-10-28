package workstations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/tags"
	_ "github.com/hashicorp/terraform-provider-google/google/services/workstations"
)

func TestAccWorkstationsWorkstationCluster_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"workstation_cluster_name": "tf-test-workstation-cluster" + randomSuffix,
		"key_short_name":           "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name":         "tf-test-value-" + acctest.RandString(t, 10),
		"random_suffix":            randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationCluster_full(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "tags", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationCluster_update(context),
			},
			{
				ResourceName:            "google_workstations_workstation_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccWorkstationsWorkstationCluster_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "tag_key" {
  parent     = "projects/${data.google_project.project.number}"
  short_name = "%{key_short_name}"
}

resource "google_tags_tag_value" "tag_value" {
  parent     = "tagKeys/${google_tags_tag_key.tag_key.name}"
  short_name = "%{value_short_name}"
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "europe-west4"
  display_name           = "cluster-display-name"

  private_cluster_config {
    enable_private_endpoint = true
  }

  domain_config {
    domain = "workstations.example.com"
  }

  workstation_authorization_url = "https://console.cloud.google.com/workstations/auth"
  workstation_launch_url        = "https://console.cloud.google.com/workstations/launch"

  labels = {
    "label" = "key"
  }

  annotations = {
    label-one = "value-one"
  }

  tags = {
    "${data.google_project.project.project_id}/${google_tags_tag_key.tag_key.short_name}" = "${google_tags_tag_value.tag_value.short_name}"
  }
}

resource "google_compute_network" "default" {
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "europe-west4"
  network       = google_compute_network.default.name
}
`, context)
}

func testAccWorkstationsWorkstationCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "tag_key" {
  parent     = "projects/${data.google_project.project.number}"
  short_name = "%{key_short_name}"
}

resource "google_tags_tag_value" "tag_value" {
  parent     = "tagKeys/${google_tags_tag_key.tag_key.name}"
  short_name = "%{value_short_name}"
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "europe-west4"
  display_name           = "cluster-display-name-updated"

  private_cluster_config {
    allowed_projects        = ["${data.google_project.project.project_id}"]
    enable_private_endpoint = true
  }

  domain_config {
    domain = "workstations.example.com"
  }

  workstation_authorization_url = "https://console.cloud.google.com/workstations/auth-updated"
  workstation_launch_url        = "https://console.cloud.google.com/workstations/launch-updated"

  labels = {
    foo = "bar"
  }

  annotations = {
    label-one = "value-updated"
  }

  tags = {
    "${data.google_project.project.project_id}/${google_tags_tag_key.tag_key.short_name}" = "${google_tags_tag_value.tag_value.short_name}"
  }
}

resource "google_compute_network" "default" {
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "europe-west4"
  network       = google_compute_network.default.name
}
`, context)
}
