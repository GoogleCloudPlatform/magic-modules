package privateca_test

import (	
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
)

func TestAccPrivatecaCaPool_privatecaCapoolUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolStart(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEnd(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolStart(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = true
  }
  labels = {
    foo = "bar"
  }
  issuance_policy {
    allowed_key_types {
      elliptic_curve {
        signature_algorithm = "ECDSA_P256"
      }
    }
    allowed_key_types {
      rsa {
        min_modulus_size = 5
        max_modulus_size = 10
      }
    }
    maximum_lifetime = "50000s"
    allowed_issuance_modes {
      allow_csr_based_issuance = true
      allow_config_based_issuance = false
    }
    identity_constraints {
      allow_subject_passthrough = false
      allow_subject_alt_names_passthrough = true
      cel_expression {
        expression = "subject_alt_names.all(san, san.type == DNS || san.type == EMAIL )"
        title = "My title"
      }
    }
    baseline_values {
      aia_ocsp_servers = ["example.com"]
      additional_extensions {
        critical = true
        value = "asdf"
        object_id {
          object_id_path = [1, 5]
        }
      }
      policy_ids {
        object_id_path = [1, 7]
      }
      policy_ids {
        object_id_path = [1,5,7]
      }
      ca_options {
        is_ca = true
        max_issuer_path_length = 10
      }
      key_usage {
        base_key_usage {
          digital_signature = true
          content_commitment = true
          key_encipherment = false
          data_encipherment = true
          key_agreement = true
          cert_sign = false
          crl_sign = true
          decipher_only = true
        }
        extended_key_usage {
          server_auth = true
          client_auth = false
          email_protection = true
          code_signing = true
          time_stamping = true
        }
      }
    }
  }
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolEnd(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = true
    publish_crl = true
  }
  labels = {
    foo = "bar"
    baz = "qux"
  }
  issuance_policy {
    allowed_key_types {
      elliptic_curve {
        signature_algorithm = "ECDSA_P256"
      }
    }
    allowed_key_types {
      rsa {
        min_modulus_size = 6
      }
    }
    maximum_lifetime = "3000s"
    allowed_issuance_modes {
      allow_csr_based_issuance = true
      allow_config_based_issuance = true
    }
    identity_constraints {
      allow_subject_passthrough = true
      allow_subject_alt_names_passthrough = true
      cel_expression {
        expression = "subject_alt_names.all(san, san.type == DNS || san.type == EMAIL )"
        title = "My title3"
      }
    }
    baseline_values {
      aia_ocsp_servers = ["example.com", "hashicorp.com"]
      additional_extensions {
        critical = true
        value = "asdf"
        object_id {
          object_id_path = [1, 7]
        }
      }
      policy_ids {
        object_id_path = [1, 5]
      }
      policy_ids {
        object_id_path = [1, 7]
      }
      ca_options {
        is_ca = true
        max_issuer_path_length = 10
      }
      key_usage {
        base_key_usage {
          digital_signature = true
          content_commitment = true
          key_encipherment = false
          data_encipherment = true
          key_agreement = false
          cert_sign = false
          crl_sign = true
          decipher_only = false
        }
        extended_key_usage {
          server_auth = false
          client_auth = true
          email_protection = true
          code_signing = true
          time_stamping = false
        }
      }
    }
  }
}
`, context)
}

func TestAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = true
  }
  labels = {
    foo = "bar"
  }
  issuance_policy {
    baseline_values {
      additional_extensions {
        critical = false
        value = "asdf"
        object_id {
          object_id_path = [1, 6]
        }
      }
      ca_options {
        is_ca = false
      }
      key_usage {
        base_key_usage {
          digital_signature = false
        }
        extended_key_usage {
          server_auth = false
        }
      }
    }
  }
}
`, context)
}

func TestAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = false
  }
  labels = {
    foo = "bar"
  }
}
`, context)
}

func TestAccPrivatecaCaPool_updateCaOption(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsTrueAndMaxPathIsPositive(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsFalse(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionMaxIssuerPathLenghIsZero(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsTrueAndMaxPathIsPositive(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        is_ca = true
        max_issuer_path_length = 10
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsFalse(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        non_ca = true
        is_ca = false
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionMaxIssuerPathLenghIsZero(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        zero_max_issuer_path_length = true
        max_issuer_path_length = 0
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}


func TestAccPrivatecaCaPool_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "pca-capool-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "pca-capool-tagvalue", tagKey)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPoolTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_privateca_ca_pool.default", "tags.%"),
					testAccCheckPrivatecaCaPoolHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCheckPrivatecaCaPoolHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_privateca_ca_pool" {
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
				if strings.HasPrefix(key, "tags.") && key != "tags.#" {
					tagKeyNamespacedName = strings.TrimPrefix(key, "tags.")
					tagValueShortName = val
					if tagValueShortName != "" {
						// Assumes tag key is in format {org_id}/{key_short_name} or {project_id}/{key_short_name}
						// The tagKeyNamespacedName from state aelready seems to have this.
						configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tagKeyNamespacedName, tagValueShortName)
						break
					}
				}
			}

			if configuredTagValueNamespacedName == "" {
				return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
			}

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

			// 3. Get the tag bindings from the Private CA CaPool.
			// Example ID: projects/my-project/locations/us-central1/caPools/my-pool
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format for CaPool: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			caPoolId := parts[5]

			// The parent resource name for TagBindings API for a CaPool
			parentURL := fmt.Sprintf("//privateca.googleapis.com/projects/%s/locations/%s/caPools/%s", project, location, caPoolId)
			// TagBindings API is regional.
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    listBindingsURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error calling TagBindings API for %s: %v", parentURL, err)
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

func testAccPrivatecaCaPoolTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name     = "tf-test-my-capool-%{random_suffix}"
  location   = "us-east1"
  tier     = "ENTERPRISE"

  publishing_options {
    publish_ca_cert = true
    publish_crl     = true
  }

  labels = {
    env = "test"
  }

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

