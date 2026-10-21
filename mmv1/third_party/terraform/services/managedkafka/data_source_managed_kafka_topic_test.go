package managedkafka_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/managedkafka"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccManagedKafkaTopicDatasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckManagedKafkaTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKafkaTopicDatasource_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_managed_kafka_topic.default",
						"google_managed_kafka_topic.example",
					),
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_managed_kafka_topic.without_optional",
						"google_managed_kafka_topic.example",
					),
				),
			},
		},
	})
}

func testAccManagedKafkaTopicDatasource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
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

resource "google_managed_kafka_topic" "example" {
  topic_id           = "tf-test-my-topic%{random_suffix}"
  cluster            = google_managed_kafka_cluster.cluster.cluster_id
  location           = "us-central1"
  partition_count    = 2
  replication_factor = 3
  configs = {
    "cleanup.policy" = "compact"
  }
}

data "google_project" "project" {
}

data "google_managed_kafka_topic" "default" {
  topic_id = google_managed_kafka_topic.example.topic_id
  cluster  = google_managed_kafka_topic.example.cluster
  location = google_managed_kafka_topic.example.location
  project  = google_managed_kafka_topic.example.project
}

data "google_managed_kafka_topic" "without_optional" {
  topic_id = google_managed_kafka_topic.example.topic_id
  cluster  = google_managed_kafka_topic.example.cluster
}
`, context)
}
