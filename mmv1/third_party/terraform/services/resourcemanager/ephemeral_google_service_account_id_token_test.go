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

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "target", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withDelegates(targetServiceAccountEmail),
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

func testAccEphemeralServiceAccountIdToken_withDelegates(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_id_token" "token" {
  target_service_account = "%s"
  target_audience       = "https://example.com"
  delegates            = ["%[1]s"]
}
`, serviceAccountEmail)
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
