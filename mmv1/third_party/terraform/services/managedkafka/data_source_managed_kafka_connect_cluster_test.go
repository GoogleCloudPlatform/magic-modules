package managedkafka_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/managedkafka"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccManagedKafkaConnectClusterDatasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckManagedKafkaConnectClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKafkaConnectClusterDatasource_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_managed_kafka_connect_cluster.default",
						"google_managed_kafka_connect_cluster.example",
					),
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_managed_kafka_connect_cluster.without_optional",
						"google_managed_kafka_connect_cluster.example",
					),
				),
			},
		},
	})
}

func testAccManagedKafkaConnectClusterDatasource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_compute_subnetwork" "mkc_secondary_subnet" {
  project       = data.google_project.project.project_id
  name          = "tf-test-my-secondary-subnetwork%{random_suffix}"
  ip_cidr_range = "10.3.0.0/16"
  region        = "us-central1"
  network       = "default"
}

resource "google_managed_kafka_cluster" "cluster" {
  cluster_id = "tf-test-my-cluster%{random_suffix}"
  location   = "us-central1"
  capacity_config {
    vcpu_count   = 3
    memory_bytes = 3221225472
  }
  gcp_config {
    access_config {
      network_configs {
        subnet = "projects/${data.google_project.project.number}/regions/us-central1/subnetworks/default"
      }
    }
  }
}

resource "google_managed_kafka_connect_cluster" "example" {
  connect_cluster_id = "tf-test-my-connect-cluster%{random_suffix}"
  kafka_cluster      = "projects/${data.google_project.project.project_id}/locations/us-central1/clusters/${google_managed_kafka_cluster.cluster.cluster_id}"
  location           = "us-central1"

  capacity_config {
    vcpu_count   = 12
    memory_bytes = 21474836480
  }
  gcp_config {
    access_config {
      network_configs {
        primary_subnet     = "projects/${data.google_project.project.number}/regions/us-central1/subnetworks/default"
        additional_subnets = [google_compute_subnetwork.mkc_secondary_subnet.id]
        dns_domain_names   = ["${google_managed_kafka_cluster.cluster.cluster_id}.us-central1.managedkafka.${data.google_project.project.project_id}.cloud.goog"]
      }
    }
  }
  labels = {
    key = "value"
  }
}

data "google_managed_kafka_connect_cluster" "default" {
  connect_cluster_id = google_managed_kafka_connect_cluster.example.connect_cluster_id
  location           = google_managed_kafka_connect_cluster.example.location
  project            = google_managed_kafka_connect_cluster.example.project
}

data "google_managed_kafka_connect_cluster" "without_optional" {
  connect_cluster_id = google_managed_kafka_connect_cluster.example.connect_cluster_id
}
`, context)
}
