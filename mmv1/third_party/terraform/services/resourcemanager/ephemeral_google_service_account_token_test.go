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
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountToken_withDelegates(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "delegates", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withDelegates(targetServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountToken_withCustomLifetime(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "lifetime", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withCustomLifetime(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountToken_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_token" "token" {
  target_service_account = "%s"
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountToken_withDelegates(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_token" "token" {
  target_service_account = "%s"
  delegates             = ["%[1]s"]
  lifetime             = "1200s"
  scopes               = ["https://www.googleapis.com/auth/cloud-platform"]
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountToken_withCustomLifetime(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_token" "token" {
  target_service_account = "%s"
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime              = "3600s"
}
`, serviceAccountEmail)
}
