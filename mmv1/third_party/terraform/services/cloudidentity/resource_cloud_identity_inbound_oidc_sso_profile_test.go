package cloudidentity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// The Inbound SSO APIs authorize callers with Google Workspace / Cloud Identity admin
// privileges (the "SSO settings" admin privilege) on the target customer, not Cloud IAM. These
// tests require the test credentials to hold that privilege on GOOGLE_CUST_ID.

func TestAccCloudIdentityInboundOidcSsoProfile_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityInboundOidcSsoProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityInboundOidcSsoProfile_minimal(context),
			},
			{
				ResourceName:      "google_cloud_identity_inbound_oidc_sso_profile.oidc_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIdentityInboundOidcSsoProfile_full(context),
			},
			{
				ResourceName:            "google_cloud_identity_inbound_oidc_sso_profile.oidc_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rp_config.0.client_secret"},
			},
			{
				ResourceName:       "google_cloud_identity_inbound_oidc_sso_profile.oidc_profile",
				ImportState:        true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCloudIdentityInboundOidcSsoProfile_minimal(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_inbound_oidc_sso_profile" "oidc_profile" {
  customer = "customers/%{cust_id}"

  idp_config {
    issuer_uri = "https://idp.example.com/tf-test-%{random_suffix}"
  }
}
`, context)
}

func testAccCloudIdentityInboundOidcSsoProfile_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_inbound_oidc_sso_profile" "oidc_profile" {
  customer     = "customers/%{cust_id}"
  display_name = "tf-test-oidc-%{random_suffix}-updated"

  idp_config {
    issuer_uri          = "https://idp.example.com/tf-test-%{random_suffix}-updated"
    change_password_uri = "https://idp.example.com/change-password"
  }

  rp_config {
    client_id     = "tf-test-client-%{random_suffix}"
    client_secret = "tf-test-secret-%{random_suffix}"
  }
}
`, context)
}
