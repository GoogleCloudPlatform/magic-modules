package compute_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl"
	cai2hcl_testing "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/testing"
)

func TestComputeForwardingRule(t *testing.T) {
	cai2hcl_testing.AssertTestFiles(
		t,
		cai2hcl.ConverterMap,
		"./testdata",
		[]string{"compute_forwarding_rule"})
}
