package vertexai_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func bootstrapTestKmsKey(t *testing.T) string {
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}
	project := envvar.GetTestProjectFromEnv()
	keyName := fmt.Sprintf("projects/%s/locations/us-central1/keyRings/tftest-shared-keyring-1/cryptoKeys/tf-bootstrap-re-key1", project)
	kmsClient := kms.NewClient(config, config.UserAgent)
	if _, err := kmsClient.Projects.Locations.KeyRings.CryptoKeys.Get(keyName).Do(); err == nil {
		return keyName
	}
	return kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-re-key1").CryptoKey.Name
}

func bootstrapTestParentModel(t *testing.T) string {
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}
	project := envvar.GetTestProjectFromEnv()
	location := "us-central1"

	baseUrl := transport_tpg.BaseUrl(vertexai.Product, config)

	// Check if a parent model already exists
	listUrl := fmt.Sprintf("%sprojects/%s/locations/%s/models", baseUrl, project, location)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    listUrl,
		UserAgent: config.UserAgent,
	})
	if err == nil {
		if models, ok := res["models"].([]interface{}); ok {
			for _, m := range models {
				if mObj, ok := m.(map[string]interface{}); ok {
					if name, ok := mObj["name"].(string); ok {
						if dName, ok := mObj["displayName"].(string); ok && strings.Contains(dName, "parent") {
							return name
						}
					}
				}
			}
		}
	}

	uploadUrl := fmt.Sprintf("%sprojects/%s/locations/%s/models:upload", baseUrl, project, location)
	uploadBody := map[string]interface{}{
		"model": map[string]interface{}{
			"displayName": "tf-test-parent-model",
			"artifactUri": "gs://cloud-samples-data/vertex-ai/prediction/test-models-requests/tensorflow",
			"containerSpec": map[string]interface{}{
				"imageUri": "us-docker.pkg.dev/vertex-ai/prediction/tf2-cpu.2-8:latest",
			},
		},
	}
	opRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    uploadUrl,
		UserAgent: config.UserAgent,
		Body:      uploadBody,
	})
	if err != nil {
		return fmt.Sprintf("projects/%s/locations/%s/models/540124641784496128", project, location)
	}

	var opResponse map[string]interface{}
	err = vertexai.VertexAIOperationWaitTimeWithResponse(config, opRes, &opResponse, project, location, "Uploading parent model", config.UserAgent, 5*time.Minute)
	if err != nil {
		return fmt.Sprintf("projects/%s/locations/%s/models/540124641784496128", project, location)
	}
	if modelName, ok := opResponse["model"].(string); ok && modelName != "" {
		return modelName
	}

	return fmt.Sprintf("projects/%s/locations/%s/models/540124641784496128", project, location)
}

func TestAccVertexAITrainingPipeline_labelsForceNew(t *testing.T) {
	t.Parallel()

	pipelineName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAITrainingPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAITrainingPipeline_labels(pipelineName, "initial"),
			},
			{
				ResourceName:            "google_vertex_ai_training_pipeline.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "region", "start_time", "state", "terraform_labels", "training_task_metadata", "update_time"},
			},
			{
				Config: testAccVertexAITrainingPipeline_labels(pipelineName, "updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vertex_ai_training_pipeline.pipeline", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				ResourceName:            "google_vertex_ai_training_pipeline.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "region", "start_time", "state", "terraform_labels", "training_task_metadata", "update_time"},
			},
		},
	})
}

func testAccVertexAITrainingPipeline_labels(displayName, labelValue string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  labels = {
    test_key = "%s"
  }
}
`, displayName, labelValue)
}

func TestAccVertexAITrainingPipeline_inputDataConfigSplits(t *testing.T) {
	t.Parallel()

	pipelineName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	datasetName := fmt.Sprintf("tf-test-ds-%s", acctest.RandString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAITrainingPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAITrainingPipeline_predefinedSplit(pipelineName, datasetName, bucketName),
			},
			{
				ResourceName:            "google_vertex_ai_training_pipeline.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "region", "start_time", "state", "terraform_labels", "training_task_metadata", "update_time"},
			},
			{
				Config: testAccVertexAITrainingPipeline_timestampSplit(pipelineName, datasetName, bucketName),
			},
			{
				ResourceName:            "google_vertex_ai_training_pipeline.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "region", "start_time", "state", "terraform_labels", "training_task_metadata", "update_time"},
			},
			{
				Config: testAccVertexAITrainingPipeline_stratifiedSplit(pipelineName, datasetName, bucketName),
			},
			{
				ResourceName:            "google_vertex_ai_training_pipeline.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "region", "start_time", "state", "terraform_labels", "training_task_metadata", "update_time"},
			},
			{
				Config:      testAccVertexAITrainingPipeline_bigqueryDestination(pipelineName, datasetName),
				ExpectError: regexp.MustCompile("Please provide TrainingPipeline.input_data_config.gcs_destination"),
			},
			{
				Config:      testAccVertexAITrainingPipeline_savedQueryId(pipelineName, datasetName, bucketName),
				ExpectError: regexp.MustCompile("SavedQuery .* not found"),
			},
		},
	})
}

func testAccVertexAITrainingPipeline_predefinedSplit(pipelineName, datasetName, bucketName string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  input_data_config {
    dataset_id = google_vertex_ai_dataset.dataset.name

    predefined_split {
      key = "split"
    }

    gcs_destination {
      output_uri_prefix = "gs://${google_storage_bucket.bucket.name}/dataset-export"
    }
  }

  labels = {
    env = "test"
  }
}

resource "google_vertex_ai_dataset" "dataset" {
  display_name        = "%s"
  metadata_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/metadata/tabular_1.0.0.yaml"
  region              = "us-central1"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}
`, pipelineName, datasetName, bucketName)
}

func testAccVertexAITrainingPipeline_timestampSplit(pipelineName, datasetName, bucketName string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  input_data_config {
    dataset_id = google_vertex_ai_dataset.dataset.name

    timestamp_split {
      key                 = "ts"
      training_fraction   = 0.8
      validation_fraction = 0.1
      test_fraction       = 0.1
    }

    gcs_destination {
      output_uri_prefix = "gs://${google_storage_bucket.bucket.name}/dataset-export"
    }
  }

  labels = {
    env = "test"
  }
}

resource "google_vertex_ai_dataset" "dataset" {
  display_name        = "%s"
  metadata_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/metadata/tabular_1.0.0.yaml"
  region              = "us-central1"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}
`, pipelineName, datasetName, bucketName)
}

func testAccVertexAITrainingPipeline_stratifiedSplit(pipelineName, datasetName, bucketName string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  input_data_config {
    dataset_id = google_vertex_ai_dataset.dataset.name

    stratified_split {
      key                 = "label"
      training_fraction   = 0.8
      validation_fraction = 0.1
      test_fraction       = 0.1
    }

    gcs_destination {
      output_uri_prefix = "gs://${google_storage_bucket.bucket.name}/dataset-export"
    }
  }

  labels = {
    env = "test"
  }
}

resource "google_vertex_ai_dataset" "dataset" {
  display_name        = "%s"
  metadata_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/metadata/tabular_1.0.0.yaml"
  region              = "us-central1"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}
`, pipelineName, datasetName, bucketName)
}

func testAccVertexAITrainingPipeline_bigqueryDestination(pipelineName, datasetName string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  input_data_config {
    dataset_id = google_vertex_ai_dataset.dataset.name

    timestamp_split {
      key                 = "ts"
      training_fraction   = 0.8
      validation_fraction = 0.1
      test_fraction       = 0.1
    }

    bigquery_destination {
      output_uri = "bq://graphite-test-sam-again"
    }
  }

  labels = {
    env = "test"
  }
}

resource "google_vertex_ai_dataset" "dataset" {
  display_name        = "%s"
  metadata_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/metadata/tabular_1.0.0.yaml"
  region              = "us-central1"
}
`, pipelineName, datasetName)
}

func testAccVertexAITrainingPipeline_savedQueryId(pipelineName, datasetName, bucketName string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_training_pipeline" "pipeline" {
  display_name             = "%s"
  region                   = "us-central1"
  training_task_definition = "gs://google-cloud-aiplatform/schema/trainingjob/definition/custom_task_1.0.0.yaml"

  training_task_inputs = jsonencode({
    workerPoolSpecs = [
      {
        machineSpec = {
          machineType = "n1-standard-4"
        }
        replicaCount = "1"
        diskSpec = {
          bootDiskType   = "pd-ssd"
          bootDiskSizeGb = 100
        }
        containerSpec = {
          imageUri = "gcr.io/cloud-marketplace/google/debian11"
        }
      }
    ]
  })

  input_data_config {
    dataset_id            = google_vertex_ai_dataset.dataset.name
    annotation_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/annotation/image_classification_1.0.0.yaml"
    saved_query_id        = "123456789"

    gcs_destination {
      output_uri_prefix = "gs://${google_storage_bucket.bucket.name}/dataset-export"
    }
  }

  labels = {
    env = "test"
  }
}

resource "google_vertex_ai_dataset" "dataset" {
  display_name        = "%s"
  metadata_schema_uri = "gs://google-cloud-aiplatform/schema/dataset/metadata/image_1.0.0.yaml"
  region              = "us-central1"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}
`, pipelineName, datasetName, bucketName)
}
