package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRegionSecurityPolicyDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicyDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_region_security_policy.default", "google_compute_region_security_policy.default"),
				),
			},
		},
	})
}

func testAccComputeRegionSecurityPolicyDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "default" {
  name        = "tf-test-region-sec-policy-%{random_suffix}"
  region      = "us-west2"
  description = "basic region security policy"
  type        = "CLOUD_ARMOR"
}

data "google_compute_region_security_policy" "default" {
  name   = google_compute_region_security_policy.default.name
  region = "us-west2"
}
`, context)
}
