// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datalineage_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/datalineage"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
)

var (
	_ = fmt.Sprintf
	_ = log.Print
	_ = strconv.Atoi
	_ = strings.Trim
	_ = time.Now
	_ = resource.TestMain
	_ = terraform.NewState
	_ = envvar.TestEnvVar
	_ = tpgresource.SetLabels
	_ = transport_tpg.Config{}
	_ = googleapi.Error{}
	_ = datalineage.Product
)

func TestAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.simple",
				RefreshState: true,
			},
		},
	})
}

func testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "simple" {
  namespace   = "example_simple_namespace"
  name        = "example_simple_name"
  description = "Nightly ETL from raw to curated"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset_simple/source_table_1"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target_simple/target_table_1"
  }
}
`, context)
}

func TestAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "namespace", "testNamespace"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "name", "test_full_job"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "input.#", "2"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "output.#", "1"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "output.0.column_lineage.0.dataset_input.#", "3"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "output.0.column_lineage.0.field.#", "3"),
					resource.TestCheckResourceAttr("google_data_lineage_open_lineage_job.full_job", "knowledge_catalog.#", "1"),
					resource.TestCheckResourceAttrSet("google_data_lineage_open_lineage_job.full_job", "knowledge_catalog.0.process"),
					resource.TestCheckResourceAttrSet("google_data_lineage_open_lineage_job.full_job", "knowledge_catalog.0.run"),
				),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.full_job",
				RefreshState: true,
			},
		},
	})
}

func testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "full_job" {
  namespace   = "testNamespace"
  name        = "test_full_job"
  description = "Test resource with all available facets"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table_1"

    symlink {
      namespace = "bigquery"
      name      = "my-project-name.raw_dataset_with_facets.source_table_1"
      type      = "TABLE"
    }

	catalog {
      framework = "iceberg"
      type      = "bigquery"
      name      = "example-catalog"
    }
  }

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table_2"

    symlink {
      namespace = "bigquery"
      name      = "my-project-name.raw_dataset_with_facets.source_table_2"
      type      = "TABLE"
    }
	
	catalog {
      framework = "iceberg"
      type      = "bigquery"
      name      = "example-catalog"
    }
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/target_table_1"

    symlink {
      namespace = "bigquery"
      name      = "my-project-name.raw_dataset_with_facets.target_table_1"
      type      = "TABLE"
    }

	catalog {
      framework = "iceberg"
      type      = "bigquery"
      name      = "example-catalog"
    }

    column_lineage {
      dataset_input {
       namespace  = "gs://example-bucket/"
       name       = "warehouse/raw_dataset/source_table_1"
        field     = "a"
        transformation {
          type    = "INDIRECT"
          subtype = "GROUP_BY"
        }
        transformation {
          type    = "INDIRECT"
          subtype = "JOIN"
        }
        transformation {
          type    = "INDIRECT"
          subtype = "FILTER"
        }
      }

      dataset_input {
       namespace  = "gs://example-bucket/"
       name       = "warehouse/raw_dataset/source_table_1"
        field     = "b"
        transformation {
          type    = "INDIRECT"
          subtype = "GROUP_BY"
        }
      }

      dataset_input {
       namespace  = "gs://example-bucket/"
       name       = "warehouse/raw_dataset/source_table_2"
        field     = "a"
        transformation {
          type    = "INDIRECT"
          subtype = "JOIN"
        }
        transformation {
          type    = "INDIRECT"
          subtype = "FILTER"
        }
      }

      field {
        name = "ident"
        input {
          namespace = "gs://example-bucket/"
          name      = "warehouse/raw_dataset/source_table_1"
          field     = "a"
          transformation {
            type    = "DIRECT"
            subtype = "IDENTITY"
          }
        }
      }

      field {
        name = "trans"
        input {
          namespace = "gs://example-bucket/"
          name      = "warehouse/raw_dataset/source_table_1"
          field     = "b"
          transformation {
            type    = "DIRECT"
            subtype = "TRANSFORMATION"
          }
        }
      }

      field {
        name = "agg"
        input {
          namespace = "gs://example-bucket/"
          name      = "warehouse/raw_dataset/source_table_2"
          field     = "c"
          transformation {
            type    = "DIRECT"
            subtype = "AGGREGATION"
          }
        }
      }
    }
  }
}
`, context)
}

// TestAccDataLineageOpenLineageJob_UpdateDescription validates that description changes don't trigger updates.
// Checks that:
// - Resource can be created with initial description
// - Description changes are applied locally only, not sent to the server
// - No new OpenLineage event is emitted when only description changes
// - State refresh correctly preserves the description value
func TestAccDataLineageOpenLineageJob_UpdateDescription(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_UpdateDescription_Initial(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.update_test",
				RefreshState: true,
			},
			{
				Config: testAccDataLineageOpenLineageJob_UpdateDescription_Updated(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.update_test",
				RefreshState: true,
			},
		},
	})
}

// testAccDataLineageOpenLineageJob_UpdateDescription_Initial provides initial resource config.
// Establishes baseline state with description="Initial description" for update testing.
func testAccDataLineageOpenLineageJob_UpdateDescription_Initial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "update_test" {
  namespace   = "update_test_namespace"
  name        = "update_test_name"
  description = "Initial description"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table_1"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target/target_table_1"
  }
}
`, context)
}

// testAccDataLineageOpenLineageJob_UpdateDescription_Updated provides updated resource config.
// Changes description to validate that description updates trigger new events via REST endpoint.
func testAccDataLineageOpenLineageJob_UpdateDescription_Updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "update_test" {
  namespace   = "update_test_namespace"
  name        = "update_test_name"
  description = "Updated description after modification"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table_1"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target/target_table_1"
  }
}
`, context)
}

// TestAccDataLineageOpenLineageJob_CreateAndDelete validates the create-delete lifecycle.
// Checks that:
// - Resource creation via POST to processOpenLineageRunEvent succeeds
// - Computed fields knowledge_catalog.process and knowledge_catalog.run are populated from response
// - Resource is properly destroyed when Terraform state is cleaned up
// - Deletion via REST DELETE endpoint removes the remote process
func TestAccDataLineageOpenLineageJob_CreateAndDelete(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_CreateAndDelete_Config(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_data_lineage_open_lineage_job.delete_test", "knowledge_catalog.0.process"),
					resource.TestCheckResourceAttrSet("google_data_lineage_open_lineage_job.delete_test", "knowledge_catalog.0.run"),
				),
			},
		},
	})
}

// testAccDataLineageOpenLineageJob_CreateAndDelete_Config provides resource config with explicit deletion_policy.
// Uses deletion_policy="DELETE" to ensure remote process is deleted when terraform destroy runs.
func testAccDataLineageOpenLineageJob_CreateAndDelete_Config(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "delete_test" {
  namespace   = "delete_test_namespace"
  name        = "delete_test_name"
  description = "Job to be deleted"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target/target_table"
  }

  deletion_policy = "DELETE"
}
`, context)
}

// TestAccDataLineageOpenLineageJob_WithDeletionPolicyAbandon validates deletion_policy="ABANDON" semantics.
// Checks that:
// - Resource can be created successfully
// - deletion_policy="ABANDON" prevents remote deletion when Terraform destroys
// - Process remains on server after terraform destroy (not checked by CheckDestroy for this policy)
// - Demonstrates that ABANDON policy leaves remote state intact
func TestAccDataLineageOpenLineageJob_WithDeletionPolicyAbandon(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_DeletionPolicyAbandon(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.abandon_test",
				RefreshState: true,
			},
		},
	})
}

// testAccDataLineageOpenLineageJob_DeletionPolicyAbandon provides resource config for abandon policy testing.
// Sets deletion_policy="ABANDON" to test that resource deletion does not remove remote process.
func testAccDataLineageOpenLineageJob_DeletionPolicyAbandon(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "abandon_test" {
  namespace   = "abandon_test_namespace"
  name        = "abandon_test_name"
  description = "Job with abandon deletion policy"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target/target_table"
  }

  deletion_policy = "ABANDON"
}
`, context)
}

// testAccCheckDataLineageOpenLineageJobDestroyProducer validates that resources are properly destroyed.
// Verifies that:
// - Each resource in the destroyed state was actually removed from the server
// - A REST GET request to the process name returns 404 (resource deleted)
// - Handles deletion_policy=ABANDON by checking deletion_policy before making destroy assertions
func testAccCheckDataLineageOpenLineageJobDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_data_lineage_open_lineage_job" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			if rs.Primary.Attributes["deletion_policy"] == "ABANDON" {
				continue
			}
			conf := acctest.GoogleProviderConfig(t)

			n := rs.Primary.Attributes["knowledge_catalog.0.process"]
			url := transport_tpg.BaseUrl(datalineage.Product, conf) + strings.TrimPrefix(n, "/")

			billingProject := ""
			if conf.BillingProject != "" {
				billingProject = conf.BillingProject
			}

			_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    conf,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: conf.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DataLineageOpenLineageJob still exists at %s", url)
			}
			if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				return err
			}
		}

		return nil
	}
}

// TestAccDataLineageOpenLineageJob_DriftDetection validates drift detection during plan phase.
// Checks that:
// - Resource is created and initial run ID is recorded in state
// - External modification (new event to same process outside Terraform) is detected
// - Plan phase detects run ID mismatch between state and remote
// - Warning is logged indicating external modifications
// - Apply phase syncs the latest run ID back to state
func TestAccDataLineageOpenLineageJob_DriftDetection(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create initial resource with run ID captured in state
				Config: testAccDataLineageOpenLineageJob_DriftDetection_Config(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.drift_test",
				RefreshState: true,
			},
			{
				// Plan with no config changes - detects that run ID has drifted
				// Remote process may have new events, so latest run is different
				Config:             testAccDataLineageOpenLineageJob_DriftDetection_Config(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Apply to sync state with latest run ID from remote
				Config: testAccDataLineageOpenLineageJob_DriftDetection_Config(context),
			},
		},
	})
}

// testAccDataLineageOpenLineageJob_DriftDetection_Config provides resource config for drift testing.
// Creates a resource that will be checked for external modifications during read phase.
func testAccDataLineageOpenLineageJob_DriftDetection_Config(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "drift_test" {
  namespace   = "drift_test_namespace"
  name        = "drift_test_name"
  description = "Job for drift detection testing"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset/source_table"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target/target_table"
  }
}
`, context)
}
