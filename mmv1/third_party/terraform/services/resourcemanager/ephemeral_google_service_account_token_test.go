package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestEphemeralServiceAccountToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "acctoken", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_basic(targetServiceAccountEmail, serviceAccount),
			},
		},
	})
}

func TestEphemeralServiceAccountToken_withDelegates(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "acctoken-delegate", serviceAccount)
	delegateServiceAccountEmail := acctest.BootstrapServiceAccount(t, "acctoken-delegate-sa", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withDelegates(targetServiceAccountEmail, delegateServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountToken_withCustomLifetime(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "acctoken-lifetime", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withCustomLifetime(targetServiceAccountEmail, serviceAccount),
			},
		},
	})
}

func testAccEphemeralServiceAccountToken_basic(serviceAccountEmail, serviceAccountId string) string {
	return fmt.Sprintf(`
resource "google_service_account_iam_member" "token_creator" {
  service_account_id = "projects/%s/serviceAccounts/%s"
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:%s"
}

// Add a time delay to allow IAM changes to propagate
resource "time_sleep" "wait_30_seconds" {
  depends_on = [google_service_account_iam_member.token_creator]
  create_duration = "10s"
}

ephemeral "google_service_account_token" "token" {
  target_service_account = %q
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime              = "3600s"
  depends_on = [time_sleep.wait_30_seconds]
}
`, envvar.GetTestProjectFromEnv(), serviceAccountEmail, serviceAccountId, serviceAccountEmail)
}

func testAccEphemeralServiceAccountToken_withDelegates(serviceAccountEmail, delegateEmail string) string {
	return fmt.Sprintf(`
resource "google_service_account_iam_member" "token_creator" {
  service_account_id = "projects/%s/serviceAccounts/%s"
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:%s"
}

// Add a time delay to allow IAM changes to propagate
resource "time_sleep" "wait_30_seconds" {
  depends_on = [google_service_account_iam_member.token_creator]
  create_duration = "10s"
}

ephemeral "google_service_account_token" "token" {
  target_service_account = %q
  delegates             = [%q]
  lifetime             = "1200s"
  scopes               = ["https://www.googleapis.com/auth/cloud-platform"]
  depends_on = [time_sleep.wait_30_seconds]
}
`, envvar.GetTestProjectFromEnv(), serviceAccountEmail, delegateEmail, serviceAccountEmail, delegateEmail)
}

func testAccEphemeralServiceAccountToken_withCustomLifetime(serviceAccountEmail, serviceAccountId string) string {
	return fmt.Sprintf(`
resource "google_service_account_iam_member" "token_creator" {
  service_account_id = "projects/%s/serviceAccounts/%s"
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:%s"
}

// Add a time delay to allow IAM changes to propagate
resource "time_sleep" "wait_30_seconds" {
  depends_on = [google_service_account_iam_member.token_creator]
  create_duration = "10s"
}

ephemeral "google_service_account_token" "token" {
  target_service_account = %q
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime              = "3600s"
  depends_on = [time_sleep.wait_30_seconds]
}
`, envvar.GetTestProjectFromEnv(), serviceAccountEmail, serviceAccountId, serviceAccountEmail)
}
