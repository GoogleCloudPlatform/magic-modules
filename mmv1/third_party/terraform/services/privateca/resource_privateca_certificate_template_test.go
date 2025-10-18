package privateca_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
)

func TestAccPrivatecaCertificateTemplate_BasicCertificateTemplate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplate(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplateUpdate0(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPrivatecaCertificateTemplate_BasicCertificateTemplateLongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplateLongForm(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplateLongFormUpdate0(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
		},
	})
}

func TestAccPrivatecaCertificateTemplate_updateCaOption(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsTrueAndMaxPathIsPositive(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsFalse(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsNull(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionMaxIssuerPathLenghIsZero(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage", "labels", "terraform_labels", "project", "location", "name"},
			},
		},
	})
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "86400s"
  description      = "An updated sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true

    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }

    known_extensions = ["EXTENDED_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }

      value    = "c3RyaW5nCg=="
      critical = true
    }

    aia_ocsp_servers = ["string"]

    ca_options {
      is_ca                  = false
      max_issuer_path_length = 6
    }

    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }

      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }

      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }

    policy_ids {
      object_id_path = [1, 6]
    }
  }

  project = "%{project_name}"

  labels = {
    label-two = "value-two"
  }
}


`, context)
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplateUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "172800s"
  description      = "A sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = false
    allow_subject_passthrough           = false

    cel_expression {
      description = "Always false"
      expression  = "false"
      location    = "update.certificate_template.json"
      title       = "New sample expression"
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 7]
    }

    known_extensions = ["BASE_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 7]
      }

      value    = "bmV3LXN0cmluZw=="
      critical = false
    }

    aia_ocsp_servers = ["new-string"]

    ca_options {
      is_ca                  = true
      max_issuer_path_length = 7
    }

    key_usage {
      base_key_usage {
        cert_sign          = true
        content_commitment = false
        crl_sign           = true
        data_encipherment  = false
        decipher_only      = false
        digital_signature  = false
        encipher_only      = false
        key_agreement      = false
        key_encipherment   = false
      }

      extended_key_usage {
        client_auth      = false
        code_signing     = false
        email_protection = false
        ocsp_signing     = false
        server_auth      = false
        time_stamping    = false
      }

      unknown_extended_key_usages {
        object_id_path = [1, 7]
      }
    }

    policy_ids {
      object_id_path = [1, 7]
    }
  }

  project = "%{project_name}"

  labels = {
    label-one = "value-one"
  }
}


`, context)
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplateLongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location    = "long/form/%{region}"
  name        = "long/form/tf-test-template%{random_suffix}"
  description = "An updated sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true

    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }

    known_extensions = ["EXTENDED_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }

      value    = "c3RyaW5nCg=="
      critical = true
    }

    aia_ocsp_servers = ["string"]

    ca_options {
      is_ca                  = false
      max_issuer_path_length = 6
    }

    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }

      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }

      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }

    policy_ids {
      object_id_path = [1, 6]
    }
  }

  project = "projects/%{project_name}"

  labels = {
    label-two = "value-two"
  }
}


`, context)
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplateLongFormUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location    = "long/form/%{region}"
  name        = "long/form/tf-test-template%{random_suffix}"
  description = "A sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = false
    allow_subject_passthrough           = false

    cel_expression {
      description = "Always false"
      expression  = "false"
      location    = "update.certificate_template.json"
      title       = "New sample expression"
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 7]
    }

    known_extensions = ["BASE_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 7]
      }

      value    = "bmV3LXN0cmluZw=="
      critical = false
    }

    aia_ocsp_servers = ["new-string"]

    ca_options {
      is_ca                  = true
      max_issuer_path_length = 7
    }

    key_usage {
      base_key_usage {
        cert_sign          = true
        content_commitment = false
        crl_sign           = true
        data_encipherment  = false
        decipher_only      = false
        digital_signature  = false
        encipher_only      = false
        key_agreement      = false
        key_encipherment   = false
      }

      extended_key_usage {
        client_auth      = false
        code_signing     = false
        email_protection = false
        ocsp_signing     = false
        server_auth      = false
        time_stamping    = false
      }

      unknown_extended_key_usages {
        object_id_path = [1, 7]
      }
    }

    policy_ids {
      object_id_path = [1, 7]
    }
  }

  project = "projects/%{project_name}"

  labels = {
    label-one = "value-one"
  }
}


`, context)
}

func testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsTrueAndMaxPathIsPositive(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "86400s"
  description      = "A sample certificate template"
  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true
    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }
  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }
    known_extensions = ["EXTENDED_KEY_USAGE"]
  }
  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }
      value    = "c3RyaW5nCg=="
      critical = true
    }
    aia_ocsp_servers = ["string"]
    ca_options {
      is_ca                  = true
      max_issuer_path_length = 6
    }
    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }
      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }
      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }
    policy_ids {
      object_id_path = [1, 6]
    }
  }
  project = "%{project_name}"
  labels = {
    label-two = "value-two"
  }
}
`, context)
}

func testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsFalse(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "86400s"
  description      = "An updated sample certificate template"
  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true
    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }
  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }
    known_extensions = ["EXTENDED_KEY_USAGE"]
  }
  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }
      value    = "c3RyaW5nCg=="
      critical = true
    }
    aia_ocsp_servers = ["string"]
    ca_options {
      is_ca = false
    }
    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }
      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }
      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }
    policy_ids {
      object_id_path = [1, 6]
    }
  }
  project = "%{project_name}"
  labels = {
    label-two = "value-two"
  }
}
`, context)
}

func testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionIsCaIsNull(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "86400s"
  description      = "An updated sample certificate template"
  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true
    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }
  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }
    known_extensions = ["EXTENDED_KEY_USAGE"]
  }
  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }
      value    = "c3RyaW5nCg=="
      critical = true
    }
    aia_ocsp_servers = ["string"]
    ca_options {
      null_ca = true
      is_ca = false
    }
    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }
      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }
      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }
    policy_ids {
      object_id_path = [1, 6]
    }
  }
  project = "%{project_name}"
  labels = {
    label-two = "value-two"
  }
}
`, context)
}

func testAccPrivatecaCertificateTemplate_CertificateTemplateCaOptionMaxIssuerPathLenghIsZero(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location         = "%{region}"
  name             = "tf-test-template%{random_suffix}"
  maximum_lifetime = "86400s"
  description      = "Another updated sample certificate template"
  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true
    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }
  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }
    known_extensions = ["EXTENDED_KEY_USAGE"]
  }
  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }
      value    = "c3RyaW5nCg=="
      critical = true
    }
    aia_ocsp_servers = ["string"]
    ca_options {
      zero_max_issuer_path_length = true
      max_issuer_path_length = 0
    }
    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }
      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }
      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }
    policy_ids {
      object_id_path = [1, 6]
    }
  }
  project = "%{project_name}"
  labels = {
    label-two = "value-two"
  }
}
`, context)
}

func TestAccPrivatecaCertificateTemplate_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "pca-certtmpl-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "pca-certtmpl-tagvalue", tagKey)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateTemplateTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_privateca_certificate_template.default", "tags.%"),
					testAccCheckPrivatecaCertificateTemplateHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_privateca_certificate_template.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCheckPrivatecaCertificateTemplateHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_privateca_certificate_template" {
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

			// 3. Get the tag bindings from the Private CA Certificate Template.
			// Example ID: projects/my-project/locations/us-central1/certificateTemplates/my-template
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format for CertificateTemplate: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			templateId := parts[5]

			// The parent resource name for TagBindings API for a Certificate Template
			parentURL := fmt.Sprintf("//privateca.googleapis.com/projects/%s/locations/%s/certificateTemplates/%s", project, location, templateId)
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

func testAccPrivatecaCertificateTemplateTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "default" {
  name        = "tf-test-my-tmpl-%{random_suffix}"
  location   = "us-east1"
  description = "An example certificate template"

  identity_constraints {
    allow_subject_passthrough = false
    allow_subject_alt_names_passthrough = false
    cel_expression {
      expression = "subject.organization == \"Google LLC\""
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 2, 3]
    }
  }

  predefined_values {
    aia_ocsp_servers = ["http://example.com"]
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
