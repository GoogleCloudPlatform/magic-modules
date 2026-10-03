package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Exercises the mutable fields of the OnlineEvaluator (display_name and the
// config sampling settings), which the generated example tests do not cover.
func TestAccVertexAIOnlineEvaluator_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	first := map[string]interface{}{
		"random_suffix": suffix,
		"display_name":  "tf-test-online-evaluator" + suffix,
		"max_samples":   "100",
		"percentage":    10,
	}
	second := map[string]interface{}{
		"random_suffix": suffix,
		"display_name":  "tf-test-online-evaluator-updated" + suffix,
		"max_samples":   "200",
		"percentage":    25,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIOnlineEvaluatorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIOnlineEvaluator_update(first),
			},
			{
				ResourceName:            "google_vertex_ai_online_evaluator.evaluator",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_observability.0.open_telemetry", "region"},
			},
			{
				Config: testAccVertexAIOnlineEvaluator_update(second),
			},
			{
				ResourceName:            "google_vertex_ai_online_evaluator.evaluator",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_observability.0.open_telemetry", "region"},
			},
		},
	})
}

func testAccVertexAIOnlineEvaluator_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_reasoning_engine" "engine" {
  display_name = "tf-test-reasoning-engine%{random_suffix}"
  description  = "Reasoning engine evaluated by the online evaluator"
  region       = "us-central1"
}

resource "google_vertex_ai_online_evaluator" "evaluator" {
  region         = "us-central1"
  display_name   = "%{display_name}"
  agent_resource = google_vertex_ai_reasoning_engine.engine.id

  config {
    max_evaluated_samples_per_run = "%{max_samples}"
    random_sampling {
      percentage = %{percentage}
    }
  }

  cloud_observability {
    open_telemetry {
      semconv_version = "1.39.0"
    }

    trace_scope {
      filter {
        duration {
          comparison_operator = "GREATER"
          value               = 0
        }
      }
    }
  }

  metric_sources {
    metric = jsonencode({
      predefinedMetricSpec = {
        metricSpecName = "safety_v1"
      }
    })
  }
}
`, context)
}
