package infrastructuremanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func TestAccInfrastructureManagerDeployment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckInfrastructureManagerDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInfrastructureManagerDeployment_basic(context),
			},
			{
				ResourceName:            "google_infrastructure_manager_deployment.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "force_destroy", "labels", "annotations"},
			},
			{
				Config: testAccInfrastructureManagerDeployment_update(context),
			},
			{
				ResourceName:            "google_infrastructure_manager_deployment.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "force_destroy", "labels", "annotations"},
			},
		},
	})
}

func testAccInfrastructureManagerDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "sa" {
  account_id   = "im-basic-test-sa-%{random_suffix}"
  display_name = "Infra Manager Basic Test SA"
}

resource "google_project_iam_member" "binding" {
  project = "%{project}"
  role    = "roles/config.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "network_admin" {
  project = "%{project}"
  role    = "roles/compute.networkAdmin"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_infrastructure_manager_deployment" "basic" {
  name            = "basic-deployment-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("test-network-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}
`, context)
}

func testAccInfrastructureManagerDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "sa" {
  account_id   = "im-basic-test-sa-%{random_suffix}"
  display_name = "Infra Manager Basic Test SA"
}

resource "google_project_iam_member" "binding" {
  project = "%{project}"
  role    = "roles/config.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "network_admin" {
  project = "%{project}"
  role    = "roles/compute.networkAdmin"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_infrastructure_manager_deployment" "basic" {
  name            = "basic-deployment-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  labels = {
    env = "test"
  }

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("test-network-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}
`, context)
}

func TestAccInfrastructureManagerDeployment_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckInfrastructureManagerDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInfrastructureManagerDeployment_full(context),
			},
			{
				ResourceName:            "google_infrastructure_manager_deployment.full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "force_destroy", "labels", "annotations"},
			},
		},
	})
}

func testAccInfrastructureManagerDeployment_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "sa" {
  account_id   = "im-full-test-sa-%{random_suffix}"
  display_name = "Infra Manager Full Test SA"
}

resource "google_project_iam_member" "binding" {
  project = "%{project}"
  role    = "roles/config.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "storage_viewer" {
  project = "%{project}"
  role    = "roles/storage.objectViewer"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_storage_bucket" "blueprint_bucket" {
  name          = "im-blueprint-bucket-%{random_suffix}"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "blueprint_object" {
  name   = "blueprint.zip"
  bucket = google_storage_bucket.blueprint_bucket.name
  source = "test-fixtures/blueprint.zip"
}

resource "google_storage_bucket" "artifacts_bucket" {
  name          = "im-artifacts-bucket-%{random_suffix}"
  location      = "US"
  force_destroy = true
}

resource "google_infrastructure_manager_deployment" "full" {
  name            = "full-deployment-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true
  
  labels = {
    environment = "test"
  }
  
  annotations = {
    purpose = "full-field-testing"
  }

  terraform_blueprint {
    gcs_source = "gs://${google_storage_bucket.blueprint_bucket.name}/${google_storage_bucket_object.blueprint_object.name}"
    
    input_values {
      variable_name = "instance_name"
      input_value   = jsonencode("test-instance-%{random_suffix}")
    }
  }

  artifacts_gcs_bucket = "gs://${google_storage_bucket.artifacts_bucket.name}"

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.storage_viewer
  ]
}
`, context)
}

func testAccCheckInfrastructureManagerDeploymentDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_infrastructure_manager_deployment" {
				continue
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("Unable to verify delete of deployment ID is empty")
			}

			project, err := acctest.GetTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			parts := strings.Split(rs.Primary.ID, "/")
			deployment_id := parts[len(parts)-1]
			location := rs.Primary.Attributes["location"]

			url := fmt.Sprintf("https://config.googleapis.com/v1/projects/%s/locations/%s/deployments/%s", project, location, deployment_id)
			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   project,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					return nil
				} else if ok {
					return fmt.Errorf("Error making GCP platform call: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
				}
				return fmt.Errorf("Error making GCP platform call: %s", err.Error())
			}
			return fmt.Errorf("Deployment still exists")
		}

		return nil
	}
}
