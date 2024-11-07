package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestEphemeralServiceAccountJwt_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountJwt_withDelegates(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "delegates", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_withDelegates(targetServiceAccountEmail),
			},
		},
	})
}

func TestEphemeralServiceAccountJwt_withExpiresIn(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "expiry", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountJwt_withExpiresIn(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountJwt_basic(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountJwt_withDelegates(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  delegates             = ["%[1]s"]
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
`, serviceAccountEmail)
}

func testAccEphemeralServiceAccountJwt_withExpiresIn(serviceAccountEmail string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_jwt" "jwt" {
  target_service_account = "%s"
  expires_in            = 3600
  payload               = jsonencode({
    "sub": "%[1]s",
    "aud": "https://example.com"
  })
}
`, serviceAccountEmail)
}
