package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// WARNING: RECORDING this test creates a real, non-deletable multi-year commitment
// (RegionCommitment has exclude_delete: true). Use VCR REPLAYING in CI; only RECORD
// on a sacrificial project where orphaned commitments are acceptable.
func TestAccComputeRegionCommitment_resourceManagerTags(t *testing.T) {

	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	suffixName := acctest.RandString(t, 10)
	tagKeyResult := acctest.BootstrapSharedTestTagKeyDetails(t, "crm-region-commitment-tagkey", "organizations/"+org, make(map[string]interface{}))
	sharedTagkey, _ := tagKeyResult["shared_tag_key"]
	tagValueResult := acctest.BootstrapSharedTestTagValueDetails(t, "crm-region-commitment-tagvalue", sharedTagkey, org)

	context := map[string]interface{}{
		"commitment_name": fmt.Sprintf("tf-test-commitment-rmt-%s", suffixName),
		"tag_key_id":      tagKeyResult["name"],
		"tag_value_id":    tagValueResult["name"],
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// No CheckDestroy: RegionCommitment has exclude_delete: true.
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionCommitment_resourceManagerTags(context),
			},
		},
	})
}

func testAccComputeRegionCommitment_resourceManagerTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_commitment" "foobar" {
  name = "%{commitment_name}"
  plan = "TWELVE_MONTH"
  resources {
    type   = "VCPU"
    amount = "1"
  }
  resources {
    type   = "MEMORY"
    amount = "1"
  }
  params {
    resource_manager_tags = {
      "%{tag_key_id}" = "%{tag_value_id}"
    }
  }
}
`, context)
}
