package certificatemanager_test

import (
	"context"
	"fmt"
	"testing"

	cmanager "cloud.google.com/go/certificatemanager/apiv1"
	certificatemanagerpb "cloud.google.com/go/certificatemanager/apiv1/certificatemanagerpb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerIssuanceConfig_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "issuance-config-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "issuance-config-tagvalue", tagKey)

	randomSuffix := acctest.RandString(t, 10)
	resourceName := "google_certificate_manager_certificate_issuance_config.test"

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerIssuanceConfigWithTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkCertificateManagerIssuanceConfigWithTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_authority_config", "tags"},
			},
		},
	})
}

func testAccCertificateManagerIssuanceConfigWithTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_certificate_manager_certificate_issuance_config" "test" {
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
  name     = "ca-pool%{random_suffix}"
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
}`, testContext)
}

// This function gets the issuance config via the Certificate Manager API and inspects its labels.
func checkCertificateManagerIssuanceConfigWithTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["location"]
		issuanceConfigName := rs.Primary.Attributes["name"]

		ctx := context.Background()

		certificateManagerClient, err := cmanager.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create certificate manager client: %v", err)
		}
		defer certificateManagerClient.Close()

		// Construct the request to get the issuance config details
		req := &certificatemanagerpb.GetCertificateIssuanceConfigRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/certificateIssuanceConfigs/%s", project, location, issuanceConfigName),
		}

		// Get the Certificate Manager issuance config
		issuanceConfig, err := certificateManagerClient.GetCertificateIssuanceConfig(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get issuance config '%s': %v", req.Name, err)
		}

		// Check the issuance config's labels for the expected tag
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		labels := issuanceConfig.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on issuance config '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		return fmt.Errorf("expected tag key '%s' not found on issuance config '%s'", expectedTagKey, req.Name)
	}
}
