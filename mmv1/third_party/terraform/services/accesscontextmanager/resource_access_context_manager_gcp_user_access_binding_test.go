package accesscontextmanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/accesscontextmanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/cloudidentity"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerGcpUserAccessBinding_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingNamedExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
  "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}

resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  group_key       = trimprefix(google_cloud_identity_group.group.id, "groups/")
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
  dry_run_access_levels = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
  session_settings {
    session_length = "3600s"
    session_length_enabled = true
    session_reauth_method = "LOGIN"
    use_oidc_max_age = false
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
          client_id = "TEST_APPLICATION"
         }
      }
    }
    active_settings {
      access_levels = [
        google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
      ]
      session_settings {
        session_length = "3600s"
        session_length_enabled = true
        session_reauth_method = "LOGIN"
        use_oidc_max_age = false
      }
    }
    dry_run_settings {
      access_levels = [
        google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
      ]
    }
  }
}
`, context)
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingNamedExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
  "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}

resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  group_key       = trimprefix(google_cloud_identity_group.group.id, "groups/")
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
  session_settings {
    session_length = "3600s"
    session_length_enabled = true
    session_reauth_method = "LOGIN"
    use_oidc_max_age = false
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
          name = "Cloud Console"
         }
      }
    }
    active_settings {
      access_levels = [
        google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
      ]
      session_settings {
        session_length = "3600s"
        session_length_enabled = true
        session_reauth_method = "LOGIN"
        use_oidc_max_age = false
      }
    }
    dry_run_settings {
      access_levels = [
        google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
      ]
    }
  }
}
`, context)
}

func testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_gcp_user_access_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(accesscontextmanager.Product, config)+"organizations/{{organization_id}}/gcpUserAccessBindings/{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AccessContextManagerGcpUserAccessBinding still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerGcpUserAccessBinding_principalTest(t *testing.T) {

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingPrincipalExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingPrincipalExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "test_sa" {
  account_id   = "tf-test-sa-%{random_suffix}"
  display_name = "Test Service Account for User Access Binding"
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
        "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}

resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  principal {
    service_account = google_service_account.test_sa.email
  }
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
}
`, context)
}

func testAccAccessContextManagerGcpUserAccessBindingRestrictedProjectScopeTest(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"org_id":         envvar.GetTestOrgFromEnv(t),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBindingRestrictedProjectScopeExample(context, "3600s"),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding_restricted_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessContextManagerGcpUserAccessBindingRestrictedProjectScopeExample(context, "7200s"),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding_restricted_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBindingRestrictedProjectScopeExample(context map[string]any, sessionLength string) string {
	context["session_length"] = sessionLength
	return acctest.Nprintf(`
resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding_restricted_project" {
  organization_id = "%{org_id}"
  principal {
    federated_principal = "principalSet://cloudresourcemanager.googleapis.com/organizations/%{org_id}/type/WorkforcePool"
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_project {
          name = "projects/%{project_number}"
        }
      }
    }
    active_settings {
      session_settings {
        session_length         = "%{session_length}"
        session_length_enabled = true
        session_reauth_method  = "LOGIN"
      }
    }
  }
}
`, context)
}

func testAccAccessContextManagerGcpUserAccessBindingRestrictedClientApplicationScopeTest(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBindingRestrictedClientApplicationScopeExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding_restricted_client_application",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBindingRestrictedClientApplicationScopeExample(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding_restricted_client_application" {
  organization_id = "%{org_id}"
  principal {
    federated_principal = "principalSet://cloudresourcemanager.googleapis.com/organizations/%{org_id}/type/WorkforcePool"
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
          name = "Gemini Enterprise mobile Android"
        }
      }
    }
    active_settings {
      session_settings {
        session_length         = "7776000s"
        session_length_enabled = true
        session_reauth_method  = "LOGIN"
      }
    }
  }
}
`, context)
}

func testAccAccessContextManagerGcpUserAccessBindingMultiScopeTest(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"org_id":         envvar.GetTestOrgFromEnv(t),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBindingMultiScopeExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding_multiscope",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBindingMultiScopeExample(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding_multiscope" {
  organization_id = "%{org_id}"
  principal {
    federated_principal = "principalSet://cloudresourcemanager.googleapis.com/organizations/%{org_id}/type/WorkforcePool"
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_project {
          name = "projects/%{project_number}"
        }
      }
    }
    active_settings {
      session_settings {
        session_length         = "7200s"
        session_length_enabled = true
        session_reauth_method  = "LOGIN"
      }
    }
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
          name = "Gemini Enterprise mobile Android"
        }
      }
    }
    active_settings {
      session_settings {
        session_length         = "7776000s"
        session_length_enabled = true
        session_reauth_method  = "LOGIN"
      }
    }
  }
}
`, context)
}
