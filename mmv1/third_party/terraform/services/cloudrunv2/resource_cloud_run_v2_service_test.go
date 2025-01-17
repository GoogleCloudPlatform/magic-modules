package cloudrunv2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCloudRunV2Service_cloudrunv2ServiceFunctionExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"zip_path":      "./test-fixtures/function-source.zip",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceFunctionExample_full(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "deletion_protection", "labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceFunctionExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_cloud_run_v2_service.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "deletion_protection", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceFunctionExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
  build_config {
    source_location = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    function_target = "helloHttp"
    image_uri = "us-docker.pkg.dev/cloudrun/container/hello"
    base_image = "us-central1-docker.pkg.dev/serverless-runtimes/google-22-full/runtimes/nodejs22"
    enable_automatic_updates = true
    worker_pool = "worker-pool"
    environment_variables = {
      FOO_KEY = "FOO_VALUE"
      BAR_KEY = "BAR_VALUE"
    }
    service_account = google_service_account.cloudbuild_service_account.id
  }
  depends_on = [
    google_project_iam_member.act_as,
    google_project_iam_member.logs_writer
  ]
}

data "google_project" "project" {
}

resource "google_storage_bucket" "bucket" {
  name     = "${data.google_project.project.project_id}-tf-test-gcf-source%{random_suffix}"  # Every bucket name must be globally unique
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%{zip_path}"  # Add path to the zipped function source code
}

resource "google_service_account" "cloudbuild_service_account" {
  account_id = "tf-test-build-sa%{random_suffix}"
}

resource "google_project_iam_member" "act_as" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.cloudbuild_service_account.email}"
}

resource "google_project_iam_member" "logs_writer" {
  project = data.google_project.project.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudbuild_service_account.email}"
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceFunctionExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
  build_config {
    source_location = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    function_target = "helloHttp"
    image_uri = "gcr.io/cloudrun/hello:latest"
    base_image = "us-central1-docker.pkg.dev/serverless-runtimes/google-22-full/runtimes/nodejs20"
    enable_automatic_updates = false
    worker_pool = "worker-pool-2"
    environment_variables = {
      FOO_KEY_FOO = "FOO_VALUE_FOO"
      BAR_KEY_BAR = "BAR_VALUE_BAR"
    }
    service_account = google_service_account.cloudbuild_service_account.id
  }
  depends_on = [
    google_project_iam_member.act_as,
    google_project_iam_member.logs_writer
  ]
}

data "google_project" "project" {
}

resource "google_storage_bucket" "bucket" {
  name     = "${data.google_project.project.project_id}-tf-test-gcf-source%{random_suffix}"  # Every bucket name must be globally unique
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source-updated.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%{zip_path}"  # Add path to the zipped function source code
}

resource "google_service_account" "cloudbuild_service_account" {
  account_id = "tf-test-build-sa-updated%{random_suffix}"
}

resource "google_project_iam_member" "act_as" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.cloudbuild_service_account.email}"
}

resource "google_project_iam_member" "logs_writer" {
  project = data.google_project.project.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudbuild_service_account.email}"
}
`, context)
}
