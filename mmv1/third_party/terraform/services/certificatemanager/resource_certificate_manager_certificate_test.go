package certificatemanager_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCertificateManagerCertificate_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "cert-manager-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "cert-manager-tagvalue", tagKey),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_certificate_manager_certificate.test", "tags.%"),
					checkCertificateManagerCertificateTags(t),
				),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func checkCertificateManagerCertificateTags(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_certificate_manager_certificate" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// 1. Get the configured tag key and value from the state.
			var configuredTagValueNamespacedName string
			var tagKeyNamespacedName, tagValueShortName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.%" {
					tagKeyNamespacedName = strings.TrimPrefix(key, "tags.")
					tagValueShortName = val
					if tagValueShortName != "" {
						configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tagKeyNamespacedName, tagValueShortName)
						break
					}
				}
			}

			if configuredTagValueNamespacedName == "" {
				return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
			}

			// Check if placeholders are still present.
			if strings.Contains(configuredTagValueNamespacedName, "%{") {
				return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q. Ensure the context map in the test step is populated", configuredTagValueNamespacedName)
			}

			// 2. Describe the tag value using the namespaced name to get its full resource name.
			safeNamespacedName := url.QueryEscape(configuredTagValueNamespacedName)
			describeTagValueURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagValues/namespaced?name=%s", safeNamespacedName)

			respDescribe, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    describeTagValueURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error describing tag value using namespaced name %q: %v", configuredTagValueNamespacedName, err)
			}

			fullTagValueName, ok := respDescribe["name"].(string)
			if !ok || fullTagValueName == "" {
				return fmt.Errorf("tag value details (name) not found in response for namespaced name: %q, response: %v", configuredTagValueNamespacedName, respDescribe)
			}

			// 3. Get the tag bindings from the Certificate Manager Certificates.
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			instance_id := parts[5]

			parentURL := fmt.Sprintf("//certificatemanager.googleapis.com/projects/%s/locations/%s/certificates/%s", project, location, instance_id)
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    listBindingsURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error calling TagBindings API: %v", err)
			}

			tagBindingsVal, exists := resp["tagBindings"]
			if !exists {
				tagBindingsVal = []interface{}{}
			}

			tagBindings, ok := tagBindingsVal.([]interface{})
			if !ok {
				return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			// 4. Perform the comparison.
			foundMatch := false
			for _, binding := range tagBindings {
				bindingMap, ok := binding.(map[string]interface{})
				if !ok {
					continue
				}
				if bindingMap["tagValue"] == fullTagValueName {
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				return fmt.Errorf("expected tag value %s (from namespaced %q) not found in tag bindings for resource %s. Bindings: %v", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID, tagBindings)
			}

			t.Logf("Successfully found matching tag binding for %s with tagValue %s", rs.Primary.ID, fullTagValueName)
		}

		return nil
	}
}

func testAccCertificateManagerCertificateTags(context map[string]interface{}) string {
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
	}`, context)
}
