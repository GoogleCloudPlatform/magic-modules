package iambeta_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/iambeta"
)

func TestAccDataSourceIAMBetaWorkloadIdentityPoolJwks_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		// Pre-existing Google-owned organization (google.com) with Agent System Pool
		"org_id": "433637338589",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIAMBetaWorkloadIdentityPoolJwksBasic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "jwks_json", regexp.MustCompile(`(?s)^\{\s*"keys":\s*\[.*\]\s*\}$`)),
					testAccCheckJWKSKeys("data.google_iam_workload_identity_pool_jwks.example"),
				),
			},
		},
	})
}

func testAccCheckJWKSKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		keysCountStr, ok := rs.Primary.Attributes["keys.#"]
		if !ok {
			return fmt.Errorf("Attribute 'keys.#' not found")
		}

		keysCount, err := strconv.Atoi(keysCountStr)
		if err != nil {
			return fmt.Errorf("Error parsing keys count: %s", err)
		}

		// Only check that key attributes are present when keys list is not empty
		if keysCount > 0 {
			fields := []string{"kty", "use", "kid", "n", "e", "alg"}
			for _, field := range fields {
				attrKey := fmt.Sprintf("keys.0.%s", field)
				if val, exists := rs.Primary.Attributes[attrKey]; !exists || val == "" {
					return fmt.Errorf("Expected %s to be present and non-empty in keys.0", attrKey)
				}
			}
		}

		return nil
	}
}

func testAccDataSourceIAMBetaWorkloadIdentityPoolJwksBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_workload_identity_pool_openid_config" "oidc" {
	resource_name = "https://sts.googleapis.com/v1/organizations/%{org_id}/locations/global/workloadIdentityPools/agents.global.org-%{org_id}.system.id.goog/.well-known/openid-configuration"
}

data "google_iam_workload_identity_pool_jwks" "example" {
	resource_name = data.google_iam_workload_identity_pool_openid_config.oidc.jwks_uri
}
`, context)
}
