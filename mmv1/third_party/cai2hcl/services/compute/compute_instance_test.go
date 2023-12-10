package compute_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl"
	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/testing"
)

func TestComputeInstance(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		cai2hcl.ConverterMap,
		"./testdata",
		[]string{
			"full_compute_instance",
			"compute_instance_iam",
		})
}
