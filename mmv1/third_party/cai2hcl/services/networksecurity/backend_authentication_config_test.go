package networksecurity_test

import (
	"testing"

	cai2hcl_testing "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cai2hcl/testing"
)

func TestBackendAuthenticationConfig(t *testing.T) {
	cai2hcl_testing.AssertTestFiles(
		t,
		"./testdata",
		[]string{"backend_authentication_config"})
}
