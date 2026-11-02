// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
// ----------------------------------------------------------------------------
//
//	***     HANDWRITTEN CODE     ***
//
// ----------------------------------------------------------------------------
package sql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSqlCustomerManagedCertificates_basic(t *testing.T) {
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
				Config: testAccDataSourceSqlCustomerManagedCertificates_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.#", "2"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.0.bundle", "bundle1"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.0.intent", "PG_LOGICAL"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.0.primary_dns", "peer1.example.com"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.0.primary_ipv4_address", "10.0.0.1"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.1.bundle", "bundle2"),
					resource.TestCheckResourceAttr("data.google_sql_customer_managed_certificates.all", "customer_managed_certificates.1.intent", "PG_LOGICAL"),
				),
			},
		},
	})
}

func testAccDataSourceSqlCustomerManagedCertificates_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name                = "tf-test-instance-%{random_suffix}"
  database_version    = "POSTGRES_15"
  region              = "us-central1"
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

resource "google_sql_customer_managed_certificate" "bundle1" {
  bundle               = "bundle1"
  instance             = google_sql_database_instance.main.name
  intent               = "PG_LOGICAL"
  primary_dns          = "peer1.example.com"
  primary_ipv4_address = "10.0.0.1"
  server_ca_cert       = tls_self_signed_cert.ca_cert.cert_pem
}

resource "google_sql_customer_managed_certificate" "bundle2" {
  bundle               = "bundle2"
  instance             = google_sql_database_instance.main.name
  intent               = "PG_LOGICAL"
  primary_dns          = "peer2.example.com"
  primary_ipv4_address = "10.0.0.2"
  server_ca_cert       = tls_self_signed_cert.ca_cert.cert_pem
}

data "google_sql_customer_managed_certificates" "all" {
  instance = google_sql_database_instance.main.name
  depends_on = [
    google_sql_customer_managed_certificate.bundle1,
    google_sql_customer_managed_certificate.bundle2,
  ]
}
`, context)
}
