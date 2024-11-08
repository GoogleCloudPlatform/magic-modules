package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestEphemeralServiceAccountIdToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountIdToken_withDelegates(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "target", delegateServiceAccountEmailTwo)         // SA_4

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_delegatesSetup(initialServiceAccount, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project),
			},
			{
				Config: testAccEphemeralServiceAccountIdToken_withDelegates(initialServiceAccount, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project),
			},
		},
	})
}

func TestEphemeralServiceAccountIdToken_withIncludeEmail(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken-email", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withIncludeEmail(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountIdToken_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  target_audience       = "https://example.com"
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountIdToken_withDelegates(initialServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project string) string {
	return fmt.Sprintf(`
resource "google_service_account_iam_binding" "sa2_to_sa3" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[4]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[3]s"
  ]
  depends_on = [google_service_account_iam_binding.sa1_to_sa2]
}

resource "google_service_account_iam_binding" "sa1_to_sa2" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[3]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[2]s"
  ]
  depends_on = [google_service_account_iam_binding.terraform_to_delegate1]
}

resource "google_service_account_iam_binding" "terraform_to_delegate1" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[2]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[1]s"
  ]
  depends_on = [google_project_iam_member.terraform_sa_token_creator]
}

resource "google_project_iam_member" "terraform_sa_token_creator" {
  project = "%[5]s"
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:%[1]s"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [
    google_service_account_iam_binding.sa1_to_sa2,
    google_service_account_iam_binding.sa2_to_sa3,
    google_project_iam_member.terraform_sa_token_creator,
  ]
  create_duration = "60s"
}

ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%[4]s"
  delegates = [
    "%[3]s",
    "%[2]s",
  ]
  target_audience       = "https://example.com"
}

# The delegation chain is:
# SA_1 (initialServiceAccountEmail) -> SA_2 (delegateServiceAccountEmailOne) -> SA_3 (delegateServiceAccountEmailTwo) -> SA_4 (targetServiceAccountEmail)
`, initialServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project)
}

func testAccEphemeralServiceAccountIdToken_delegatesSetup(initialServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project string) string {
	return fmt.Sprintf(`
resource "google_service_account_iam_binding" "sa2_to_sa3" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[4]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[3]s"
  ]
  depends_on = [google_service_account_iam_binding.sa1_to_sa2]
}

resource "google_service_account_iam_binding" "sa1_to_sa2" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[3]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[2]s"
  ]
  depends_on = [google_service_account_iam_binding.terraform_to_delegate1]
}

resource "google_service_account_iam_binding" "terraform_to_delegate1" {
  service_account_id = "projects/%[5]s/serviceAccounts/%[2]s"
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = [
    "serviceAccount:%[1]s"
  ]
  depends_on = [google_project_iam_member.terraform_sa_token_creator]
}

resource "google_project_iam_member" "terraform_sa_token_creator" {
  project = "%[5]s"
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:%[1]s"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [
    google_service_account_iam_binding.sa1_to_sa2,
    google_service_account_iam_binding.sa2_to_sa3,
    google_project_iam_member.terraform_sa_token_creator,
  ]
  create_duration = "60s"
}

# The delegation chain is:
# SA_1 (initialServiceAccountEmail) -> SA_2 (delegateServiceAccountEmailOne) -> SA_3 (delegateServiceAccountEmailTwo) -> SA_4 (targetServiceAccountEmail)
`, initialServiceAccountEmail, delegateServiceAccountEmailOne, delegateServiceAccountEmailTwo, targetServiceAccountEmail, project)
}

func testAccEphemeralServiceAccountIdToken_withIncludeEmail(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  target_audience       = "https://example.com"
  include_email        = true
}
`, serviceAccountEmail)
}
