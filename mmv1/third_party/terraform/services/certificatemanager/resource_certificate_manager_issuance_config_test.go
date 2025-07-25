package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerIssuanceConfig_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "certificate_manager_issuance_config-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "certificate_manager_issuance_config-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerIssuanceConfigTags(context),
				Check: resource.TestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"google_certificate_manager_certificate_issuance_config.issuanceconfig", "tags.%"),
				),
			},
			{
				ResourceName:            "google_certificate_manager_certificate_issuance_config.issuanceconfig",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCertificateManagerIssuanceConfigTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate_issuance_config" "issuanceconfig" {
  name    = "tf-test-issuance-config%{random_suffix}"
  description = "sample description for the certificate issuanceConfigs"
  certificate_authority_config {
    certificate_authority_service_config {
        ca_pool = google_privateca_ca_pool.pool.id
    }
  }
  lifetime = "1814400s"
  rotation_window_percentage = 34
  key_algorithm = "ECDSA_P256"
  labels = { "name": "wrench", "count": "3" }

  depends_on=[google_privateca_certificate_authority.ca_authority]
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}

resource "google_privateca_ca_pool" "pool" {
  name     = "ca-pool"
  location = "us-central1"
  tier     = "ENTERPRISE"
}

resource "google_privateca_certificate_authority" "ca_authority" {
  location = "us-central1"
  pool = google_privateca_ca_pool.pool.name
  certificate_authority_id = "ca-authority"
  config {
    subject_config {
      subject {
        organization = "HashiCorp"
        common_name = "my-certificate-authority"
      }
      subject_alt_name {
        dns_names = ["hashicorp.com"]
      }
    }
    x509_config {
      ca_options {
        is_ca = true
      }
      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }

  // Disable CA deletion related safe checks for easier cleanup.
  deletion_protection                    = false
  skip_grace_period                      = true
  ignore_active_certificates_on_deletion = true
}
`, context)
}
