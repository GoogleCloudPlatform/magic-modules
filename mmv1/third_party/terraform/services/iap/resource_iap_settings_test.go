package iap_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIapSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIapSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIapSettings_basic(context),
			},
			{
				ResourceName:      "google_iap_settings.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
						{
				Config: testAccIapSettings_update(context),
			},
			{
				ResourceName:      "google_iap_settings.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
						{
				Config: testAccIapSettings_delete(context),
			},
			{
				ResourceName:      "google_iap_settings.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}



func testAccIapSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_backend_service" "default" {
  name                            = "backend-service"
  region                          = "us-central1"
  health_checks                   = [google_compute_health_check.default.id]
  connection_draining_timeout_sec = 10
  session_affinity                = "CLIENT_IP"
}

resource "google_compute_health_check" "default" {
  name               = "rbs-health-check"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}

resource "google_iap_settings" "default" {
  name = "projects/test_project_id/iap_web/compute-us-central1/services/${google_compute_region_backend_service.default.name}"
  access_settings {
    gcip_settings {
      tenant_ids = ["tenant-id-1234"]
      login_page_uri = "test_uri"
    }
    cors_settings {
      allow_http_options = "true"
    }
    oauth_settings {
      login_hint = "test"
      programmatic_clients = ["client-ids"]
    }
    reauth_settings {
      method = "LOGIN"
      max_age = "3.5s"
      policy_type = "MINIMUM"
    }
  }
  application_settings {
    cms_settings {
      rctoken_aud = "audience"
    }
    cookie_domain = "site-cookie" 
    access_denied_page_settings {
      access_denied_page_uri = "access-denied-uri"
      generate_troubleshooting_uri = true
      remediation_token_generation_enabled = true
    }   
  }
}
`, context)
}

func testAccIapSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_iap_settings" "default" {
  name = "projects/test_project_id/iap_web/compute-us-central1/services/${google_compute_region_backend_service.default.name}"
  access_settings {
    cors_settings {
      allow_http_options = "false"
    }
    allowed_domain_settings {
      domains = ["domains-list"]
      enable = true
    }
    workforce_identity_settings {
      workforce_pools = ["wif-resource"]
      oauth2 {
        client_id = "test-client-id"
        client_Secret = "test-client-secret"
      }
    }
    identity_source = "IDENTITY_SOURCE_UNSPECIFIED"
  }
  application_settings {
    attribute_propagation_settings {
      output_credentials = "HEADER"
      enable = true
    }  
  }
}

`, context)
}

func testAccIapSettings_delete(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_iap_settings" "default" {
  name = "projects/test_project_id/iap_web/compute-us-central1/services/${google_compute_region_backend_service.default.name}"
  access_settings {
    gcip_settings {
      tenant_ids = ["tenant-id-1234"]
      login_page_uri = "test_uri"
    }
  }
}
`, context)
}

