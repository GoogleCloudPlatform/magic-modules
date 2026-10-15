package vertexai_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIPipelineJob_full(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"pipeline_job_id": fmt.Sprintf("tf-test-full-job-%s", randomSuffix),
		"storage_bucket":  fmt.Sprintf("tf-test-full-bucket-%s", randomSuffix),
		"network_name":    fmt.Sprintf("tf-test-full-net-%s", randomSuffix),
		"random_suffix":   randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIPipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIPipelineJob_full(context),
			},
			{
				ResourceName:            "google_vertex_ai_pipeline_job.pipeline_job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"end_time", "error", "job_detail", "labels", "pipeline_job_id", "pipeline_spec", "preflight_validations", "region", "start_time", "state", "terraform_labels", "update_time"},
			},
		},
	})
}

func testAccVertexAIPipelineJob_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_storage_bucket" "bucket" {
  name                        = "%{storage_bucket}"
  location                    = "us-central1"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_compute_network" "my_network" {
  name                    = "%{network_name}"
  auto_create_subnetworks = false
}

resource "google_vertex_ai_pipeline_job" "pipeline_job" {
  pipeline_job_id       = "%{pipeline_job_id}"
  display_name          = "test-pipeline-job-full"
  region                = "us-central1"
  service_account       = "${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  network               = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.my_network.name}"
  reserved_ip_ranges    = ["vertex-ai-ip-range"]
  preflight_validations = false

  labels = {
    test-label = "test-value"
  }

  psc_interface_config {
    network_attachment = "projects/${data.google_project.project.project_id}/regions/us-central1/networkAttachments/my-attachment"
    dns_peering_configs {
      domain         = "my-internal-domain.corp."
      target_network = google_compute_network.my_network.id
      target_project = data.google_project.project.project_id
    }
  }

  pipeline_spec = jsonencode({
    pipelineInfo = {
      name = "hello-world"
    }
    root = {
      dag = {
        tasks = {}
      }
    }
    schemaVersion = "2.1.0"
    sdkVersion    = "kfp-2.0.0"
  })

  runtime_config {
    gcs_output_directory = "gs://${google_storage_bucket.bucket.name}/pipeline_root"
    failure_policy       = "PIPELINE_FAILURE_POLICY_FAIL_FAST"
    parameter_values     = {
      param1 = "val1"
    }
    input_artifacts {
      key         = "my-artifact"
      artifact_id = "test-artifact"
    }
  }
}
`, context)
}

func TestAccVertexAIPipelineJob_encryptionSpec(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"pipeline_job_id": fmt.Sprintf("tf-test-enc-job-%s", randomSuffix),
		"storage_bucket":  fmt.Sprintf("tf-test-enc-bucket-%s", randomSuffix),
		"kms_key_name":    "invalid-key",
		"random_suffix":   randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIPipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVertexAIPipelineJob_encryptionSpec(context),
				ExpectError: regexp.MustCompile("The provided KMS crypto key reference is not valid|in_use_customer_managed_encryption_keys"),
			},
		},
	})
}

func testAccVertexAIPipelineJob_encryptionSpec(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%{storage_bucket}"
  location                    = "us-central1"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_vertex_ai_pipeline_job" "pipeline_job" {
  pipeline_job_id = "%{pipeline_job_id}"
  display_name    = "test-pipeline-job-enc"
  region          = "us-central1"

  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }

  pipeline_spec = jsonencode({
    pipelineInfo = {
      name = "hello-world"
    }
    root = {
      dag = {
        tasks = {}
      }
    }
    schemaVersion = "2.1.0"
    sdkVersion    = "kfp-2.0.0"
  })

  runtime_config {
    gcs_output_directory = "gs://${google_storage_bucket.bucket.name}/pipeline_root"
  }
}
`, context)
}

func TestAccVertexAIPipelineJob_templateUri(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"pipeline_job_id": fmt.Sprintf("tf-test-tmpl-job-%s", randomSuffix),
		"storage_bucket":  fmt.Sprintf("tf-test-tmpl-bucket-%s", randomSuffix),
		"random_suffix":   randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIPipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVertexAIPipelineJob_templateUri(context),
				ExpectError: regexp.MustCompile("Failed to download package from uri|does not have permission to get ArtifactRegistry tag"),
			},
		},
	})
}

func testAccVertexAIPipelineJob_templateUri(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%{storage_bucket}"
  location                    = "us-central1"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_vertex_ai_pipeline_job" "pipeline_job" {
  pipeline_job_id = "%{pipeline_job_id}"
  display_name    = "test-template-uri"
  region          = "us-central1"
  template_uri    = "https://us-kfp.pkg.dev/proj/repo/template/v1"

  runtime_config {
    gcs_output_directory = "gs://${google_storage_bucket.bucket.name}/pipeline_root"
  }
}
`, context)
}
