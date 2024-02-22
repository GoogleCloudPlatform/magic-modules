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
  key_management_mode 	= "MANUAL"
  service_resolvers  	{
      service_directory_service  = "projects/315636579862/secrets/external-servicedirectoryservice/versions/latest"
      hostname 			 = "projects/315636579862/secrets/external-uri/versions/latest"
      server_certificates        {
      		raw_der	= "projects/315636579862/secrets/external-rawder/versions/latest"
      	}
    }
}
`, context)
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location     		= "us-central1"
  key_management_mode 	= "CLOUD_KMS"
  crypto_space_path	= "v0/longlived/crypto-space-placeholder"
}
`, context)
}
