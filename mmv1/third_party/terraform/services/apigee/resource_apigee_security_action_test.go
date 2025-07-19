package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccApigeeSecurityAction_apigeeSecurityActionFull(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeSecurityActionDestroyProducer(t),
		/* allow, deny and flag are mutually exclusive, so we test them in sequence */
		/* also all conditions except ip_address_ranges and bot_reasons seem to be mutually exclusive, so we test them in sequence */
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullAllow(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullDeny(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullHttpMethods(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullFlag(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullApiKeys(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullAccessTokens(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullApiProducts(context),
			},

			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullDeveloperApps(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullDevelopers(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullUserAgents(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullRegionCodes(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSecurityAction_apigeeSecurityActionFullAsns(context),
			},
			{
				ResourceName:      "google_apigee_security_action.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccApigeeSecurityAction_apigeeSecurityActionFullTTL(context),
				ExpectNonEmptyPlan: true, // ttl change enforces recreation of the resource
			},
		},
	})
}

func testAccApigeeSecurityAction_apigeeBase(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_client_config" "current" {}

resource "google_compute_network" "apigee_network" {
    name = "tf-test-network-%{random_suffix}"
}

resource "google_compute_global_address" "apigee_range" {
    name          = "tf-test-address-%{random_suffix}"
    purpose       = "VPC_PEERING"
    address_type  = "INTERNAL"
    prefix_length = 16
    network       = google_compute_network.apigee_network.id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
    network                 = google_compute_network.apigee_network.id
    service                 = "servicenetworking.googleapis.com"
    reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
}

resource "google_apigee_organization" "apigee_org" {
    analytics_region   = "us-central1"
    project_id         = data.google_client_config.current.project
    authorized_network = google_compute_network.apigee_network.id
    depends_on         = [google_service_networking_connection.apigee_vpc_connection]
}

resource "google_apigee_environment" "env" {
    name         = "tf-test-env-%{random_suffix}"
    description  = "Apigee Environment"
    display_name = "environment-1"
    org_id       = google_apigee_organization.apigee_org.id
}

resource "google_apigee_addons_config" "apigee_org_security_addons_config" {
    org = google_apigee_organization.apigee_org.name
    addons_config {
        api_security_config {
            enabled = true
        }
    }
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullAllow(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        ip_address_ranges = [
            "100.0.220.1",
            "200.0.0.1",
        ]

        bot_reasons = [
            "Flooder",
            "Public Cloud Azure",
            "Public Cloud AWS",
        ]
    }

    allow {}

    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullFlag(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        ip_address_ranges = [
            "100.0.220.1",
            "200.0.0.1",
        ]

        bot_reasons = [
            "Flooder",
            "Public Cloud Azure",
            "Public Cloud AWS",
        ]
    }

    flag {
        headers {
			name  = "X-Flag-Header"
			value = "flag-value"
		}
        headers {
			name  = "X-Flag-Header-2"
			value = "flag-value-2"
		}
    }

    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullDeny(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        ip_address_ranges = [
            "100.0.220.1",
            "200.0.0.1",
        ]

        bot_reasons = [
            "Flooder",
            "Public Cloud Azure",
            "Public Cloud AWS",
        ]
    }
	
	deny {
		response_code = 403
	}

    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullHttpMethods(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        http_methods = [
			"GET",
			"POST",
			"PUT",
		]
    }
	
	deny {
		response_code = 403
	}

    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullApiKeys(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
		api_keys = [
			"foo-key",
			"bar-key",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullAccessTokens(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        access_tokens = [
			"foo-token",
			"bar-token",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullApiProducts(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        api_products = [
			"foo-product",
			"bar-product",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullDeveloperApps(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
		developer_apps = [
			"foo-app",
			"bar-app",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullDevelopers(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        developers = [
			"foo-developer",
			"bar-developer",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullUserAgents(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        user_agents = [
			"Mozilla/5.0",
			"curl/7.64.1",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullRegionCodes(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
        region_codes = [
			"US",
			"CA",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullAsns(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
		asns = [
			"23",
			"42",
		]
    }
	
	deny {
		response_code = 403
	}
	
    expire_time = "2032-12-31T23:59:59Z"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}

func testAccApigeeSecurityAction_apigeeSecurityActionFullTTL(context map[string]interface{}) string {
	return testAccApigeeSecurityAction_apigeeBase(context) + acctest.Nprintf(`
resource "google_apigee_security_action" "default" {
    security_action_id = "tf-test-%{random_suffix}"
    org_id             = google_apigee_organization.apigee_org.name
    env_id             = google_apigee_environment.env.name
    description        = "Apigee Security Action"
    state              = "ENABLED"

    condition_config {
		asns = [
			"23",
			"42",
		]
    }
	
	deny {
		response_code = 403
	}
	
    ttl 		= "3600s"
    depends_on  = [
        google_apigee_addons_config.apigee_org_security_addons_config
    ]
}
`, context)
}
