package compute

import (
	"testing"

	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/testing"
)

func TestComputeForwardingRule(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		ConverterNames, ConverterMap,
		"./testdata",
		[]string{
			"full_compute_forwarding_rule",
		})
}
