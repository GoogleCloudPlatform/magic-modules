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

func TestAccCertificateManagerTrustConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfig_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccCertificateManagerTrustConfig_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerTrustConfig_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config"
  location = "global"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
  }

  allowlisted_certificates  {
    pem_certificate = file("test-fixtures/cert.pem") 
  }

  labels = {
    "foo" = "bar"
  }
}
`, context)
}

func testAccCertificateManagerTrustConfig_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config 2"
  location    = "global"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
  }

  allowlisted_certificates  {
    pem_certificate = file("test-fixtures/cert.pem") 
  }

  labels = {
    "bar" = "foo"
  }
}
`, context)
}

func TestAccCertificateManagerTrustConfig_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "trust-config-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "trust-config-tagvalue", tagKey)

	randomSuffix := acctest.RandString(t, 10)
	resourceName := "google_certificate_manager_trust_config.test"

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfigWithTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkCertificateManagerTrustConfigWithTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"trust_store", "tags"},
			},
		},
	})
}

func testAccCertificateManagerTrustConfigWithTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_certificate_manager_trust_config" "test" {
	  name        = "tf-test-trust-config%{random_suffix}"
        description = "sample description for the trust config 2"
        location    = "global"
        allowlisted_certificates  {
          pem_certificate = file("test-fixtures/cert.pem") 
        }
tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
	}`, testContext)
}

func checkCertificateManagerTrustConfigWithTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		trustConfigName := rs.Primary.Attributes["name"]

		ctx := context.Background()

		certificateManagerClient, err := cmanager.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create certificate manager client: %v", err)
		}
		defer certificateManagerClient.Close()

		// Construct the request to get the trust config details
		req := &certificatemanagerpb.GetTrustConfigRequest{
			Name: fmt.Sprintf("projects/%s/locations/global/trustConfigs/%s", project, trustConfigName),
		}

		// Get the Certificate Manager trust config
		trustConfig, err := certificateManagerClient.GetTrustConfig(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get trust config '%s': %v", req.Name, err)
		}

		// Check the trust config's labels for the expected tag
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		labels := trustConfig.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on trust config '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		return fmt.Errorf("expected tag key '%s' not found on trust config '%s'", expectedTagKey, req.Name)
	}
}
