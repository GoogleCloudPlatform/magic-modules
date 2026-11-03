package cloudidentity_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var regexpAllowlistedDomainName = regexp.MustCompile(`^allowlistedDomains/[^/]+$`)

var allowlistedDomainScopes = []interface{}{
	"https://www.googleapis.com/auth/cloud-identity.allowlisteddomains",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
}

// The Allowlisted Domains API acts on the customer of the caller, which must be a
// Google Workspace / Cloud Identity admin with the allowlisted domain management
// privilege, using the cloud-identity.allowlisteddomains (or cloud-identity) scope.
func TestAccCloudIdentityAllowlistedDomain_basic(t *testing.T) {
	// Not parallel: every create/delete rewrites the customer's whole allowlist,
	// and the API rejects a mutation while another one is pending.

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityAllowlistedDomainDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityAllowlistedDomain_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_identity_allowlisted_domain.allowlisted_domain", "domain", fmt.Sprintf("tf-test-%s.com", context["random_suffix"])),
					resource.TestMatchResourceAttr("google_cloud_identity_allowlisted_domain.allowlisted_domain", "name", regexpAllowlistedDomainName),
				),
			},
			{
				ResourceName:      "google_cloud_identity_allowlisted_domain.allowlisted_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:       "google_cloud_identity_allowlisted_domain.allowlisted_domain",
				ImportState:        true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCloudIdentityAllowlistedDomain_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  scopes = [
    "https://www.googleapis.com/auth/cloud-identity.allowlisteddomains",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_cloud_identity_allowlisted_domain" "allowlisted_domain" {
  domain = "tf-test-%{random_suffix}.com"
}
`, context)
}

func testAccCheckCloudIdentityAllowlistedDomainDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloud_identity_allowlisted_domain" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config, err := allowlistedDomainTestConfig(t)
			if err != nil {
				return err
			}

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CloudIdentityBasePath}}{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("CloudIdentityAllowlistedDomain still exists at %s", url)
			}
			if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				return fmt.Errorf("Error checking CloudIdentityAllowlistedDomain %s: %w", url, err)
			}
		}

		return nil
	}
}

// allowlistedDomainTestConfig returns a provider config whose credentials carry
// the Allowlisted Domains scope. Under VCR, the cached config of the test's own
// (scoped) provider is reused so requests are recorded/replayed; otherwise the
// default test config would only have the cloud-platform scope, which the API
// rejects.
func allowlistedDomainTestConfig(t *testing.T) (*transport_tpg.Config, error) {
	if acctest.IsVcrEnabled() {
		return acctest.GoogleProviderConfig(t), nil
	}
	sdkProvider := provider.Provider()
	if diags := sdkProvider.Configure(context.Background(), sdkterraform.NewResourceConfigRaw(map[string]interface{}{
		"scopes": allowlistedDomainScopes,
	})); diags.HasError() {
		return nil, fmt.Errorf("Error configuring provider: %v", diags)
	}
	return sdkProvider.Meta().(*transport_tpg.Config), nil
}
