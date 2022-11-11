package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVertexAITensorboard_Full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVertexAITensorboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAITensorboard_Full(context),
			},
		},
	})
}

func testAccVertexAITensorboard_Full(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "project" {
}

resource "google_kms_key_ring" "keyring" {
  name     = "keyring-%{random_suffix}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "example-key" {
  name            = "crypto-key-%{random_suffix}"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"  
  lifecycle {
    prevent_destroy = false
  }
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_encrypt" {
  crypto_key_id = google_kms_crypto_key.example-key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"  
  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_decrypt" {
  crypto_key_id = google_kms_crypto_key.example-key.id
  role          = "roles/cloudkms.cryptoKeyDecrypter"
  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
  ]
}

resource "google_vertex_ai_tensorboard" "tensorboard" {
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key_encrypt, google_kms_crypto_key_iam_binding.crypto_key_decrypt]
  display_name = "terraform%{random_suffix}"
  description  = "sample description"
  labels       = {
    "key1" : "value1",
    "key2" : "value2"
  }
  region       = "us-central1"
  encryption_spec {
	kms_key_name = google_kms_crypto_key.example-key.id
  }
}
`, context)
}
