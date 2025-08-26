package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIReasoningEngine_vertexAiReasoningEngineUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"bucket_name": acctest.TestBucketName(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIReasoningEngine_vertexAiReasoningEngineBasic(context),
			},
			{
				ResourceName:            "google_vertex_ai_reasoning_engine.reasoning_engine",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIReasoningEngine_vertexAiReasoningEngineUpdate(context),
			},
			{
				ResourceName:            "google_vertex_ai_reasoning_engine.reasoning_engine",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIReasoningEngine_vertexAiReasoningEngineBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_reasoning_engine" "reasoning_engine" {
  display_name = "sample-reasoning-engine"
  description  = "A basic reasoning engine"
  region       = "us-central1"

  spec {
    agent_framework = "google-adk"
    class_methods   = []

    deployment_spec {

      # This is references inside the pickle file.
      env {
        name  = "PROJECT_ID"
        value = data.google_project.project.id
      }

      secret_env {
        name = "secret_var_1"

        secret_ref {
          secret  = google_secret_manager_secret.secret.secret_id
          version = "latest"
        }
      }
    }

    package_spec {
      dependency_files_gcs_uri = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_dependencies_tar_gz.name}"
      pickle_object_gcs_uri    = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_code_pkl.name}"
      python_version           = "3.11"
      requirements_gcs_uri     = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_requirements_txt.name}"
    }
  }

  depends_on = [
    google_secret_manager_secret_iam_member.secret_access,
    google_secret_manager_secret_version.secret_version
  ]
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "test"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "secret_access" {
  secret_id  = google_secret_manager_secret.secret.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform-re.iam.gserviceaccount.com"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%{bucket_name}"
  location                    = "us-central1"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "bucket_obj_requirements_txt" {
  name   = "requirements.txt"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/requirements.txt"
}

resource "google_storage_bucket_object" "bucket_obj_code_pkl" {
  name   = "code.pkl"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/code.pkl"
}

resource "google_storage_bucket_object" "bucket_obj_dependencies_tar_gz" {
  name   = "dependencies.tar.gz"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/dependencies.tar.gz"
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIReasoningEngine_vertexAiReasoningEngineUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_reasoning_engine" "reasoning_engine" {
  display_name = "sample-reasoning-engine-updated"
  description  = "A basic reasoning engine updated"
  region       = "us-central1"

  spec {
    agent_framework = "langchain"
    class_methods   = []

    deployment_spec {

      # This is references inside the pickle file.
      env {
        name  = "PROJECT_ID"
        value = data.google_project.project.id
      }

      env {
        name  = "REGION"
        value = "us-central1"
      }

      secret_env {
        name = "secret_var_1"

        secret_ref {
          secret  = google_secret_manager_secret.secret.secret_id
          version = "latest"
        }
      }

      secret_env {
        name = "secret_var_2"

        secret_ref {
          secret  = google_secret_manager_secret.secret_new.secret_id
          version = "2"
        }
      }
    }

    package_spec {
      dependency_files_gcs_uri = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_dependencies_tar_gz.name}"
      pickle_object_gcs_uri    = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_code_pkl.name}"
      python_version           = "3.12"
      requirements_gcs_uri     = "${google_storage_bucket.bucket.url}/${google_storage_bucket_object.bucket_obj_requirements_txt.name}"
    }
  }

  depends_on = [
    google_secret_manager_secret_iam_member.secret_access,
    google_secret_manager_secret_version.secret_version,
    google_secret_manager_secret_iam_member.secret_access_new,
    google_secret_manager_secret_version.secret_version_new_2
  ]
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "test"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version_new_1" {
  secret      = google_secret_manager_secret.secret_new.id
  secret_data = "test"
}

resource "google_secret_manager_secret_version" "secret_version_new_2" {
  secret      = google_secret_manager_secret.secret_new.id
  secret_data = "test update"

  depends_on = [
    google_secret_manager_secret_version.secret_version_new_1
  ]
}

resource "google_secret_manager_secret" "secret_new" {
  secret_id = "secret-new"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "secret_access" {
  secret_id  = google_secret_manager_secret.secret.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform-re.iam.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "secret_access_new" {
  secret_id  = google_secret_manager_secret.secret_new.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform-re.iam.gserviceaccount.com"
}

resource "google_storage_bucket" "bucket" {
  name                        = "%{bucket_name}"
  location                    = "us-central1"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "bucket_obj_requirements_txt" {
  name   = "requirements_langchain.txt"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/requirements_langchain.txt"
}

resource "google_storage_bucket_object" "bucket_obj_code_pkl" {
  name   = "code_langchain.pkl"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/code_langchain.pkl"
}

resource "google_storage_bucket_object" "bucket_obj_dependencies_tar_gz" {
  name   = "dependencies_langchain.tar.gz"
  bucket = google_storage_bucket.bucket.id
  source = "./test-fixtures/dependencies_langchain.tar.gz"
}

data "google_project" "project" {}
`, context)
}
