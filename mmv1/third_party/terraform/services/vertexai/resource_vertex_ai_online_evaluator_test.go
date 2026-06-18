// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vertexai_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVertexAIOnlineEvaluator_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckVertexAIOnlineEvaluatorDestroyHandwritten(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIOnlineEvaluator_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "display_name", "eval-"+randomSuffix),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "config.0.max_evaluated_samples_per_run", "100"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "config.0.random_sampling.0.percentage", "10"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "metric_sources.#", "1"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "cloud_observability.0.open_telemetry.0.semconv_version", "1.39.0"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "cloud_observability.0.trace_scope.0.filter.#", "2"),
				),
			},
			{
				ResourceName:            "google_vertex_ai_online_evaluator.evaluator",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAIOnlineEvaluator_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "display_name", "eval-upd-"+randomSuffix),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "config.0.max_evaluated_samples_per_run", "200"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "config.0.random_sampling.0.percentage", "20"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "metric_sources.#", "2"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "cloud_observability.0.open_telemetry.0.semconv_version", "1.39.0"),
					resource.TestCheckResourceAttr("google_vertex_ai_online_evaluator.evaluator", "cloud_observability.0.trace_scope.0.filter.#", "2"),
				),
			},
			{
				ResourceName:            "google_vertex_ai_online_evaluator.evaluator",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
		},
	})
}

func testAccVertexAIOnlineEvaluator_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_reasoning_engine" "engine" {
  display_name = "re-%{random_suffix}"
  description  = "A basic reasoning engine"
  labels       = {
    "key" = "value"
  }
  region       = "us-central1"
}

resource "google_vertex_ai_online_evaluator" "evaluator" {
  region         = "us-central1"
  display_name   = "eval-%{random_suffix}"
  
  agent_resource = google_vertex_ai_reasoning_engine.engine.id

  config {
    max_evaluated_samples_per_run = "100"
    random_sampling {
      percentage = 10
    }
  }

  metric_sources {
    metric = jsonencode({
      "predefinedMetricSpec" = {
        "metricSpecName" = "safety_v1"
      }
    })
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
      filter {
        total_token_usage {
          comparison_operator = "LESS"
          value               = 1000
        }
      }
    }
  }
}
`, context)
}

func testAccVertexAIOnlineEvaluator_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_reasoning_engine" "engine" {
  display_name = "re-%{random_suffix}"
  description  = "A basic reasoning engine"
  labels       = {
    "key" = "value"
  }
  region       = "us-central1"
}

resource "google_vertex_ai_online_evaluator" "evaluator" {
  region         = "us-central1"
  display_name   = "eval-upd-%{random_suffix}"
  
  agent_resource = google_vertex_ai_reasoning_engine.engine.id

  config {
    max_evaluated_samples_per_run = "200"
    random_sampling {
      percentage = 20
    }
  }

  metric_sources {
    metric = jsonencode({
      "predefinedMetricSpec" = {
        "metricSpecName" = "safety_v1"
      }
    })
  }

  metric_sources {
    metric = jsonencode({
      "predefinedMetricSpec" = {
        "metricSpecName" = "hallucination_v1"
      }
    })
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
      filter {
        total_token_usage {
          comparison_operator = "LESS"
          value               = 1000
        }
      }
    }
  }
}
`, context)
}

func testAccCheckVertexAIOnlineEvaluatorDestroyHandwritten(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_online_evaluator" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(vertexai.Product, config)+"projects/{{project}}/locations/{{region}}/onlineEvaluators/{{name}}")
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
				return fmt.Errorf("VertexAIOnlineEvaluator still exists at %s", url)
			}
		}

		return nil
	}
}
