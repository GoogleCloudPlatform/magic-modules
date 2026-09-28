package geminidataanalytics_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGeminiDataAnalyticsDataAgent_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGeminiDataAnalyticsDataAgentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDataAnalyticsDataAgent_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_gemini_data_analytics_data_agent.agent", "data_analytics_agent.0.published_context.0.system_instruction", "Answer questions about sales orders."),
					resource.TestCheckResourceAttrPair("google_gemini_data_analytics_data_agent.agent", "data_analytics_agent.0.published_context.0.datasource_references.0.bq.0.table_references.0.project_id", "google_bigquery_table.table", "project"),
					resource.TestCheckResourceAttrPair("google_gemini_data_analytics_data_agent.agent", "data_analytics_agent.0.published_context.0.datasource_references.0.bq.0.table_references.0.dataset_id", "google_bigquery_dataset.dataset", "dataset_id"),
					resource.TestCheckResourceAttrPair("google_gemini_data_analytics_data_agent.agent", "data_analytics_agent.0.published_context.0.datasource_references.0.bq.0.table_references.0.table_id", "google_bigquery_table.table", "table_id"),
				),
			},
			{
				ResourceName:            "google_gemini_data_analytics_data_agent.agent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGeminiDataAnalyticsDataAgent_updated(context),
			},
			{
				ResourceName:            "google_gemini_data_analytics_data_agent.agent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccGeminiDataAnalyticsDataAgent_bigQueryResources(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id                 = "tf_test_data_agent_ds_%{random_suffix}"
  location                   = "US"
  delete_contents_on_destroy = true
}

resource "google_bigquery_table" "table" {
  dataset_id          = google_bigquery_dataset.dataset.dataset_id
  table_id            = "tf_test_sales_%{random_suffix}"
  deletion_protection = false

  schema = jsonencode([
    {
      name = "order_id"
      type = "STRING"
      mode = "REQUIRED"
    },
    {
      name = "order_total"
      type = "FLOAT"
      mode = "NULLABLE"
    }
  ])
}

resource "google_bigquery_routine" "routine" {
  dataset_id      = google_bigquery_dataset.dataset.dataset_id
  routine_id      = "tf_test_add_tax_%{random_suffix}"
  routine_type    = "SCALAR_FUNCTION"
  language        = "SQL"
  definition_body = "amount * 1.1"

  arguments {
    name      = "amount"
    data_type = jsonencode({ "typeKind" : "FLOAT64" })
  }

  return_type = jsonencode({ "typeKind" : "FLOAT64" })
}
`, context)
}

func testAccGeminiDataAnalyticsDataAgent_basic(context map[string]interface{}) string {
	return testAccGeminiDataAnalyticsDataAgent_bigQueryResources(context) + acctest.Nprintf(`
resource "google_gemini_data_analytics_data_agent" "agent" {
  location      = "global"
  data_agent_id = "tf-test-data-agent-%{random_suffix}"
  display_name  = "BigQuery data agent"
  description   = "Example Gemini Data Analytics data agent backed by BigQuery."

  labels = {
    env = "test"
  }

  data_analytics_agent {
    published_context {
      system_instruction = "Answer questions about sales orders."

      datasource_references {
        bq {
          table_references {
            project_id = google_bigquery_dataset.dataset.project
            dataset_id = google_bigquery_dataset.dataset.dataset_id
            table_id   = google_bigquery_table.table.table_id

            schema {
              description = "Customer sales orders."
              synonyms    = ["orders", "purchases"]
              tags        = ["sales"]

              fields {
                name        = "order_total"
                description = "Total order value in USD."
                synonyms    = ["amount"]
                tags        = ["currency"]
              }
            }
          }
        }
      }

      options {
        model = "LATEST_GA_MODEL"
        analysis {
          python {
            enabled = true
          }
        }
        datasource {
          big_query_max_billed_bytes = 1000000000
        }
      }

      example_queries {
        natural_language_question = "What is the total order value?"
        sql_query                 = "SELECT SUM(order_total) AS total_order_value FROM tf_test_sales_%{random_suffix}"
      }

      example_queries {
        natural_language_question = "What is the total for an order?"
        sql_query                 = "SELECT order_total FROM tf_test_sales_%{random_suffix} WHERE order_id = @order_id"
        parameters {
          name        = "order_id"
          description = "The order ID."
          data_type   = "STRING"
        }
      }

      glossary_terms {
        display_name = "Order"
        description  = "A customer purchase transaction."
        labels       = ["sales", "transaction"]
      }

      schema_relationships {
        left_schema_paths {
          table_fqn = "${google_bigquery_dataset.dataset.project}.${google_bigquery_dataset.dataset.dataset_id}.${google_bigquery_table.table.table_id}"
          paths     = ["order_id"]
        }
        right_schema_paths {
          table_fqn = "${google_bigquery_dataset.dataset.project}.${google_bigquery_dataset.dataset.dataset_id}.${google_bigquery_table.table.table_id}"
          paths     = ["order_total"]
        }
        sources          = ["LLM_SUGGESTED"]
        confidence_score = 0.9
      }

      user_functions {
        bq_routines {
          routine_reference {
            project_id = google_bigquery_routine.routine.project
            dataset_id = google_bigquery_routine.routine.dataset_id
            routine_id = google_bigquery_routine.routine.routine_id
          }
          description = "Adds 10% tax to an amount."
        }
      }
    }
  }
}
`, context)
}

func testAccGeminiDataAnalyticsDataAgent_updated(context map[string]interface{}) string {
	return testAccGeminiDataAnalyticsDataAgent_bigQueryResources(context) + acctest.Nprintf(`
resource "google_gemini_data_analytics_data_agent" "agent" {
  location      = "global"
  data_agent_id = "tf-test-data-agent-%{random_suffix}"
  display_name  = "Updated data agent"
  description   = "Updated Gemini Data Analytics data agent."

  labels = {
    env     = "test"
    updated = "true"
  }

  data_analytics_agent {
    staging_context {
      system_instruction = "Answer concise questions about the sample Shakespeare corpus."

      datasource_references {
        bq {
          table_references {
            project_id = "bigquery-public-data"
            dataset_id = "samples"
            table_id   = "shakespeare"

            schema {
              description = "Word counts in Shakespeare's works."
              synonyms    = ["plays"]
              tags        = ["literature"]

              fields {
                name        = "word_count"
                description = "Number of times the word appears in the corpus."
                synonyms    = ["count"]
                tags        = ["metric"]
              }
            }
          }
        }
      }

      options {
        analysis {
          python {
            enabled = true
          }
        }
        datasource {
          big_query_max_billed_bytes = 2000000000
        }
      }

      example_queries {
        natural_language_question = "How many times does a word appear in Hamlet?"
        sql_query                 = "SELECT SUM(word_count) FROM samples.shakespeare WHERE corpus = 'hamlet' AND word = @word"
        parameters {
          name      = "word"
          data_type = "STRING"
        }
      }

      glossary_terms {
        display_name = "Corpus"
        description  = "A collection of written works."
        labels       = ["literature", "text"]
      }

      schema_relationships {
        left_schema_paths {
          table_fqn = "bigquery-public-data.samples.shakespeare"
          paths     = ["word"]
        }
        right_schema_paths {
          table_fqn = "bigquery-public-data.samples.shakespeare"
          paths     = ["word_count"]
        }
        sources          = ["LLM_SUGGESTED"]
        confidence_score = 0.9
      }

      user_functions {
        bq_routines {
          routine_reference {
            project_id = google_bigquery_routine.routine.project
            dataset_id = google_bigquery_routine.routine.dataset_id
            routine_id = google_bigquery_routine.routine.routine_id
          }
          description = "Adds tax to an amount."
        }
      }
    }
  }
}
`, context)
}
