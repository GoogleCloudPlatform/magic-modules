package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVertexAIModel_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_model.model",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccVertexAIModel_advanced(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_advanced(context),
			},
			{
				ResourceName:            "google_vertex_ai_model.model",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccVertexAIModel_advancedUpdate(context),
			},
			{
				ResourceName:            "google_vertex_ai_model.model",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccVertexAIModel_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_vertex_ai_model" "model" {
  name = "tf-test-model-name%{random_suffix}"
  container_spec {
    image_uri = "gcr.io/cloud-ml-service-public/cloud-ml-online-prediction-model-server-cpu:v1_15py3cmle_op_images_20200229_0210_RC00"
  }
  display_name = "sample-model"
  location     = "us-central1"
}
`, context)
}

func testAccVertexAIModel_advanced(context map[string]interface{}) string {
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
  display_name = "sample-model"
  location     = "us-central1"
  artifact_uri = "gs://cloud-samples-data/vertex-ai/google-cloud-aiplatform-ci-artifacts/models/iris_xgboost/"
  description  = "A sample model"
  labels = {
    label-one = "value-one"
  }
  version_aliases     = ["default", "v1", "v2"]
  version_description = "A sample model version"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
  depends_on   = [
    google_kms_crypto_key_iam_member.crypto_key
  ]
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIModel_advancedUpdate(context map[string]interface{}) string {
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
  location     = "us-central1"
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
  depends_on   = [
    google_kms_crypto_key_iam_member.crypto_key
  ]
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

data "google_project" "project" {}
`, context)
}

func testAccCheckVertexAIModelDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_model" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{VertexAIBasePath}}projects/{{project}}/locations/{{location}}/models/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("VertexAIModel still exists at %s", url)
			}
		}

		return nil
	}
}
