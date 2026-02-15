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

func TestAccCertificateManagerCertificateMap_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "cert-map-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "cert-map-tagvalue", tagKey)

	randomSuffix := acctest.RandString(t, 10)
	resourceName := "google_certificate_manager_certificate_map.test"

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateMapWithTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkCertificateManagerCertificateMapWithTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gclb_target", "tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateMapWithTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_certificate_manager_certificate_map" "test" {
	  name    = "tf-test-cert-map-%{random_suffix}"
	  description = "A basic certificate map for testing"
	  tags = {
	    "%{org}/%{tagKey}" = "%{tagValue}"
	  }
	}`, testContext)
}

// This function gets the certificate map via the Certificate Manager API and inspects its labels.
func checkCertificateManagerCertificateMapWithTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		certificateMapName := rs.Primary.Attributes["name"]

		ctx := context.Background()

		certificateManagerClient, err := cmanager.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create certificate manager client: %v", err)
		}
		defer certificateManagerClient.Close()

		// Construct the request to get the certificate map details
		req := &certificatemanagerpb.GetCertificateMapRequest{
			Name: fmt.Sprintf("projects/%s/locations/global/certificateMaps/%s", project, certificateMapName),
		}

		// Get the Certificate Manager certificate map
		certificateMap, err := certificateManagerClient.GetCertificateMap(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get certificate map '%s': %v", req.Name, err)
		}

		// Check the certificate map's labels for the expected tag
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		labels := certificateMap.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on certificate map '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		return fmt.Errorf("expected tag key '%s' not found on certificate map '%s'", expectedTagKey, req.Name)
	}
}
