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
		"ephemeral_reference":     "ephemeral.google_service_account_id_token.token",
		"target_service_account":  targetServiceAccountEmail,
		"target_audience":         "https://example.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_basic(context),
				Check: resource.ComposeTestCheckFunc(
					// Assert exact values
					resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.target_service_account", context["target_service_account"].(string)),
					resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.target_audience", context["target_audience"].(string)),
					// Assert set
					resource.TestCheckResourceAttrSet(acctest.EchoResourceName, "data.id_token"),
					// Assert unset (is unset/null in resources)
					resource.TestCheckNoResourceAttr(acctest.EchoResourceName, "data.include_email"),
				),
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
		"ephemeral_reference":     "ephemeral.google_service_account_id_token.token",
		"target_service_account":  targetServiceAccountEmail,
		"delegate_1":              delegateServiceAccountEmailOne,
		"delegate_2":              delegateServiceAccountEmailTwo,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withDelegates(context),
				Check: resource.ComposeTestCheckFunc(
					// Assert exact values
					resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.delegates.0", context["delegate_1"].(string)),
					resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.delegates.1", context["delegate_2"].(string)),
					// Assert set
					resource.TestCheckResourceAttrSet(acctest.EchoResourceName, "data.id_token"),
					// Assert unset (is unset/null in resources)
					resource.TestCheckNoResourceAttr(acctest.EchoResourceName, "data.include_email"),
				),
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
		"ephemeral_reference":     "ephemeral.google_service_account_id_token.token",
		"target_service_account":  targetServiceAccountEmail,
		"include_email":           "true",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountIdToken_withIncludeEmail(context),
				Check: resource.ComposeTestCheckFunc(
					// Assert exact values
					resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.include_email", context["include_email"].(string)),
					// Assert set
					resource.TestCheckResourceAttrSet(acctest.EchoResourceName, "data.id_token"),
				),
			},
		},
	})
}

func testAccEphemeralServiceAccountIdToken_basic(context map[string]interface{}) string {
	return acctest.EchoResourceConfig(context["ephemeral_reference"].(string)) + acctest.Nprintf(`
ephemeral "google_service_account_id_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  target_audience       = "%{target_audience}"
}
`, context)
}

func testAccEphemeralServiceAccountIdToken_withDelegates(context map[string]interface{}) string {
	return acctest.EchoResourceConfig(context["ephemeral_reference"].(string)) + acctest.Nprintf(`
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
	return acctest.EchoResourceConfig(context["ephemeral_reference"].(string)) + acctest.Nprintf(`
ephemeral "google_service_account_id_token" "%{ephemeral_resource_name}" {
  target_service_account = "%{target_service_account}"
  target_audience       = "https://example.com"
  include_email        = %{include_email}
}
`, context)
}
