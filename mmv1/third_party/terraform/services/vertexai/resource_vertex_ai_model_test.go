package vertexai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVertexAIModel_postCreationUpdates(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"model_id":      fmt.Sprintf("tf-test-test-model%s", randomString),
		"random_suffix": randomString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_modelIdProvided_create(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_id", context["model_id"].(string)),
				),
			},
			{
				Config: testAccVertexAIModel_modelIdProvided_update(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_id", context["model_id"].(string)),
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "description", "updated"),
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_name", "updated"),
				),
			},
		},
	})
}

func TestAccVertexAIModel_modelIdNotProvidedAtCreateTime(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_modelIdNotProvided_create(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_vertex_ai_model.model", "model_id"),
				),
			},
		},
	})
}

func testAccVertexAIModel_modelIdNotProvided_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  project = "%{project_name}"
  source_model = google_vertex_ai_model.model_upload.name

  region       = "us-central1"
}

resource "google_vertex_ai_model" "model_upload" {
  model_name = "tf_test_model_upload_source_%{random_suffix}"
  description  = "basic upload model for source testing"
  region       = "us-central1"

    container_spec {
      image_uri             = "us-docker.pkg.dev/vertex-ai/prediction/tf2-cpu.2-12:latest"
    command     = ["/usr/bin/python3"]
    args        = ["model_server.py"]
    
    env {
      name  = "MODEL_NAME"
      value = "example-model"
    }
    
    health_route = "/health"
    predict_route = "/predict"
    ports {
      container_port = 8080
    }
  }
  artifact_uri = "gs://cloud-samples-data/ai-platform/mnist_tfrecord/pretrained"

}
`, context)
}

func TestAccVertexAIModel_modelIdProvidedAtCreateTime(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"model_id":      fmt.Sprintf("tf-test-test-model%s", randomString),
		"random_suffix": randomString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_modelIdProvided_create(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_id", context["model_id"].(string)),
				),
			},
		},
	})
}

func testAccVertexAIModel_modelIdProvided_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = google_vertex_ai_model.model_upload.name

  region       = "us-central1"
}

resource "google_vertex_ai_model" "model_upload" {
  model_name = "tf_test_model_upload_source_%{random_suffix}"
  description  = "basic upload model for source testing"
  region       = "us-central1"

  container_spec {
      image_uri             = "us-docker.pkg.dev/vertex-ai/prediction/tf2-cpu.2-12:latest"
    command     = ["/usr/bin/python3"]
    args        = ["model_server.py"]
    
    env {
      name  = "MODEL_NAME"
      value = "example-model"
    }
    
    health_route = "/health"
    predict_route = "/predict"
    ports {
      container_port = 8080
    }
  }
  artifact_uri = "gs://cloud-samples-data/ai-platform/mnist_tfrecord/pretrained"

}
`, context)
}

func testAccVertexAIModel_modelIdProvided_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = google_vertex_ai_model.model_upload.name

  region       = "us-central1"

  description = "updated"
  model_name = "updated"
}

resource "google_vertex_ai_model" "model_upload" {
  model_name = "tf_test_model_upload_source_%{random_suffix}"
  description  = "basic upload model for source testing"
  region       = "us-central1"

  container_spec {
      image_uri             = "us-docker.pkg.dev/vertex-ai/prediction/tf2-cpu.2-12:latest"
    command     = ["/usr/bin/python3"]
    args        = ["model_server.py"]
    
    env {
      name  = "MODEL_NAME"
      value = "example-model"
    }
    
    health_route = "/health"
    predict_route = "/predict"
    ports {
      container_port = 8080
    }
  }
  artifact_uri = "gs://cloud-samples-data/ai-platform/mnist_tfrecord/pretrained"

}
`, context)
}
