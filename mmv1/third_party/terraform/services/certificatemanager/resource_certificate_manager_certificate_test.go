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

func TestAccCertificateManagerCertificate_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "cert-manager-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "cert-manager-tagvalue", tagKey)

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "google_certificate_manager_certificate.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkCertificateManagerCertificateTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_certificate_manager_certificate" "test" {
	  name    = "tf-test-cert-%{random_suffix}"
	  description = "Global cert"
	  location = "us-east1"
	  self_managed {
	    pem_certificate = file("test-fixtures/cert.pem")
    	    pem_private_key = file("test-fixtures/private-key.pem")
	  }
	  tags = {
	    "%{org}/%{tagKey}" = "%{tagValue}"
	  }
	}`, testContext)
}

// This function gets the certificate via the Certificate Manager API and inspects its tags.
func checkCertificateManagerCertificateTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["location"]
		certificateName := rs.Primary.Attributes["name"]

		// Construct the expected full tag key
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		// This `ctx` variable is now a `context.Context` object
		ctx := context.Background()

		// Create a Certificate Manager client

		certificateManagerClient, err := cmanager.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create certificate manager client: %v", err)
		}
		defer certificateManagerClient.Close()
		// Construct the request to get the certificate details
		req := &certificatemanagerpb.GetCertificateRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/certificates/%s", project, location, certificateName),
		}

		// Get the Certificate Manager certificate
		certificate, err := certificateManagerClient.GetCertificate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get certificate '%s': %v", req.Name, err)
		}

		// Check the certificate's labels for the expected tag
		// In the Certificate Manager API, tags are represented as labels.
		labels := certificate.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on certificate '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		// If we reach here, the tag key was not found.
		return fmt.Errorf("expected tag key '%s' not found on certificate '%s'", expectedTagKey, req.Name)
	}
}
