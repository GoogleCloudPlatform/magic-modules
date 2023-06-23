package acctest

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestEnvVar = envvar.TestEnvVar

// ProviderConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func ProviderConfigEnvNames() []string {
	return envvar.ProviderConfigEnvNames()
}

var CredsEnvVars = envvar.CredsEnvVars

var ProjectNumberEnvVars = envvar.ProjectNumberEnvVars

var ProjectEnvVars = envvar.ProjectEnvVars

var FirestoreProjectEnvVars = envvar.FirestoreProjectEnvVars

var RegionEnvVars = envvar.RegionEnvVars

var ZoneEnvVars = envvar.ZoneEnvVars

var OrgEnvVars = envvar.OrgEnvVars

// This value is the Customer ID of the GOOGLE_ORG_DOMAIN workspace.
// See https://admin.google.com/ac/accountsettings when logged into an org admin for the value.
var CustIdEnvVars = envvar.CustIdEnvVars

// This value is the username of an identity account within the GOOGLE_ORG_DOMAIN workspace.
// For example in the org example.com with a user "foo@example.com", this would be set to "foo".
// See https://admin.google.com/ac/users when logged into an org admin for a list.
var IdentityUserEnvVars = envvar.IdentityUserEnvVars

var OrgEnvDomainVars = envvar.OrgEnvDomainVars

var ServiceAccountEnvVars = envvar.ServiceAccountEnvVars

var OrgTargetEnvVars = envvar.OrgTargetEnvVars

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
var BillingAccountEnvVars = envvar.BillingAccountEnvVars

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
var MasterBillingAccountEnvVars = envvar.MasterBillingAccountEnvVars

// This value is the description used for test PublicAdvertisedPrefix setup to avoid required DNS
// setup. This is only used during integration tests and would be invalid to surface to users
var PapDescriptionEnvVars = envvar.PapDescriptionEnvVars

func AccTestPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	if v := transport_tpg.MultiEnvSearch(CredsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(CredsEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ProjectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ProjectEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(RegionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(RegionEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ZoneEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ZoneEnvVars, ", "))
	}
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectNumberFromEnv() string {
	return envvar.GetTestProjectNumberFromEnv()
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return envvar.GetTestProjectFromEnv()
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	return envvar.GetTestCredsFromEnv()
}

// AccTestPreCheck ensures at least one of the region env variables is set.
func GetTestRegionFromEnv() string {
	return envvar.GetTestRegionFromEnv()
}

func GetTestZoneFromEnv() string {
	return envvar.GetTestZoneFromEnv()
}

func GetTestCustIdFromEnv(t *testing.T) string {
	return envvar.GetTestCustIdFromEnv(t)
}

func GetTestIdentityUserFromEnv(t *testing.T) string {
	return envvar.GetTestIdentityUserFromEnv(t)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
func GetTestFirestoreProjectFromEnv(t *testing.T) string {
	return envvar.GetTestFirestoreProjectFromEnv(t)
}

// Returns the raw organization id like 1234567890, skipping the test if one is
// not found.
func GetTestOrgFromEnv(t *testing.T) string {
	return envvar.GetTestOrgFromEnv(t)
}

// Alternative to GetTestOrgFromEnv that doesn't need *testing.T
// If using this, you need to process unset values at the call site
func UnsafeGetTestOrgFromEnv() string {
	return envvar.UnsafeGetTestOrgFromEnv()
}

func GetTestOrgDomainFromEnv(t *testing.T) string {
	return envvar.GetTestOrgDomainFromEnv(t)
}

func GetTestOrgTargetFromEnv(t *testing.T) string {
	return envvar.GetTestOrgTargetFromEnv(t)
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
func GetTestBillingAccountFromEnv(t *testing.T) string {
	return envvar.GetTestBillingAccountFromEnv(t)
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
func GetTestMasterBillingAccountFromEnv(t *testing.T) string {
	return envvar.GetTestMasterBillingAccountFromEnv(t)
}

func GetTestServiceAccountFromEnv(t *testing.T) string {
	return envvar.GetTestServiceAccountFromEnv(t)
}

func GetTestPublicAdvertisedPrefixDescriptionFromEnv(t *testing.T) string {
	return envvar.GetTestPublicAdvertisedPrefixDescriptionFromEnv(t)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	envvar.SkipIfEnvNotSet(t, envs...)
}
