package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccEventarcPipeline_update(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"region":                  region,
		"project_id":              envvar.GetTestProjectFromEnv(),
		"project_number":          envvar.GetTestProjectNumberFromEnv(),
		"key1_name":               acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-pipeline-key1").CryptoKey.Name,
		"key2_name":               acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-pipeline-key2").CryptoKey.Name,
		"network_attachment_name": acctest.BootstrapNetworkAttachment(t, "tf-test-eventarc-pipeline-na", acctest.BootstrapSubnet(t, "tf-test-eventarc-pipeline-subnet", acctest.BootstrapSharedTestNetwork(t, "tf-test-eventarc-pipeline-network"))),
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcPipelineDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcPipeline_full(context),
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "pipeline_id", "terraform_labels"},
			},
			{
				Config: testAccEventarcPipeline_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_pipeline.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "pipeline_id", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcPipeline_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key1_member" {
  crypto_key_id = "%{key1_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "key2_member" {
  crypto_key_id = "%{key2_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_pipeline" "primary" {
  location        = "%{region}"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  crypto_key_name = "%{key1_name}"
  destinations {
    http_endpoint {
      uri = "https://10.77.0.0:80/route"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/%{region}/networkAttachments/%{network_attachment_name}"
    }
  }
  depends_on = [google_kms_crypto_key_iam_member.key1_member, google_kms_crypto_key_iam_member.key2_member]
}
`, context)
}

func testAccEventarcPipeline_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key1_member" {
  crypto_key_id = "%{key1_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "key2_member" {
  crypto_key_id = "%{key2_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "time_sleep" "pre_wait_600s" {
  create_duration = "600s"
}

resource "google_eventarc_pipeline" "primary" {
  location        = "%{region}"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  crypto_key_name = "%{key2_name}"
  destinations {
    http_endpoint {
      uri = "https://10.77.0.0:80/route"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/%{region}/networkAttachments/%{network_attachment_name}"
    }
  }
  depends_on = [time_sleep.pre_wait_600s, google_kms_crypto_key_iam_member.key1_member, google_kms_crypto_key_iam_member.key2_member]
}

resource "time_sleep" "post_wait_600s" {
  create_duration = "600s"
  depends_on      = [google_eventarc_pipeline.primary]
}
`, context)
}
