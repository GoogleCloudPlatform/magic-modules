package kms_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context),
			},
			{
				ResourceName:            "google_kms_ekm_connection.example-ekmconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location		= "us-central1"
  keyManagementMode 	= MANUAL
  serviceResolvers: [
    {
      serviceDirectoryService 	= "projects/data.google_project.project.name/locations/us-central1/namespaces/google_service_directory_namespace.sd_namespace.id/services/google_service_directory_service.sd_service.id"
      hostname 			= "example.cloud.goog"
      serverCertificates: [
      	{
      		rawDer		= "chykm91dGVygoogexamplechym89"
      	}
      ]
    }
  ]
}

resource "google_service_directory_namespace" "sd_namespace" {
  namespace_id = "ekm-namespace"
  location     = "us-central1"
  project      = data.google_project.project.number
}

resource "google_service_directory_service" "sd_service" {
  service_id = "ekm-service"
  namespace  = google_service_directory_namespace.sd_namespace.id

  metadata = {
    region = "us-central1"
  }
}

data "google_project" "project" {}
`, context)
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location		= "us-central1"
  keyManagementMode 	= CLOUD_KMS
  cryptoSpacePath	= "v0/longlived/crypto-space-placeholder"
  serviceResolvers: [
    {
      serviceDirectoryService 	= "projects/data.google_project.project.name/locations/us-central1/namespaces/ekm-namespace/services/ekm-service"
      hostname 			= "example.cloud.goog"
      serverCertificates: [
      	{
      		rawDer		= "chykm91dGVygoogexamplechym89"
      	}
      ]
    }
  ]
}

data "google_project" "project" {}
`, context)
}
