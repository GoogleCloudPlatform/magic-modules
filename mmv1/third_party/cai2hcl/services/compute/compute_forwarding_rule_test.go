package compute

import (
	"testing"

	cai2hcl_testing "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/testing"
)

func TestComputeForwardingRule(t *testing.T) {
	cai2hcl_testing.AssertTestFiles(
		t,
		UberConverter,
		"./testdata",
		[]string{"compute_forwarding_rule"})
}
