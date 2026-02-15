package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCertificateManagerTrustConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfig_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccCertificateManagerTrustConfig_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerTrustConfig_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config"
  location = "global"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
  }

  allowlisted_certificates  {
    pem_certificate = file("test-fixtures/cert.pem") 
  }

  labels = {
    "foo" = "bar"
  }
}
`, context)
}

func testAccCertificateManagerTrustConfig_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config 2"
  location    = "global"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
  }

  allowlisted_certificates  {
    pem_certificate = file("test-fixtures/cert.pem") 
  }

  labels = {
    "bar" = "foo"
  }
}
`, context)
}

func TestAccCertificateManagerTrustConfig_tags(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfig_tags(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCertificateManagerTrustConfig_tags(context map[string]interface{}) string {
	r := fmt.Sprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config"
  location = "global"
  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
  }
  allowlisted_certificates  {
    pem_certificate = file("test-fixtures/cert.pem") 
  }
  labels = {
    "foo" = "bar"
  }
tags = {`, context)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

