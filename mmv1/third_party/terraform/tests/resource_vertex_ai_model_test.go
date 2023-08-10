package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVertexAIModel_vertexAiModelAdvancedExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_vertexAiModelAdvancedExample(context),
			},
			{
				ResourceName:            "google_vertex_ai_model.model",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccVertexAIModel_vertexAiModelAdvancedExampleUpdate(context),
			},
			{
				ResourceName:            "google_vertex_ai_model.model",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccVertexAIModel_vertexAiModelAdvancedExampleUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_vertex_ai_model" "model" {
  name = "tf-test-model-name%{random_suffix}"
  container_spec {
    image_uri = "us-docker.pkg.dev/vertex-ai/prediction/xgboost-cpu.1-5:latest"
    args      = ["sample", "args"]
    command   = ["sample", "command"]
    env {
      name  = "env_one"
      value = "value_one"
    }
    health_route = "/health"
    ports {
      container_port = 8080
    }
    predict_route = "/predict"
  }
  display_name = "new-sample-model"
  region       = "us-central1"
  artifact_uri = "gs://cloud-samples-data/vertex-ai/google-cloud-aiplatform-ci-artifacts/models/iris_xgboost/"
  description  = "An updated sample model"
  labels = {
    label-two = "value-two"
  }
  version_aliases     = ["default", "v1", "v2"]
  version_description = "A sample model version"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}
`, context)
}
