package iam3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/iam3"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAM3OrganizationsPolicyBindingDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_full(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},

			{
				Config: testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},
		},
	})
}

func testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on = [google_iam_principal_access_boundary_policy.pab_policy]
}


resource "google_iam_organizations_policy_binding" "my_org_binding" {
  depends_on = [time_sleep.wait_60_seconds]
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
}
`, context)
}

func testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on = [google_iam_principal_access_boundary_policy.pab_policy]
}

resource "google_iam_organizations_policy_binding" "my_org_binding" {
  depends_on = [time_sleep.wait_60_seconds]
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  annotations    = {"foo": "bar"}
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
  condition {
    description  = "test condition"
    expression   = "principal.subject == 'al@a.com'"
    location     = "test location"
    title        = "test title"
  }
}
`, context)
}

func TestAccIAM3OrganizationsPolicyBinding_iamAccessPolicyBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAM3OrganizationsPolicyBindingDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIAM3OrganizationsPolicyBinding_iamAccessPolicyBinding(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},
			{
				Config: testAccIAM3OrganizationsPolicyBinding_iamAccessPolicyBinding_update(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},
		},
	})
}

func testAccIAM3OrganizationsPolicyBinding_iamAccessPolicyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "test_sa" {
  account_id   = "tf-test-sa%{random_suffix}"
  display_name = "Test Service Account for Access Policy"
}

resource "google_iam_organization_access_policy" "access_policy" {
  organization     = "%{org_id}"
  location         = "global"
  access_policy_id = "tf-test-org-policy%{random_suffix}"
  details {
    rules {
      effect      = "ALLOW"
      principals  = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test_sa.email}"]
      operation {
        permissions = ["eventarc.googleapis.com/messageBuses.publish"]
      }
    }
  }
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on      = [google_iam_organization_access_policy.access_policy]
}

resource "google_iam_organizations_policy_binding" "my_org_binding" {
  depends_on        = [time_sleep.wait_60_seconds]
  organization      = "%{org_id}"
  location          = "global"
  display_name      = "test org binding%{random_suffix}"
  policy_kind       = "ACCESS"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy            = "organizations/%{org_id}/locations/global/accessPolicies/${google_iam_organization_access_policy.access_policy.access_policy_id}"
  target {
    resource = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
}
`, context)
}

func testAccIAM3OrganizationsPolicyBinding_iamAccessPolicyBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "test_sa" {
  account_id   = "tf-test-sa%{random_suffix}"
  display_name = "Test Service Account for Access Policy"
}

resource "google_iam_organization_access_policy" "access_policy" {
  organization     = "%{org_id}"
  location         = "global"
  access_policy_id = "tf-test-org-policy%{random_suffix}"
  details {
    rules {
      effect      = "ALLOW"
      principals  = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test_sa.email}"]
      operation {
        permissions = ["eventarc.googleapis.com/messageBuses.publish"]
      }
    }
  }
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on      = [google_iam_organization_access_policy.access_policy]
}

resource "google_iam_organizations_policy_binding" "my_org_binding" {
  depends_on        = [time_sleep.wait_60_seconds]
  organization      = "%{org_id}"
  location          = "global"
  display_name      = "test org binding%{random_suffix}"
  policy_kind       = "ACCESS"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy            = "organizations/%{org_id}/locations/global/accessPolicies/${google_iam_organization_access_policy.access_policy.access_policy_id}"
  annotations       = {"foo": "bar"}
  target {
    resource = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
}
`, context)
}
