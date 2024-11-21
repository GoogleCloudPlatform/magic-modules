package resourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountToken_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "basic", serviceAccount)

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
		"scope_1":                 "https://www.googleapis.com/auth/cloud-platform",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_basic(context),
			},
		},
	})
}

func TestAccEphemeralServiceAccountToken_withDelegates(t *testing.T) {
	t.Parallel()

	initialServiceAccount := envvar.GetTestServiceAccountFromEnv(t)
	delegateServiceAccountEmailOne := acctest.BootstrapServiceAccount(t, "delegate1", initialServiceAccount)          // SA_2
	delegateServiceAccountEmailTwo := acctest.BootstrapServiceAccount(t, "delegate2", delegateServiceAccountEmailOne) // SA_3
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "target", delegateServiceAccountEmailTwo)         // SA_4

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
		"delegate_1":              delegateServiceAccountEmailOne,
		"delegate_2":              delegateServiceAccountEmailTwo,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withDelegates(context),
			},
		},
	})
}

func TestAccEphemeralServiceAccountToken_withCustomLifetime(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "lifetime", serviceAccount)

	context := map[string]interface{}{
		"ephemeral_resource_name": "token",
		"target_service_account":  targetServiceAccountEmail,
		"lifetime":                "3600s",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountToken_withCustomLifetime(context),
			},
		},
	})
}

func testAccEphemeralServiceAccountToken_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
ephemeral "google_service_account_access_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  scopes                = ["%{scope_1}"]
}
`, context)
}

func testAccEphemeralServiceAccountToken_withDelegates(context map[string]interface{}) string {
	return acctest.Nprintf(`
ephemeral "google_service_account_access_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  delegates = [
    "%{delegate_1}",
    "%{delegate_2}",
  ]
  scopes = ["https://www.googleapis.com/auth/cloud-platform"]
}
`, context)
}

func testAccEphemeralServiceAccountToken_withCustomLifetime(context map[string]interface{}) string {
	return acctest.Nprintf(`

ephemeral "google_service_account_access_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  scopes                = ["https://www.googleapis.com/auth/cloud-platform"]
  lifetime              = "3600s"
}
`, context)
}
