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
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "jwks_json", regexp.MustCompile(`^\{"keys":\[.*\]\}$`)),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.kty", "RSA"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.alg", "RS256"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.use", "sig"),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.kid", regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.n", regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_jwks.example", "keys.0.e", "AQAB"),
				),
			},
		},
	})
}

func testAccDataSourceIAMBetaWorkloadIdentityPoolJwksBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_workload_identity_pool_openid_config" "oidc" {
	workload_identity_pool_id = "organizations/%{org_id}/locations/global/workloadIdentityPools/agents.global.org-%{org_id}.system.id.goog"
}

data "google_iam_workload_identity_pool_jwks" "example" {
	jwks_uri = data.google_iam_workload_identity_pool_openid_config.oidc.jwks_uri
}
`, context)
}
