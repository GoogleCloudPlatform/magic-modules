package resourcemanager

import (
	"testing"

	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/testing"
)

func TestComputeInstance(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		UberConverter,
		"./testdata",
		[]string{
			"project_create",
			"project_iam",
		})
}
