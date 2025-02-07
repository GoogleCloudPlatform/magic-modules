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

	context := map[string]interface{}{
		"project_id":              envvar.GetTestProjectFromEnv(),
		"project_number":          envvar.GetTestProjectNumberFromEnv(),
		"key1_name":               acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-eventarc-pipeline-key1").CryptoKey.Name,
		"key2_name":               acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-eventarc-pipeline-key2").CryptoKey.Name,
		"network_attachment_name": acctest.BootstrapNetworkAttachment(t, "tf-test-eventarc-pipeline-na", acctest.BootstrapSubnet(t, "tf-test-eventarc-pipeline-subnet", acctest.BootstrapSharedTestNetwork(t, "tf-test-eventarc-pipeline-network"))),
		"random_suffix":           acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcPipelineDestroyProducer(t),
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

resource "google_eventarc_pipeline" "primary" {
  location        = "us-central1"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  crypto_key_name = "%{key1_name}"
  destinations {
    http_endpoint {
      uri                      = "https://10.77.0.0:80/route"
      message_binding_template = "{\"headers\":{\"new-header-key\": \"new-header-value\"}}"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/%{network_attachment_name}"
    }
    output_payload_format {
      avro {
        schema_definition = "{\"type\": \"record\", \"name\": \"my_record\", \"fields\": [{\"name\": \"my_field\", \"type\": \"string\"}]}"
      }
    }
  }
  input_payload_format {
    avro {
      schema_definition = "{\"type\": \"record\", \"name\": \"my_record\", \"fields\": [{\"name\": \"my_field\", \"type\": \"string\"}]}"
    }
  }
  retry_policy {
    max_retry_delay = "50s"
    max_attempts    = 2
    min_retry_delay = "40s"
  }
  mediations {
    transformation {
      transformation_template = <<-EOF
{
"id": message.id,
"datacontenttype": "application/json",
"data": "{ \"scrubbed\": \"true\" }"
}
EOF
    }
  }
  logging_config {
    log_severity = "DEBUG"
  }
  depends_on = [google_kms_crypto_key_iam_member.key1_member]
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

resource "google_eventarc_pipeline" "primary" {
  location        = "us-central1"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  crypto_key_name = "%{key2_name}"
  destinations {
    http_endpoint {
      uri                      = "https://10.77.0.0:80/route"
      message_binding_template = "{\"headers\":{\"new-header-key\": \"new-header-value\"}}"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/%{network_attachment_name}"
    }
    output_payload_format {
      avro {
        schema_definition = "{\"type\": \"record\", \"name\": \"my_record\", \"fields\": [{\"name\": \"my_field\", \"type\": \"string\"}]}"
      }
    }
  }
  input_payload_format {
    avro {
      schema_definition = "{\"type\": \"record\", \"name\": \"my_record\", \"fields\": [{\"name\": \"my_field\", \"type\": \"string\"}]}"
    }
  }
  retry_policy {
    max_retry_delay = "50s"
    max_attempts    = 2
    min_retry_delay = "40s"
  }
  mediations {
    transformation {
      transformation_template = <<-EOF
{
"id": message.id,
"datacontenttype": "application/json",
"data": "{ \"scrubbed\": \"true\" }"
}
EOF
    }
  }
  logging_config {
    log_severity = "DEBUG"
  }
  depends_on = [google_kms_crypto_key_iam_member.key1_member, google_kms_crypto_key_iam_member.key2_member]
}
`, context)
}
