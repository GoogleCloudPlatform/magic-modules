package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceComputeSslCertificate(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeSslCertificateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeSslCertificateCheck("data.google_compute_ssl_certificate.cert", "google_compute_ssl_certificate.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceComputeSslCertificateCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		certificateAttrToCheck := []string{
			"name",
			"project",
			"description",
			"certificate",
			"certificate_id",
			"self_link",
			"creation_timestamp",
		}

		for _, attr := range certificateAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceComputeSslCertificateConfig() string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
	name		= "cert-test-%s"
	description = "really descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

data "google_compute_ssl_certificate" "cert" {
	name     = "${google_compute_ssl_certificate.foobar.name}"
}
`, acctest.RandString(10))
}
