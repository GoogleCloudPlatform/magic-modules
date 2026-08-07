package vertexai_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVertexAIAgentAnomalyDetectionScope(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic": testAccVertexAIAgentAnomalyDetectionScope_basicTest,
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccVertexAIAgentAnomalyDetectionScope_basicTest(t *testing.T) {
	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
			Role:   "roles/logging.configWriter",
		},
	})

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"project":  envvar.GetTestProjectFromEnv(),
		"scope_id": "tf-test-scope-id" + randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIAgentAnomalyDetectionScopeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIAgentAnomalyDetectionScope_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_agent_anomaly_detection_scope.agent_anomaly_detection_scope",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agent_anomaly_detection_scope_id", "region"},
			},
			{
				ResourceName:       "google_vertex_ai_agent_anomaly_detection_scope.agent_anomaly_detection_scope",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccVertexAIAgentAnomalyDetectionScope_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_agent_anomaly_detection_scope" "agent_anomaly_detection_scope" {
  agent_anomaly_detection_scope_id = "%{scope_id}"
  region                           = "us-central1"
  display_name                     = "sample-scope"
  log_buckets                      = [
    "projects/%{project}/locations/us-central1/buckets/_Default"
  ]
  observability_buckets            = [
    "projects/%{project}/locations/us-central1/buckets/_Default"
  ]
}
`, context)
}

func testAccCheckVertexAIAgentAnomalyDetectionScopeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_agent_anomaly_detection_scope" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "https://us-central1-aiplatform.googleapis.com/v1beta1/{{name}}")
			if err != nil {
				return err
			}
			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VertexAIAgentAnomalyDetectionScope still exists at %s", url)
			}
		}
		return nil
	}
}
