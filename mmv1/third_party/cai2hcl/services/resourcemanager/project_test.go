package resourcemanager_test

import (
	"testing"

	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/testing"
)

func TestComputeInstance(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		"./testdata",
		[]string{
			"project_create",
			"project_iam",
		})
}
