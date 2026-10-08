package vertexai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccVertexAIModel_postCreationUpdates(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project_name": envvar.GetTestProjectFromEnv(),
		"model_id":     fmt.Sprintf("tf-test-test-model%s", randomString),
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
				ResourceName:      "google_vertex_ai_model.model",
				ImportState:       true,
				ImportStateVerify: true,
				// The API returns the project number, and source_model is
				// write-only because it is only used when copying a model.
				ImportStateVerifyIgnore: []string{"project", "source_model"},
			},
			{
				Config: testAccVertexAIModel_modelIdProvided_update(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_id", context["model_id"].(string)),
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "description", "updated"),
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "display_name", "updated"),
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
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region = "us-central1"
}
`, context)
}

func TestAccVertexAIModel_modelIdProvidedAtCreateTime(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project_name": envvar.GetTestProjectFromEnv(),
		"model_id":     fmt.Sprintf("tf-test-test-model%s", randomString),
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

func TestAccVertexAIModel_copyWithOptionalFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"kms_key_name":    kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-vertex-ai-model-key").CryptoKey.Name,
		"model_id":        fmt.Sprintf("tf-test-test-model%s", acctest.RandString(t, 10)),
		"version_alias":   "tf-test-alias",
	}

	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModel_copyWithOptionalFields(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "custom_service_account", context["service_account"].(string)),
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "encryption_spec.0.kms_key_name", context["kms_key_name"].(string)),
				),
			},
			{
				Config: testAccVertexAIModel_copyWithOptionalFieldsUpdate(context),
				Check: resource.TestCheckResourceAttr(
					// Vertex AI retains the reserved "default" alias at index 0.
					"google_vertex_ai_model.model", "version_aliases.1", context["version_alias"].(string),
				),
			},
		},
	})
}

func testAccVertexAIModel_copyWithOptionalFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region = "us-central1"
  custom_service_account = "%{service_account}"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}
`, context)
}

func testAccVertexAIModel_copyWithOptionalFieldsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region = "us-central1"
  custom_service_account = "%{service_account}"
  version_aliases = ["%{version_alias}"]

  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}
`, context)
}

func testAccVertexAIModel_modelIdProvided_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region = "us-central1"
}
`, context)
}

func testAccVertexAIModel_modelIdProvided_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region = "us-central1"

  description = "updated"
  display_name = "updated"
}
`, context)
}
