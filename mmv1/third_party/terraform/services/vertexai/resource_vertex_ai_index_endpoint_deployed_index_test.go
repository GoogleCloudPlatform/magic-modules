package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIIndexEndpointDeployedIndex_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIIndexEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIIndexEndpointDeployedIndex_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.index_endpoint_deployed_index",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIIndexEndpoint_updated(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.index_endpoint_deployed_index",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIIndexEndpoint_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_vertex_ai_index_endpoint_deployed_index" "index_endpoint_deployed_index" {
		depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
		index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
		index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
		deployed_index_id = "tf-test-deployed-id%{random_suffix}"
		display_name = "deployed_index_mutated_display_name"
		dedicated_resources {
		  machine_spec {
			machine_type      = "n1-standard-32"
		  }
		  max_replica_count = 2
		  min_replica_count = 1
		}
	  }
	  
	  resource "google_storage_bucket" "bucket" {
		name     = "tf-test-bucket-name%{random_suffix}"
		location = "us-central1"
		uniform_bucket_level_access = true
	  }
	  
	  # The sample data comes from the following link:
	  # https://cloud.google.com/vertex-ai/docs/matching-engine/filtering#specify-namespaces-tokens
	  resource "google_storage_bucket_object" "data" {
		name   = "tf-test-storage-bucket-name%{random_suffix}"
		bucket = google_storage_bucket.bucket.name
		content = <<EOF
	  {"id": "42", "embedding": [0.5, 1.0], "restricts": [{"namespace": "class", "allow": ["cat", "pet"]},{"namespace": "category", "allow": ["feline"]}]}
	  {"id": "43", "embedding": [0.6, 1.0], "restricts": [{"namespace": "class", "allow": ["dog", "pet"]},{"namespace": "category", "allow": ["canine"]}]}
	  EOF
	  }
	  
	  resource "google_vertex_ai_index" "index" {
		labels = {
		  foo = "bar"
		}
		region   = "us-central1"
		display_name = "tf-test-index%{random_suffix}"
		description = "index for test"
		metadata {
		  contents_delta_uri = "gs://${google_storage_bucket.bucket.name}/contents"
		  config {
			dimensions = 2
			approximate_neighbors_count = 150
			shard_size = "SHARD_SIZE_MEDIUM"
			distance_measure_type = "DOT_PRODUCT_DISTANCE"
			algorithm_config {
			  tree_ah_config {
				leaf_node_embedding_count = 500
				leaf_nodes_to_search_percent = 7
			  }
			}
		  }
		}
		index_update_method = "BATCH_UPDATE"
	  }
	  
	  
	  resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
		display_name = "tf-test-index-endpoint%{random_suffix}"
		description  = "A sample vertex endpoint"
		region       = "us-central1"
		labels       = {
		  label-one = "value-one"
		}
		public_endpoint_enabled = true
	  }
	  
	  data "google_project" "project" {}
`, context)
}

func testAccVertexAIIndexEndpoint_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_vertex_ai_index_endpoint_deployed_index" "index_endpoint_deployed_index" {
		depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
		index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
		index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
		deployed_index_id = "tf-test-deployed-id%{random_suffix}"
		display_name = "deployed_index_mutated_display_name"
		dedicated_resources {
		  machine_spec {
			machine_type      = "n1-standard-32"
		  }
		  max_replica_count = 3
		  min_replica_count = 2
		}
	  }
	  
	  resource "google_storage_bucket" "bucket" {
		name     = "tf-test-bucket-name%{random_suffix}"
		location = "us-central1"
		uniform_bucket_level_access = true
	  }
	  
	  # The sample data comes from the following link:
	  # https://cloud.google.com/vertex-ai/docs/matching-engine/filtering#specify-namespaces-tokens
	  resource "google_storage_bucket_object" "data" {
		name   = "tf-test-storage-bucket-name%{random_suffix}"
		bucket = google_storage_bucket.bucket.name
		content = <<EOF
	  {"id": "42", "embedding": [0.5, 1.0], "restricts": [{"namespace": "class", "allow": ["cat", "pet"]},{"namespace": "category", "allow": ["feline"]}]}
	  {"id": "43", "embedding": [0.6, 1.0], "restricts": [{"namespace": "class", "allow": ["dog", "pet"]},{"namespace": "category", "allow": ["canine"]}]}
	  EOF
	  }
	  
	  resource "google_vertex_ai_index" "index" {
		labels = {
		  foo = "bar"
		}
		region   = "us-central1"
		display_name = "tf-test-index%{random_suffix}"
		description = "index for test"
		metadata {
		  contents_delta_uri = "gs://${google_storage_bucket.bucket.name}/contents"
		  config {
			dimensions = 2
			approximate_neighbors_count = 150
			shard_size = "SHARD_SIZE_MEDIUM"
			distance_measure_type = "DOT_PRODUCT_DISTANCE"
			algorithm_config {
			  tree_ah_config {
				leaf_node_embedding_count = 500
				leaf_nodes_to_search_percent = 7
			  }
			}
		  }
		}
		index_update_method = "BATCH_UPDATE"
	  }
	  
	  
	  resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
		display_name = "tf-test-index-endpoint%{random_suffix}"
		description  = "A sample vertex endpoint"
		region       = "us-central1"
		labels       = {
		  label-one = "value-one"
		}
		public_endpoint_enabled = true
	  }
	  
	  data "google_project" "project" {}
`, context)
}
