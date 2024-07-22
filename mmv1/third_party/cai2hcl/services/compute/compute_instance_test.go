package compute_test

import (
	"testing"

	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/testing"
)

func TestComputeInstance(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		"./testdata",
		[]string{
			"full_compute_instance",
			"compute_instance_iam",
		})
}
