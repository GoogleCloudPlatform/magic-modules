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

func TestAccCloudIdentityInboundSamlSsoProfile_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityInboundSamlSsoProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityInboundSamlSsoProfile_minimal(context),
			},
			{
				ResourceName:      "google_cloud_identity_inbound_saml_sso_profile.saml_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIdentityInboundSamlSsoProfile_full(context),
			},
			{
				ResourceName:      "google_cloud_identity_inbound_saml_sso_profile.saml_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:       "google_cloud_identity_inbound_saml_sso_profile.saml_profile",
				ImportState:        true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCloudIdentityInboundSamlSsoProfile_minimal(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_inbound_saml_sso_profile" "saml_profile" {
  customer = "customers/%{cust_id}"

  idp_config {
    entity_id                  = "https://idp.example.com/tf-test-%{random_suffix}"
    single_sign_on_service_uri = "https://idp.example.com/sso"
  }
}
`, context)
}

func testAccCloudIdentityInboundSamlSsoProfile_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_inbound_saml_sso_profile" "saml_profile" {
  customer     = "customers/%{cust_id}"
  display_name = "tf-test-saml-%{random_suffix}-updated"

  idp_config {
    entity_id                  = "https://idp.example.com/tf-test-%{random_suffix}-updated"
    single_sign_on_service_uri = "https://idp.example.com/sso-updated"
    logout_redirect_uri        = "https://idp.example.com/logout"
    change_password_uri        = "https://idp.example.com/change-password"
  }
}
`, context)
}
