package kms_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKMSEkmConnection_kmsEkmConnectionBasicExample_full(context),
			},
			{
				ResourceName:            "google_kms_ekm_connection.example-ekmconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context),
			},
			{
				ResourceName:            "google_kms_ekm_connection.example-ekmconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_secret_manager_secret_version" "raw_der" {
  secret = "playground-cert"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "hostname" {
  secret = "external-uri"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "servicedirectoryservice" {
  secret = "external-servicedirectoryservice"
  project = "315636579862"
}
data "google_project" "vpc-project" {
  project_id = "cloud-ekm-refekm-playground"
}
data "google_project" "project" {
}
resource "google_kms_ekm_connection" "example-ekmconnection" {
  count = 3
  name          = "tf_test_ekmconnection_example%{random_suffix}${count.index}"
  location		= "us-central1"
  key_management_mode 	= "MANUAL"
  service_resolvers  	{
      service_directory_service  = data.google_secret_manager_secret_version.servicedirectoryservice.secret_data
      hostname 			 = data.google_secret_manager_secret_version.hostname.secret_data
      server_certificates        {
      		raw_der	= data.google_secret_manager_secret_version.raw_der.secret_data
      }
  }
}

resource "google_kms_ekm_connection_iam_member" "add_viewer" {
  name = google_kms_ekm_connection.example-ekmconnection[0].id
  location     		= "us-central1"
  role    = "roles/viewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"

  condition {
      title       = "expires_after_2029_12_31"
      description = "Expiring at midnight of 2029-12-31"
      expression  = "request.time < timestamp(\"2030-01-01T00:00:00Z\")"
    }
  }

resource "google_kms_ekm_connection_iam_binding" "ekm_admin" {
  name = google_kms_ekm_connection.example-ekmconnection[1].id
  location     		= "us-central1"
  role = "roles/cloudkms.ekmConnectionsAdmin"
  members  = ["serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"]

  condition {
      title       = "expires_after_2029_12_31"
      description = "Expiring at midnight of 2029-12-31"
      expression  = "request.time < timestamp(\"2030-01-01T00:00:00Z\")"
    }
  }

data "google_iam_policy" "ekm_sa" {
  binding {
    role = "roles/cloudcontrolspartner.ekmServiceAgent"
    members = [
      "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com",
    ]
  }
}

resource "google_kms_ekm_connection_iam_policy" "policy" {
  name = google_kms_ekm_connection.example-ekmconnection[2].id
  policy_data = data.google_iam_policy.ekm_sa.policy_data
  location     		= "us-central1"
}
`, context)
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "vpc-project" {
  project_id = "cloud-ekm-refekm-playground"
}
data "google_project" "project" {
}
data "google_secret_manager_secret_version" "raw_der" {
  secret = "playground-cert"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "hostname" {
  secret = "external-uri"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "servicedirectoryservice" {
  secret = "external-servicedirectoryservice"
  project = "315636579862"
}
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location     		= "us-central1"
  key_management_mode 	= "CLOUD_KMS"
  crypto_space_path	= "v0/longlived/crypto-space-placeholder"
  service_resolvers  	{
      service_directory_service  = data.google_secret_manager_secret_version.servicedirectoryservice.secret_data
      hostname 			 = data.google_secret_manager_secret_version.hostname.secret_data
      server_certificates        {
      		raw_der	= data.google_secret_manager_secret_version.raw_der.secret_data
      }
  }
}
`, context)
}
