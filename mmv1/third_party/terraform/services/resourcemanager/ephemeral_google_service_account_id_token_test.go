package resourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountIdToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken", serviceAccount)

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_basic(context),
			},
		},
	})
}

func TestAccEphemeralServiceAccountIdToken_withDelegates(t *testing.T) {
	t.Parallel()

	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "id-delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "id-delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "id-target", delegateServiceAccountEmailTwo)         // SA_4

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
		"delegate_1":              delegateServiceAccountEmailOne,
		"delegate_2":              delegateServiceAccountEmailTwo,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withDelegates(context),
			},
		},
	})
}

func TestAccEphemeralServiceAccountIdToken_withIncludeEmail(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "idtoken-email", serviceAccount)

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withIncludeEmail(context),
			},
		},
	})
}

func testAccEphemeralServiceAccountIdToken_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
ephemeral "google_service_account_id_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  target_audience       = "https://example.com"
}
`, context)
}

func testAccEphemeralServiceAccountIdToken_withDelegates(context map[string]interface{}) string {
	return acctest.Nprintf(`
ephemeral "google_service_account_id_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  delegates = [
    "%{delegate_1}",
    "%{delegate_2}",
  ]
  target_audience       = "https://example.com"
}
`, context)
}

func testAccEphemeralServiceAccountIdToken_withIncludeEmail(context map[string]interface{}) string {
	return acctest.Nprintf(`
ephemeral "google_service_account_id_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  target_audience       = "https://example.com"
  include_email        = true
}
`, context)
}
