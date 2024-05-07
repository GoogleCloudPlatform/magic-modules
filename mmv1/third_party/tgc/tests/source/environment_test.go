package test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	samplePolicyPath          = "../testdata/sample_policies"
	defaultAncestry           = "organization/12345/folder/67890"
	defaultBillingAccount     = "000AA0-A0B00A-AA00AA"
	defaultCustId             = "A00ccc00a"
	defaultFolder             = "67890"
	defaultIdentityUser       = "foo"
	defaultOrganization       = "12345"
	defaultOrganizationDomain = "meep.test.com"
	defaultOrganizationTarget = "13579"
	defaultProject            = "foobar"
	defaultProviderVersion    = "5.5.0" // if dev override is enabled, the provider version is ignored in terraform execution
	defaultRegion             = "us-central1"
	defaultServiceAccount     = "meep@foobar.iam.gserviceaccount.com"
)

// AccTestPreCheck ensures at least one of the project env variables is set.
func getTestProjectFromEnv() string {
	project := multiEnvSearch([]string{"TEST_PROJECT", "GOOGLE_PROJECT"})
	if project == "" {
		log.Printf("Missing required env var TEST_PROJECT. Default (%s) will be used.", defaultProject)
		project = defaultProject
	}

	return project
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
func getTestCredsFromEnv() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current directory: %v", err)
	}

	credentials := multiEnvSearch([]string{"TEST_CREDENTIALS", "GOOGLE_APPLICATION_CREDENTIALS"})
	if credentials != "" {
		// Make credentials path relative to repo root rather than
		// test/ dir if it is a relative path.
		if !filepath.IsAbs(credentials) {
			credentials = filepath.Join(cwd, "..", credentials)
		}
	} else {
		log.Printf("missing env var TEST_CREDENTIALS, will try to use Application Default Credentials")
	}

	return credentials
}

// AccTestPreCheck ensures at least one of the region env variables is set.
func getTestRegionFromEnv() string {
	return defaultRegion
}

func getTestCustIdFromEnv(t *testing.T) string {
	return defaultCustId
}

func getTestIdentityUserFromEnv(t *testing.T) string {
	return defaultIdentityUser
}

func getTestOrgFromEnv(t *testing.T) string {
	org, ok := os.LookupEnv("TEST_ORG_ID")
	if !ok {
		log.Printf("Missing required env var TEST_ORG_ID. Default (%s) will be used.", defaultOrganization)
		org = defaultOrganization
	}

	return org
}

func getTestOrgDomainFromEnv(t *testing.T) string {
	return defaultOrganizationDomain
}

func getTestOrgTargetFromEnv(t *testing.T) string {
	return defaultOrganizationTarget
}

func getTestBillingAccountFromEnv(t *testing.T) string {
	return defaultBillingAccount
}

func getTestServiceAccountFromEnv(t *testing.T) string {
	return defaultServiceAccount
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func shouldOutputGeneratedFiles() bool {
	_, ok := os.LookupEnv("TFV_CREATE_GENERATED_FILES")
	return ok
}
