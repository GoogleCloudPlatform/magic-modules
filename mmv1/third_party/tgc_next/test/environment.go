package test

import (
	"log"
	"os"
	"path/filepath"
)

const (
	samplePolicyPath          = "../testdata/sample_policies"
	defaultAncestry           = "organization/529579013760/project/1067888929963"
	defaultBillingAccount     = "000AA0-A0B00A-AA00AA"
	defaultCustId             = "A00ccc00a"
	defaultFolder             = "67890"
	defaultIdentityUser       = "foo"
	defaultOrganization       = "529579013760"
	defaultOrganizationDomain = "meep.test.com"
	defaultOrganizationTarget = "13579"
	defaultProject            = "ci-test-project-nightly-beta"
	defaultProviderVersion    = "5.5.0" // if dev override is enabled, the provider version is ignored in terraform execution
	defaultRegion             = "us-central1"
	defaultServiceAccount     = "meep@foobar.iam.gserviceaccount.com"
)

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

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
