package iambeta_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/iambeta"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "jwks_json", regexp.MustCompile(`(?s)^\{\s*"keys":\s*\[.*\]\s*\}\s*$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.kty", regexp.MustCompile(`^RSA$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.use", regexp.MustCompile(`^sig$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.alg", regexp.MustCompile(`^RS256$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.kid", regexp.MustCompile(`.+`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.n", regexp.MustCompile(`.+`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.e", regexp.MustCompile(`^AQAB$`)),
				),
			},
		},
	})
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
