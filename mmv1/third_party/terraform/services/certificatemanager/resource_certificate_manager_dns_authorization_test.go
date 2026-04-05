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
func TestAccCertificateManagerDnsAuthorization_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorization_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
			{
				Config: testAccCertificateManagerDnsAuthorization_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorization_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "tf-test-dns-auth%{random_suffix}"
  description = "The default dnss"
	labels = {
		a = "a"
	}
  domain      = "%{random_suffix}.hashicorptest.com"
}
`, context)
}

func testAccCertificateManagerDnsAuthorization_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "tf-test-dns-auth%{random_suffix}"
  description = "The default dnss2"
	labels = {
		a = "b"
	}
  domain      = "%{random_suffix}.hashicorptest.com"
}
`, context)
}

func TestAccCertificateManagerDnsAuthorization_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "dns-authz-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "dns-authz-tagvalue", tagKey)

	randomSuffix := acctest.RandString(t, 10)
	resourceName := "google_certificate_manager_dns_authorization.test"

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorizationWithTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkCertificateManagerDnsAuthorizationWithTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dns_authorization", "tags"},
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorizationWithTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_certificate_manager_dns_authorization" "test" {
	  name          = "tf-test-dns-auth%{random_suffix}"
        description = "The default dns"
        labels = {
                a = "a"
        }
        domain          = "%{random_suffix}.hashicorptest.com"
	tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
	}`, testContext)
}

// This function gets the DNS authorization via the Certificate Manager API and inspects its labels.
func checkCertificateManagerDnsAuthorizationWithTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		dnsAuthorizationName := rs.Primary.Attributes["name"]

		ctx := context.Background()

		certificateManagerClient, err := cmanager.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create certificate manager client: %v", err)
		}
		defer certificateManagerClient.Close()

		// Construct the request to get the DNS authorization details
		req := &certificatemanagerpb.GetDnsAuthorizationRequest{
			Name: fmt.Sprintf("projects/%s/locations/global/dnsAuthorizations/%s", project, dnsAuthorizationName),
		}

		// Get the Certificate Manager DNS authorization
		dnsAuthorization, err := certificateManagerClient.GetDnsAuthorization(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get DNS authorization '%s': %v", req.Name, err)
		}

		// Check the DNS authorization's labels for the expected tag
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		labels := dnsAuthorization.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on DNS authorization '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		return fmt.Errorf("expected tag key '%s' not found on DNS authorization '%s'", expectedTagKey, req.Name)
	}
}
