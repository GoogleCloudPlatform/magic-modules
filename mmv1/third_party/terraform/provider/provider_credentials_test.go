package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestAccSdkProvider_credentials is a series of acc tests asserting how the SDK provider handles credentials arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_credentials
func TestAccSdkProvider_credentials(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"credentials can be configured as a path to a credentials JSON file":                                       testAccSdkProvider_credentials_validJsonFilePath,
		"configuring credentials as a path to a non-existent file results in an error":                             testAccSdkProvider_credentials_badJsonFilepathCausesError,
		"config takes precedence over environment variables":                                                       testAccSdkProvider_credentials_configPrecedenceOverEnvironmentVariables,
		"when credentials is unset in the config, environment variables are used in a given order":                 testAccSdkProvider_credentials_precedenceOrderEnvironmentVariables, // GOOGLE_CREDENTIALS, GOOGLE_CLOUD_KEYFILE_JSON, GCLOUD_KEYFILE_JSON, GOOGLE_APPLICATION_CREDENTIALS
		"when credentials is set to an empty string in the config the value isn't ignored and results in an error": testAccSdkProvider_credentials_emptyStringValidation,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSdkProvider_credentials_validJsonFilePath(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// unset all credentials env vars
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, "")
	}

	credentials := transport_tpg.TestFakeCredentialsPath

	context := map[string]interface{}{
		"credentials":   credentials,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             testAccSdkProvider_credentialsInProviderBlock(context),
				PlanOnly:           true, // Path to file is valid but contents aren't; apply would error
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccSdkProvider_credentials_badJsonFilepathCausesError(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// unset all credentials env vars
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, "")
	}

	pathToMissingFile := "./this/path/does/not/exist.json" // Doesn't exist

	context := map[string]interface{}{
		"credentials":   pathToMissingFile,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// credentials is a path to a json, but if that file doesn't exist so there's an error
				Config:      testAccSdkProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
	})
}

func testAccSdkProvider_credentials_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	credentials := envvar.GetTestCredsFromEnv()

	// ensure all possible credentials env vars set; show they aren't used
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, credentials)
	}

	pathToMissingFile := "./this/path/does/not/exist.json" // Doesn't exist

	context := map[string]interface{}{
		"credentials":   pathToMissingFile,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
	})
}

func testAccSdkProvider_credentials_precedenceOrderEnvironmentVariables(t *testing.T) {
	/*
		These are all the ENVs for credentials, and they are in order of precedence.
		GOOGLE_CREDENTIALS
		GOOGLE_CLOUD_KEYFILE_JSON
		GCLOUD_KEYFILE_JSON
		GOOGLE_APPLICATION_CREDENTIALS
	*/

	goodCredentials := envvar.GetTestCredsFromEnv()
	badCreds := acctest.GenerateFakeCredentialsJson("test")
	badCredsPath := "./this/path/does/not/exist.json" // Doesn't exist

	context := map[string]interface{}{
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		// ProtoV5ProviderFactories set on each step to ensure provider reconfigured each time
		Steps: []resource.TestStep{
			{
				// Error as all ENVs set to 'bad' creds
				PreConfig: func() {
					for _, v := range envvar.CredsEnvVars {
						t.Setenv(v, badCreds)
					}
				},
				Config:                   testAccSdkProvider_credentialsInEnvsOnly(context),
				ExpectError:              regexp.MustCompile("private key should be a PEM or plain PKCS1 or PKCS8"),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				// GOOGLE_CREDENTIALS is used 1st if set
				PreConfig: func() {
					// good
					t.Setenv("GOOGLE_CREDENTIALS", goodCredentials) //used
					// bad
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GCLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCredsPath) // needs to be a path
				},
				Config:                   testAccSdkProvider_credentialsInEnvsOnly(context),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 2nd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					// good
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", goodCredentials) //used
					// bad
					t.Setenv("GCLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCredsPath) // needs to be a path
				},
				Config:                   testAccSdkProvider_credentialsInEnvsOnly(context),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 3rd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", "")
					// good
					t.Setenv("GCLOUD_KEYFILE_JSON", goodCredentials) //used
					// bad
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCredsPath) // needs to be a path
				},
				Config:                   testAccSdkProvider_credentialsInEnvsOnly(context),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				// GOOGLE_APPLICATION_CREDENTIALS is used 4th
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", "")
					t.Setenv("GCLOUD_KEYFILE_JSON", "")
					// bad
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCredsPath) // used, needs to be a path
				},
				ExpectError:              regexp.MustCompile(fmt.Sprintf("%s: no such file", badCredsPath)), // Errors when tries to use GOOGLE_APPLICATION_CREDENTIALS
				Config:                   testAccSdkProvider_credentialsInEnvsOnly(context),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				// Make last step have credentials to enable deleting the resource
				PreConfig: func() {
					t.Setenv("GOOGLE_CREDENTIALS", goodCredentials)
				},
				Config:                   "// Empty config, to force deletion of resources using credentials set above",
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
		},
	})
}

func testAccSdkProvider_credentials_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	credentials := envvar.GetTestCredsFromEnv()

	// ensure all credentials env vars set
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, credentials)
	}

	context := map[string]interface{}{
		"credentials":   "", // empty string used
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

// testAccSdkProvider_credentialsInProviderBlock allows setting the credentials argument in a provider block.
// This function uses google_service_account because it is implemented with the SDK and hows how creds are handled
// in the SDK provider config.
func testAccSdkProvider_credentialsInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	credentials = "%{credentials}"
}

resource "google_service_account" "default" {
  account_id   = "%{resource_name}"
  display_name = "Testing, provisioned by testAccSdkProvider_credentialsInProviderBlock_provisionSdkResource"
}
`, context)
}

// testAccSdkProvider_credentialsInEnvsOnly allows testing when the credentials argument
// is only supplied via ENVs
func testAccSdkProvider_credentialsInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "default" {
  account_id   = "%{resource_name}"
  display_name = "Testing, provisioned by testAccSdkProvider_credentialsInEnvsOnly"
}
`, context)
}
