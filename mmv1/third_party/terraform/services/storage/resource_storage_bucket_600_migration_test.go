package storage_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageBucket_retentionPeriodUpgrade(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt(t))
	retentionPeriod := 3600
	bucket := testAccStorageBucket_retentionPeriod(bucketName, retentionPeriod)
	expectedRetentionPeriod := strconv.Itoa(retentionPeriod)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: bucket,
			},
			{
				Config:             bucket,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectKnownRetentionPeriodValue("google_storage_bucket.bucket", expectedRetentionPeriod),
					},
				},
			},
		},
	})
}

func testAccStorageBucket_retentionPeriod(bucketName string, retentionPeriod int) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = %q
  location      = "US"
  retention_policy {
    retention_period = %d
  }
}
`, bucketName, retentionPeriod)
}

var _ plancheck.PlanCheck = expectKnownRetentionPeriodValue{}

type expectKnownRetentionPeriodValue struct {
	ResourceAddress         string
	ExpectedRetentionPeriod string
}

func (e expectKnownRetentionPeriodValue) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result error
	for _, change := range req.Plan.ResourceChanges {
		if change.Address == e.ResourceAddress {
			after := change.Change.After.(map[string]any)
			retentionPolicies, ok := after["retention_policy"].([]any)
			policy, ok := retentionPolicies[0].(map[string]any)
			if !ok {
				result = fmt.Errorf("Resource %q has an invalid retention_policy block", e.ResourceAddress)
				return
			}

			retentionPeriod, ok := policy["retention_period"].(string)
			if !ok {
				result = fmt.Errorf("Resource %q has an invalid retention_period", e.ResourceAddress)
			}
			if retentionPeriod != e.ExpectedRetentionPeriod {
				result = fmt.Errorf("Resource %q has retention_period %q, expected %q", e.ResourceAddress, retentionPeriod, e.ExpectedRetentionPeriod)
			}
			return
		}
	}

	resp.Error = result
}

func ExpectKnownRetentionPeriodValue(resourceAddress string, expectedRetentionPeriod string) plancheck.PlanCheck {
	return expectKnownRetentionPeriodValue{
		ResourceAddress:         resourceAddress,
		ExpectedRetentionPeriod: expectedRetentionPeriod,
	}
}
