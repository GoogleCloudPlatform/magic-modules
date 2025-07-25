package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
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

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "certificate-manager-trust-config-tagkey", map[string]interface{}{})

	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "certificate-manager-trust-config-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfigTags(context),
				Check: resource.TestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"google_certificate_manager_trust_config.default", "tags.%"),
				),
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

func testAccCertificateManagerTrustConfigTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
        name        = "tf-test-trust-config%{random_suffix}"
        description = "sample description for the trust config 2"
        location    = "global"
        allowlisted_certificates  {
          pem_certificate = file("test-fixtures/cert.pem") 
        }
tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
