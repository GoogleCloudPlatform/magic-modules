package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func setTestCheckDataPipelinePipelineId(res string, pipelineId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceDataPipelinePipelineId(res, s)
		if err != nil {
			return err
		}
		*pipelineId = updateId
		return nil
	}
}

func testCheckDataPipelinePipelineIdAfterUpdate(res string, pipelineId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceDataPipelinePipelineId(res, s)
		if err != nil {
			return err
		}

		if pipelineId == nil {
			return fmt.Errorf("unexpected error, pipeline ID was not set")
		}

		if *pipelineId != updateId {
			return fmt.Errorf("unexpected mismatch in pipeline ID after update, resource was recreated. Initial %q, Updated %q",
				*pipelineId, updateId)
		}
		return nil
	}
}

func getTestResourceDataPipelinePipelineId(res string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return "", fmt.Errorf("not found: %s", res)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("no ID is set for %s", res)
	}

	if v, ok := rs.Primary.Attributes["id"]; ok {
		return v, nil
	}

	return "", fmt.Errorf("id not set on resource %s", res)
}

func TestAccDataPipelinePipeline_basic(t *testing.T) {
	t.Parallel()

	var generatedId string
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataPipelinePipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelinePipeline_basic(),
				Check:  setTestCheckDataPipelinePipelineId("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataPipelinePipeline_basicUpdate(),
				Check:  testCheckDataPipelinePipelineIdAfterUpdate("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccDataPipelinePipeline_basicUpdate() string {
	return `
resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline"
  display_name = "update-pipeline"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_flex_template_request {
      project_id = "my-project"
      launch_parameter {
        job_name = "my-job"
        environment {
          temp_location = "gs://my-bucket/tmp_dir"
        }
        container_spec_gcs_path = "gs://my-bucket/templates/template_file"
      }
      location      = "us-central1"
      validate_only = false
    }
  }
  schedule_info {
    schedule  = "* */2 * * *"
    time_zone = "UTC"
  }
}
`
}
func testAccDataPipelinePipeline_basic() string {
	return `
resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_flex_template_request {
      project_id = "my-project"
      launch_parameter {
        job_name = "my-job"
        environment {
          temp_location = "gs://my-bucket/tmp_dir"
        }
        container_spec_gcs_path = "gs://my-bucket/templates/template_file"
      }
      location      = "us-central1"
      validate_only = false
    }
  }
  schedule_info {
    schedule  = "* */2 * * *"
    time_zone = "UTC"
  }
}
`
}
