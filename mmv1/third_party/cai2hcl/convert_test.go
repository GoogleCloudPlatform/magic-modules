package cai2hcl_test

import (
	cai2hclTesting "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/testing"
	"testing"
)

func TestConvertCompute(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		"./services/compute/testdata",
		[]string{
			"full_compute_instance",
		})
}

func TestConvertResourcemanager(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		"./services/resourcemanager/testdata",
		[]string{
			"project_create",
		})
}

func TestConvertNetworksecurity(t *testing.T) {
	cai2hclTesting.AssertTestFiles(
		t,
		"./services/networksecurity/testdata",
		[]string{
			"server_tls_policy",
			"backend_authentication_config",
		})
}
