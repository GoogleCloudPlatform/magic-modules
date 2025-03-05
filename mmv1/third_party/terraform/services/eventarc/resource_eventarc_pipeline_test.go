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
		"service_account":         envvar.GetTestServiceAccountFromEnv(t),
		"key_name":                acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-eventarc-pipeline-key").CryptoKey.Name,
		"key2_name":               acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-eventarc-pipeline-key2").CryptoKey.Name,
		"network_attachment_name": acctest.BootstrapNetworkAttachment(t, "tf-test-eventarc-pipeline-na", acctest.BootstrapSubnet(t, "tf-test-eventarc-pipeline-subnet", acctest.BootstrapSharedTestNetwork(t, "tf-test-eventarc-pipeline-network"))),
		"random_suffix":           acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcPipeline_eventarcPipelineWithCmekAndAvroFormatExample(context),
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
						// TODO(tommyreddad): Replace with plancheck.ResourceActionUpdate once pipeline update API is working.
						plancheck.ExpectResourceAction("google_eventarc_pipeline.primary", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "pipeline_id", "terraform_labels"},
			},
			{
				Config: testAccEventarcPipeline_unset(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// TODO(tommyreddad): Replace with plancheck.ResourceActionUpdate once pipeline update API is working.
						plancheck.ExpectResourceAction("google_eventarc_pipeline.primary", plancheck.ResourceActionDestroyBeforeCreate),
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

func testAccEventarcPipeline_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_pipeline" "primary" {
  location        = "us-central1"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  crypto_key_name = "%{key2_name}"
  display_name    = "updated pipeline"
  logging_config {
    log_severity = "ALERT"
  }
  destinations {
    http_endpoint {
      uri = "https://10.77.0.1:80/route"
      message_binding_template = "{\"headers\":{\"new-header-key2\": \"new-header-value2\"}}"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/%{network_attachment_name}"
    }
    authentication_config {
      google_oidc {
        service_account = "%{service_account}"
        audience        = "http://www.example.com"
      }
    }
    output_payload_format {
      json {}
    }
  }
  input_payload_format {
    json {}
  }
  retry_policy {
    max_retry_delay = "55s"
    max_attempts    = 3
    min_retry_delay = "45s"
  }
  mediations {
    transformation {
      transformation_template = <<-EOF
{
"id": message.id,
"datacontenttype": "application/json",
"data": "{ \"scrubbed\": \"false\" }"
}
EOF
    }
  }
}
`, context)
}

func testAccEventarcPipeline_unset(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_pipeline" "primary" {
  location        = "us-central1"
  pipeline_id     = "tf-test-some-pipeline%{random_suffix}"
  destinations {
    http_endpoint {
      uri = "https://10.77.0.1:80/route"
    }
    network_config {
      network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/%{network_attachment_name}"
    }
  }
}
`, context)
}
