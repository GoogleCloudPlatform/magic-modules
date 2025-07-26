package vertexai_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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

  region       = "us-central1"
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

func testAccVertexAIModel_modelIdProvided_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region       = "us-central1"
}
`, context)
}

func testAccVertexAIModel_modelIdProvided_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/7222055265628061696"

  region       = "us-central1"

  description = "updated"
  display_name = "updated"
}
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

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VertexAIBasePath}}projects/{{project}}/locations/{{region}}/models/{{model_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VertexAIModel still exists at %s", url)
			}
		}

		return nil
	}
}
