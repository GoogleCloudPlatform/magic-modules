// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
// ----------------------------------------------------------------------------
//
//	***     HANDWRITTEN CODE     ***
//
// ----------------------------------------------------------------------------
package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccSqlCustomerManagedCertificate_postgres(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: "~> 4.0",
			},
		},
		CheckDestroy: testAccSqlCustomerManagedCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlCustomerManagedCertificate_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "bundle", "bundle1"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "intent", "PG_LOGICAL"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "primary_dns", "peer1.example.com"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "primary_ipv4_address", "10.0.0.1"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "file_paths.#", "3"),
					resource.TestCheckResourceAttrSet("google_sql_customer_managed_certificate.bundle1", "server_ca_cert_file_path"),
					resource.TestCheckResourceAttrSet("google_sql_customer_managed_certificate.bundle1", "client_cert_file_path"),
					resource.TestCheckResourceAttrSet("google_sql_customer_managed_certificate.bundle1", "client_key_file_path"),
				),
			},
			{
				ResourceName:      "google_sql_customer_managed_certificate.bundle1",
				ImportState:       true,
				ImportStateVerify: true,
				// The Cloud SQL API does not return raw PEM certificate/key content on Read, and deletion_policy is client-side only.
				ImportStateVerifyIgnore: []string{"server_ca_cert", "client_cert", "client_key", "deletion_policy"},
			},
			{
				Config: testAccSqlCustomerManagedCertificate_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "bundle", "bundle1"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "intent", "PG_LOGICAL"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "primary_dns", "peer2.example.com"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "primary_ipv4_address", "10.0.0.2"),
					resource.TestCheckResourceAttr("google_sql_customer_managed_certificate.bundle1", "file_paths.#", "3"),
				),
			},
		},
	})
}

func testAccSqlCustomerManagedCertificateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_sql_customer_managed_certificate" {
				continue
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("No ID is set")
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/customerManagedCertificates")
			if err != nil {
				return err
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
					continue
				}
				return err
			}

			bundle := rs.Primary.Attributes["bundle"]
			if rawList, ok := res["customerManagedCertificates"].([]interface{}); ok {
				for _, item := range rawList {
					if m, ok := item.(map[string]interface{}); ok {
						if b, ok := m["bundle"].(string); ok && b == bundle {
							return fmt.Errorf("CustomerManagedCertificate bundle %q still exists, should have been destroyed", bundle)
						}
					}
				}
			}
		}

		return nil
	}
}

func testAccSqlCustomerManagedCertificate_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-instance-%{random_suffix}"
  region              = "us-central1"
  database_version    = "POSTGRES_15"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "tls_private_key" "ca_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "ca_cert" {
  private_key_pem = tls_private_key.ca_key.private_key_pem

  subject {
    common_name  = "peer-ca.example.com"
    organization = "Example"
  }

  validity_period_hours = 24
  is_ca_certificate     = true

  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
  ]
}

resource "tls_private_key" "client_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "client_csr" {
  private_key_pem = tls_private_key.client_key.private_key_pem

  subject {
    common_name  = "peer-client.example.com"
    organization = "Example"
  }
}

resource "tls_locally_signed_cert" "client_cert" {
  cert_request_pem   = tls_cert_request.client_csr.cert_request_pem
  ca_private_key_pem = tls_private_key.ca_key.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca_cert.cert_pem

  validity_period_hours = 24

  allowed_uses = [
    "client_auth",
    "digital_signature",
    "key_encipherment",
  ]
}

resource "google_sql_customer_managed_certificate" "bundle1" {
  bundle               = "bundle1"
  instance             = google_sql_database_instance.instance.name
  intent               = "PG_LOGICAL"
  primary_dns          = "peer1.example.com"
  primary_ipv4_address = "10.0.0.1"
  server_ca_cert       = tls_self_signed_cert.ca_cert.cert_pem
  client_cert          = tls_locally_signed_cert.client_cert.cert_pem
  client_key           = tls_private_key.client_key.private_key_pem
}
`, context)
}

func testAccSqlCustomerManagedCertificate_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-instance-%{random_suffix}"
  region              = "us-central1"
  database_version    = "POSTGRES_15"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "tls_private_key" "ca_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "ca_cert" {
  private_key_pem = tls_private_key.ca_key.private_key_pem

  subject {
    common_name  = "peer-ca.example.com"
    organization = "Example"
  }

  validity_period_hours = 24
  is_ca_certificate     = true

  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
  ]
}

resource "tls_private_key" "client_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "client_csr" {
  private_key_pem = tls_private_key.client_key.private_key_pem

  subject {
    common_name  = "peer-client.example.com"
    organization = "Example"
  }
}

resource "tls_locally_signed_cert" "client_cert" {
  cert_request_pem   = tls_cert_request.client_csr.cert_request_pem
  ca_private_key_pem = tls_private_key.ca_key.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca_cert.cert_pem

  validity_period_hours = 24

  allowed_uses = [
    "client_auth",
    "digital_signature",
    "key_encipherment",
  ]
}

resource "google_sql_customer_managed_certificate" "bundle1" {
  bundle               = "bundle1"
  instance             = google_sql_database_instance.instance.name
  intent               = "PG_LOGICAL"
  primary_dns          = "peer2.example.com"
  primary_ipv4_address = "10.0.0.2"
  server_ca_cert       = tls_self_signed_cert.ca_cert.cert_pem
  client_cert          = tls_locally_signed_cert.client_cert.cert_pem
  client_key           = tls_private_key.client_key.private_key_pem
}
`, context)
}
