package iambeta_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/iambeta"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceIAMBetaWorkloadIdentityPoolOpenIdConfig_basic(t *testing.T) {
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
				Config: testAccDataSourceIAMBetaWorkloadIdentityPoolOpenIdConfigBasic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "issuer", regexp.MustCompile(`^https://sts\.googleapis\.com/v1/(organizations|projects)/[0-9]+/locations/global/workloadIdentityPools/agents\.global\.(org|proj)-[0-9]+\.system\.id\.goog$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "jwks_uri", regexp.MustCompile(`^https://sts\.googleapis\.com/v1/(organizations|projects)/[0-9]+/locations/global/workloadIdentityPools/agents\.global\.(org|proj)-[0-9]+\.system\.id\.goog/openid/jwks$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "authorization_endpoint", regexp.MustCompile(`^https://sts\.googleapis\.com/v1/(organizations|projects)/[0-9]+/locations/global/workloadIdentityPools/agents\.global\.(org|proj)-[0-9]+\.system\.id\.goog/authorize$`)),
					resource.TestMatchResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "token_endpoint", regexp.MustCompile(`^https://sts\.googleapis\.com/v1/(organizations|projects)/[0-9]+/locations/global/workloadIdentityPools/agents\.global\.(org|proj)-[0-9]+\.system\.id\.goog/token$`)),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "response_types_supported.#", "1"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "response_types_supported.0", "id_token"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "subject_types_supported.#", "1"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "subject_types_supported.0", "public"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "id_token_signing_alg_values_supported.#", "1"),
					resource.TestCheckResourceAttr("data.google_iam_workload_identity_pool_openid_config.example", "id_token_signing_alg_values_supported.0", "RS256"),
				),
			},
		},
	})
}

func testAccDataSourceIAMBetaWorkloadIdentityPoolOpenIdConfigBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_workload_identity_pool_openid_config" "example" {
	resource_name = "https://sts.googleapis.com/v1/organizations/%{org_id}/locations/global/workloadIdentityPools/agents.global.org-%{org_id}.system.id.goog/.well-known/openid-configuration"
}
`, context)
}
