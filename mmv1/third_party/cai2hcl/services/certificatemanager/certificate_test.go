package certificatemanager_test

import (
	cai2hcl_testing "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cai2hcl/testing"
	"testing"
)

func TestCertificate(t *testing.T) {
	cai2hcl_testing.AssertTestFiles(
		t,
		"./testdata",
		[]string{"certificate"})
}
