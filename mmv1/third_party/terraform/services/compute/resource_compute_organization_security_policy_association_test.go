package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeOrganizationSecurityPolicyAssociation_excludeFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeOrganizationSecurityPolicyAssociationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeOrganizationSecurityPolicyAssociation_excludeFieldsCreate(context),
			},
			{
				ResourceName:            "google_compute_organization_security_policy_association.policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy_id"},
			},
			{
				Config: testAccComputeOrganizationSecurityPolicyAssociation_excludeFieldsUpdate(context),
			},
			{
				ResourceName:            "google_compute_organization_security_policy_association.policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy_id"},
			},
		},
	})
}

func testAccComputeOrganizationSecurityPolicyAssociation_excludeFieldsCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "security_policy_target" {
  display_name = "tf-test-secpol-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

resource "google_compute_organization_security_policy" "policy" {
  short_name   = "tf-test%{random_suffix}"
  parent       = google_folder.security_policy_target.name
  type         = "CLOUD_ARMOR"
}

resource "google_compute_organization_security_policy_association" "policy" {
  name          = "tf-test%{random_suffix}"
  attachment_id = google_compute_organization_security_policy.policy.parent
  policy_id     = google_compute_organization_security_policy.policy.id
  excluded_projects = [
    "projects/12345678910",
    "projects/01987654321"
  ]
  excluded_folders = [
    "projects/12345678910",
    "projects/01987654321"
  ]
}
`, context)
}

func testAccComputeOrganizationSecurityPolicyAssociation_excludeFieldsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "security_policy_target" {
  display_name = "tf-test-secpol-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

resource "google_compute_organization_security_policy" "policy" {
  short_name   = "tf-test%{random_suffix}"
  parent       = google_folder.security_policy_target.name
  type         = "CLOUD_ARMOR"
}

resource "google_compute_organization_security_policy_association" "policy" {
  name          = "tf-test%{random_suffix}"
  attachment_id = google_compute_organization_security_policy.policy.parent
  policy_id     = google_compute_organization_security_policy.policy.id
  excluded_projects = [
    "projects/01987654321"
  ]
  excluded_folders = [
    "projects/00000000000",
    "projects/12345678910",
    "projects/01987654321"
  ]
}
`, context)
}
