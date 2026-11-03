package cloudidentity_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	_ "github.com/hashicorp/terraform-provider-google/google/services/cloudidentity"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// The Cloud Identity Policy API can only be called by a Google Workspace / Cloud Identity
// super administrator (service accounts must use domain-wide delegation impersonating one), and
// it only accepts the https://www.googleapis.com/auth/cloud-identity.policies OAuth scope
// (cloud-platform is rejected), so the provider is configured with that scope below.
// Create, patch and delete are only supported for rule and detector settings, so a DLP word list
// detector is used; the customer must have a Workspace edition that supports DLP.
//
// Required environment variables:
//   - GOOGLE_CUST_ID: the customer ID (without the "customers/" prefix).
//   - GOOGLE_CLOUD_IDENTITY_POLICY_ORG_UNIT: an org unit ID in that customer (without the
//     "orgUnits/" prefix) that the test policy is applied to.
const cloudIdentityPolicyOrgUnitEnvVar = "GOOGLE_CLOUD_IDENTITY_POLICY_ORG_UNIT"

var cloudIdentityPolicyScopes = []interface{}{
	"https://www.googleapis.com/auth/cloud-identity.policies",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
}

func TestAccCloudIdentityPolicy_cloudidentityPolicyBasic(t *testing.T) {
	orgUnit := os.Getenv(cloudIdentityPolicyOrgUnitEnvVar)

	context := map[string]interface{}{
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"org_unit":      orgUnit,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
			if orgUnit == "" {
				t.Fatalf("%s must be set to an org unit ID in the test customer", cloudIdentityPolicyOrgUnitEnvVar)
			}
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityPolicy_cloudidentityPolicyBasic(context),
			},
			{
				ResourceName:      "google_cloud_identity_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIdentityPolicy_cloudidentityPolicyBasic_update(context),
			},
			{
				ResourceName:      "google_cloud_identity_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:       "google_cloud_identity_policy.primary",
				ImportState:        true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCloudIdentityPolicy_cloudidentityPolicyBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  scopes = [
    "https://www.googleapis.com/auth/cloud-identity.policies",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_cloud_identity_policy" "primary" {
  customer = "customers/%{cust_id}"

  policy_query {
    org_unit = "orgUnits/%{org_unit}"
  }

  setting {
    type = "settings/detector.word_list"
    value_json = jsonencode({
      displayName = "tf-test-detector-%{random_suffix}"
      description = "Terraform acceptance test detector."
      wordList = {
        words = ["confidential", "top-secret"]
      }
    })
  }
}
`, context)
}

func testAccCloudIdentityPolicy_cloudidentityPolicyBasic_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  scopes = [
    "https://www.googleapis.com/auth/cloud-identity.policies",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_cloud_identity_policy" "primary" {
  customer = "customers/%{cust_id}"

  policy_query {
    org_unit = "orgUnits/%{org_unit}"
  }

  setting {
    type = "settings/detector.word_list"
    value_json = jsonencode({
      displayName = "tf-test-detector-%{random_suffix}-updated"
      description = "Terraform acceptance test detector (updated)."
      wordList = {
        words = ["confidential", "top-secret", "internal-only"]
      }
    })
  }
}
`, context)
}

func testAccCheckCloudIdentityPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloud_identity_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config, err := cloudIdentityPolicyTestConfig(t)
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
				return fmt.Errorf("CloudIdentityPolicy still exists at %s", url)
			}
			if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				return fmt.Errorf("Error checking CloudIdentityPolicy %s: %s", url, err)
			}
		}

		return nil
	}
}

// cloudIdentityPolicyTestConfig returns a provider config whose credentials carry the
// cloud-identity.policies scope. Under VCR, the cached config of the test's own (scoped)
// provider is reused so requests are recorded/replayed; otherwise the default test config would
// only have the cloud-platform scope, which the API rejects.
func cloudIdentityPolicyTestConfig(t *testing.T) (*transport_tpg.Config, error) {
	if acctest.IsVcrEnabled() {
		return acctest.GoogleProviderConfig(t), nil
	}
	sdkProvider := provider.Provider()
	if diags := sdkProvider.Configure(context.Background(), sdkterraform.NewResourceConfigRaw(map[string]interface{}{
		"scopes": cloudIdentityPolicyScopes,
	})); diags.HasError() {
		return nil, fmt.Errorf("Error configuring provider: %v", diags)
	}
	return sdkProvider.Meta().(*transport_tpg.Config), nil
}
