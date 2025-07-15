package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIEndpointWithModelGardenDeployment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             nil, // the resource does not yet support deletion
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_basic(context),
			},
			{
				ResourceName:      "google_vertex_ai_endpoint_with_model_garden_deployment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"etag",
					"id",
					"publisher_model_name",
					"project",
					"model_config",
				},
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "test" {
  publisher_model_name = "publishers/google/models/paligemma@paligemma-224-float32"
  location             = "us-central1"
  model_config {
    accept_eula =  true
  }
}
`, context)
}

func TestAccVertexAIEndpointWithModelGardenDeployment_withConfigs(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_withConfigs(context),
			},
			{
				ResourceName:      "google_vertex_ai_endpoint_with_model_garden_deployment.test_with_configs",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"etag",
					"id",
					"publisher_model_name",
					"project",
					"model_config",
					"deploy_config",
				},
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_withConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "test_with_configs" {
  publisher_model_name = "publishers/google/models/paligemma@paligemma-224-float32"
  location             = "us-central1"
  model_config {
    accept_eula =  true
  }
  deploy_config {
    dedicated_resources {
      machine_spec {
        machine_type      = "g2-standard-16"
        accelerator_type  = "NVIDIA_L4"
        accelerator_count = 1
      }
      min_replica_count = 1
    }
  }
}
`, context)
}
