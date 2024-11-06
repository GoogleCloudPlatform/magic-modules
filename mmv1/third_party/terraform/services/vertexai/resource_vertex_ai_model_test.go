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

func TestAccVertexAIModel_modelIdNotProvided(t *testing.T) {
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
				Config: testAccVertexAIModel_modelIdNotProvided(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_vertex_ai_model.model", "model_id"),
				),
			},
		},
	})
}

func testAccVertexAIModel_modelIdNotProvided(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/6033738282699849728"

  region       = "us-central1"
}
`, context)
}

func TestAccVertexAIModel_modelIdProvided(t *testing.T) {
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
				Config: testAccVertexAIModel_modelIdProvided(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_model.model", "model_id", context["model_id"].(string)),
				),
			},
		},
	})
}

func testAccVertexAIModel_modelIdProvided(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model" "model" {
  model_id = "%{model_id}"
  project = "%{project_name}"
  source_model = "projects/%{project_name}/locations/us-central1/models/6033738282699849728"

  region       = "us-central1"
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
