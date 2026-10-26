package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIEvaluationSet_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEvaluationSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEvaluationSet_updateStep1(context),
			},
			{
				ResourceName:            "google_vertex_ai_evaluation_set.evaluation_set",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccVertexAIEvaluationSet_updateStep3(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vertex_ai_evaluation_set.evaluation_set", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_vertex_ai_evaluation_set.evaluation_set",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccVertexAIEvaluationSet_updateStep1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_evaluation_item" "item_1" {
  region               = "us-central1"
  display_name         = "tf-test-eval-item-1-%{random_suffix}"
  evaluation_item_type = "REQUEST"
  evaluation_request {
    prompt {
      text = "hello"
    }
  }
}

resource "google_vertex_ai_evaluation_set" "evaluation_set" {
  region           = "us-central1"
  display_name     = "tf-test-eval-set-%{random_suffix}"
  evaluation_items = [
    google_vertex_ai_evaluation_item.item_1.id
  ]
  metadata = jsonencode({
    "env" = "test"
  })
}
`, context)
}

func testAccVertexAIEvaluationSet_updateStep3(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_evaluation_item" "item_1" {
  region               = "us-central1"
  display_name         = "tf-test-eval-item-1-%{random_suffix}"
  evaluation_item_type = "REQUEST"
  evaluation_request {
    prompt {
      text = "hello"
    }
  }
}

resource "google_vertex_ai_evaluation_item" "item_2" {
  region               = "us-central1"
  display_name         = "tf-test-eval-item-2-%{random_suffix}"
  evaluation_item_type = "REQUEST"
  evaluation_request {
    prompt {
      text = "world"
    }
  }
}

resource "google_vertex_ai_evaluation_set" "evaluation_set" {
  region           = "us-central1"
  display_name     = "tf-test-eval-set-%{random_suffix}"
  evaluation_items = [
    google_vertex_ai_evaluation_item.item_1.id,
    google_vertex_ai_evaluation_item.item_2.id
  ]
  metadata = jsonencode({
    "env" = "test"
  })
}
`, context)
}
