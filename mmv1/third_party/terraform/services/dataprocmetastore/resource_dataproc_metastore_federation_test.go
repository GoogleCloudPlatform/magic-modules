// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataprocmetastore_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"regexp"
	"testing"

	metastore "cloud.google.com/go/metastore/apiv1"
	metastorepb "cloud.google.com/go/metastore/apiv1/metastorepb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetastoreFederation_deletionprotection(t *testing.T) {
	t.Parallel()

	name := "tf-test-metastore-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreFederationDeletionProtection(name, "us-central1"),
			},
			{
				ResourceName:            "google_dataproc_metastore_federation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:      testAccMetastoreFederationDeletionProtection(name, "us-west2"),
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccMetastoreFederationDeletionProtectionFalse(name, "us-central1"),
			},
			{
				ResourceName:            "google_dataproc_metastore_federation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccMetastoreFederation_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "metastore-federation-tags-key", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "metastore-federation-tags-value", tagKey)

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "google_dataproc_metastore_federation.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreFederationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreFederationTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkMetastoreFederationTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
			{
				Config: testAccMetastoreFederationTagsDeletionProtection(context),
			},
		},
	})
}

func testAccMetastoreFederationDeletionProtection(name string, location string) string {

	return fmt.Sprintf(`
       resource "google_dataproc_metastore_service" "default" {
         service_id = "%s"
         location   = "us-central1"
         tier       = "DEVELOPER"
         hive_metastore_config {
           version           = "3.1.2"
           endpoint_protocol = "GRPC"
         }
         }
       resource "google_dataproc_metastore_federation" "default" {
          federation_id = "%s"
          location      = "%s"
          version       = "3.1.2"
          deletion_protection = true
          backend_metastores {
            rank           = "1"
            name           = google_dataproc_metastore_service.default.id
            metastore_type = "DATAPROC_METASTORE" 
         }
}
`, name, name, location)
}

func testAccMetastoreFederationDeletionProtectionFalse(name string, location string) string {

	return fmt.Sprintf(`
       resource "google_dataproc_metastore_service" "default" {
         service_id = "%s"
         location   = "us-central1"
         tier       = "DEVELOPER"
         hive_metastore_config {
           version           = "3.1.2"
           endpoint_protocol = "GRPC"
         }
         }
       resource "google_dataproc_metastore_federation" "default" {
          federation_id = "%s"
          location      = "%s"
          version       = "3.1.2"
          deletion_protection = false
          backend_metastores {
            rank           = "1"
            name           = google_dataproc_metastore_service.default.id
            metastore_type = "DATAPROC_METASTORE" 
         }
}
`, name, name, location)
}

func testAccMetastoreFederationTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dataproc_metastore_service" "backend" {
	  service_id = "tf-test-service-%{random_suffix}"
	  location   = "us-east1"
	  tier       = "DEVELOPER"
	  hive_metastore_config {
				version           = "3.1.2"
				endpoint_protocol = "GRPC"
			}
	}

	resource "google_dataproc_metastore_federation" "test" {
	  federation_id = "tf-test-federation-%{random_suffix}"
	  location      = "us-east1"
	  version       = "3.1.2"
	  backend_metastores {
	    name = google_dataproc_metastore_service.backend.id
	    metastore_type = "DATAPROC_METASTORE"
	    rank = "1"
	  }
	  tags = {
	    "%{org}/%{tagKey}" = "%{tagValue}"
	  }
	}`, testContext)
}

// This function gets the federation via the Metastore API and inspects its tags.
func checkMetastoreFederationTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["location"]
		federationName := rs.Primary.Attributes["federation_id"]

		// Construct the expected full tag key
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		// This `ctx` variable is now a `context.Context` object
		ctx := context.Background()

		// Create a Metastore client for Federation
		metastoreClient, err := metastore.NewDataprocMetastoreFederationClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create metastore federation client: %v", err)
		}
		defer metastoreClient.Close()

		// Construct the request to get the federation details
		req := &metastorepb.GetFederationRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/federations/%s", project, location, federationName),
		}

		// Get the Metastore federation
		federation, err := metastoreClient.GetFederation(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get metastore federation '%s': %v", req.Name, err)
		}

		labels := federation.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on federation '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		// If we reach here, the tag key was not found.
		return fmt.Errorf("expected tag key '%s' not found on federation '%s'", expectedTagKey, req.Name)
	}
}

func testAccMetastoreFederationTagsDeletionProtection(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_dataproc_metastore_service" "default" {
			service_id = "tf-test-service-%{random_suffix}"
			location   = "us-central1"
			tier       = "DEVELOPER"
			hive_metastore_config {
				version           = "3.1.2"
				endpoint_protocol = "GRPC"
			}
		}
		resource "google_dataproc_metastore_federation" "default" {
			location      = "us-central1"
			federation_id = "tf-test-federation-%{random_suffix}"
			version       = "3.1.2"
			backend_metastores {
				rank           = "1"
				name           = google_dataproc_metastore_service.default.id
				metastore_type = "DATAPROC_METASTORE"
			}
			tags = {"%{org}/%{tagKey}" = "%{tagValue}"}
			deletion_protection = false 
		}
	`, context)

}
