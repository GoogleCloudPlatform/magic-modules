package vertexai_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
)

func TestAccVertexAIEndpointDeployedModel_mutation(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	sourceModel := vertexai.BootstrapVertexAIModel(t, "us-central1")
	suffix := acctest.RandString(t, 10)
	// A fixed deployed model ID can become unusable in a project after an aborted
	// deployment (deploys then fail with INTERNAL), so use a random ID per test.
	deployedModelID := fmt.Sprintf("%d", acctest.RandIntRange(t, 1000000000, 2147483647))
	trafficSplit := fmt.Sprintf(`{"%s":100}`, deployedModelID)
	context := func(minReplicas, maxReplicas, requiredReplicas int, target interface{}, containerLogging bool, trafficSplit interface{}) map[string]interface{} {
		return map[string]interface{}{
			"project":           project,
			"source_model":      sourceModel,
			"suffix":            suffix,
			"deployed_model_id": deployedModelID,
			"min_replicas":      minReplicas,
			"max_replicas":      maxReplicas,
			"required_replicas": requiredReplicas,
			"target":            target,
			"container_logging": containerLogging,
			"traffic_split":     trafficSplit,
		}
	}
	const resourceName = "google_vertex_ai_endpoint_deployed_model.model"
	expectUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointDeployedModelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// target is omitted so the API default (60) is read back.
				Config: testAccVertexAIEndpointDeployedModel_mutation(context(1, 1, 1, nil, false, nil)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform deployed model"),
					resource.TestCheckResourceAttr(resourceName, "enable_access_logging", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_container_logging", "false"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-central1"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.required_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.autoscaling_metric_specs.0.target", "60"),
				),
			},
			{
				// Every mutable field changes together and the endpoint starts routing
				// traffic to the deployed model, without replacing the deployment.
				Config:           testAccVertexAIEndpointDeployedModel_mutation(context(2, 3, 2, 75, true, strconv.Quote(trafficSplit))),
				ConfigPlanChecks: expectUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.min_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.max_replica_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.required_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.autoscaling_metric_specs.0.target", "75"),
					resource.TestCheckResourceAttr(resourceName, "enable_container_logging", "true"),
					resource.TestCheckResourceAttr("google_vertex_ai_endpoint.endpoint", "traffic_split", trafficSplit),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// The target is set explicitly because removing it retains the
				// existing value rather than resetting to the API default.
				Config:           testAccVertexAIEndpointDeployedModel_mutation(context(1, 1, 1, 60, false, strconv.Quote(trafficSplit))),
				ConfigPlanChecks: expectUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.max_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.required_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "dedicated_resources.0.autoscaling_metric_specs.0.target", "60"),
					resource.TestCheckResourceAttr(resourceName, "enable_container_logging", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVertexAIEndpointDeployedModel_mutation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_deployed_model" "model" {
  endpoint              = google_vertex_ai_endpoint.endpoint.id
  deployed_model_id     = "%{deployed_model_id}"
  model                 = google_vertex_ai_model.model.id
  display_name          = "Terraform deployed model"
  enable_access_logging = true
  region                = "us-central1"
  enable_container_logging = %{container_logging}

  dedicated_resources {
    machine_spec {
      machine_type = "n1-standard-2"
    }

    min_replica_count      = %{min_replicas}
    max_replica_count      = %{max_replicas}
    required_replica_count = %{required_replicas}
    autoscaling_metric_specs {
      metric_name = "aiplatform.googleapis.com/prediction/online/cpu/utilization"
      target      = %{target}
    }
  }
}

resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "tf-test-deployment-%{suffix}"
  display_name = "Terraform deployed model endpoint"
  location     = "us-central1"
  region       = "us-central1"
  project      = "%{project}"
  traffic_split = %{traffic_split}
}

resource "google_vertex_ai_model" "model" {
  model_id     = "tf-test-deployment-%{suffix}"
  project      = "%{project}"
  region       = "us-central1"
  source_model = "%{source_model}"
}
`, context)
}
