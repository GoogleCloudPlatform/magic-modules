package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSdkProvider_project(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"config takes precedence over environment variables":                           testAccSdkProvider_project_configPrecedenceOverEnvironmentVariables,
		"when project is unset in the config, environment variables are used":          testAccSdkProvider_project_precedenceOrderEnvironmentVariables,
		"when project is set to an empty string in the config the value isn't ignored": testAccSdkProvider_project_emptyStringValidation,
		"when project is unknown in the config, environment variables are used":        testAccSdkProvider_project_unknownHandling,
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

func testAccSdkProvider_project_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	project := envvar.GetTestProjectFromEnv()

	// set all possible project env vars to other value; show they aren't used
	for _, v := range envvar.ProjectEnvVars {
		t.Setenv(v, "foobar")
	}

	context := map[string]interface{}{
		"project": project,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_projectInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", project),
				),
			},
		},
	})
}

func testAccSdkProvider_project_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	/*
		These are all the ENVs for project, and they are in order of precedence.
		GOOGLE_PROJECT
		GOOGLE_CLOUD_PROJECT
		GCLOUD_PROJECT
		CLOUDSDK_CORE_PROJECT
	*/

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Unset all ENVs for project
				PreConfig: func() {
					for _, v := range envvar.ProjectEnvVars {
						t.Setenv(v, "")
					}
				},
				Config: testAccSdkProvider_projectInEnvsOnly(),
				Check: resource.ComposeTestCheckFunc(
					// Differing behavior between SDK and PF; the attribute is found here.
					// This reflects the different type systems used in the SDKv2 and the plugin-framework
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", ""),
				),
			},
			{
				// GOOGLE_PROJECT is used 1st if set
				PreConfig: func() {
					t.Setenv("GOOGLE_PROJECT", "GOOGLE_PROJECT") // used
					t.Setenv("GOOGLE_CLOUD_PROJECT", "GOOGLE_CLOUD_PROJECT")
					t.Setenv("GCLOUD_PROJECT", "GCLOUD_PROJECT")
					t.Setenv("CLOUDSDK_CORE_PROJECT", "CLOUDSDK_CORE_PROJECT")
				},
				Config: testAccSdkProvider_projectInEnvsOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", "GOOGLE_PROJECT"),
				),
			},
			{
				// GOOGLE_CLOUD_PROJECT is used 2nd if set
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_PROJECT", "")
					// set
					t.Setenv("GOOGLE_CLOUD_PROJECT", "GOOGLE_CLOUD_PROJECT") //used
					t.Setenv("GCLOUD_PROJECT", "GOOGLE_CLOUD_PROJECT")
					t.Setenv("CLOUDSDK_CORE_PROJECT", "CLOUDSDK_CORE_PROJECT")
				},
				Config: testAccSdkProvider_projectInEnvsOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", "GOOGLE_CLOUD_PROJECT"),
				),
			},
			{
				// GCLOUD_PROJECT is used 3rd if set
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_PROJECT", "")
					t.Setenv("GOOGLE_CLOUD_PROJECT", "")
					// set
					t.Setenv("GCLOUD_PROJECT", "GCLOUD_PROJECT") // used
					t.Setenv("CLOUDSDK_CORE_PROJECT", "CLOUDSDK_CORE_PROJECT")
				},
				Config: testAccSdkProvider_projectInEnvsOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", "GCLOUD_PROJECT"),
				),
			},
			{
				// CLOUDSDK_CORE_PROJECT is used 4th if set
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_PROJECT", "")
					t.Setenv("GOOGLE_CLOUD_PROJECT", "")
					t.Setenv("GCLOUD_PROJECT", "")
					// set
					t.Setenv("CLOUDSDK_CORE_PROJECT", "CLOUDSDK_CORE_PROJECT")
				},
				Config: testAccSdkProvider_projectInEnvsOnly(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// When not explicitly waiting for the unknown value to resolve the project could be the ENV or could be from the provisioned project
					// In practice it might be that this data source would be refreshed during an apply (?)
					resource.TestMatchResourceAttr("data.google_provider_config_sdk.default", "project", regexp.MustCompile(fmt.Sprintf("(tf-test-[0-9a-z]{16}|%s)", project))),

					// When explicitly made to wait, returns the provisioned project's id
					resource.TestMatchResourceAttr("data.google_provider_config_sdk.wait", "project", regexp.MustCompile("tf-test-[0-9a-z]{16}")),
				),
			},
		},
	})
}

func testAccSdkProvider_project_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"project": "",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_projectInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccSdkProvider_project_unknownHandling(t *testing.T) {

	project := envvar.GetTestProjectFromEnv()
	context := map[string]interface{}{
		"org_id":             envvar.GetTestOrgFromEnv(t),
		"billing_account_id": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_projectUnknownHandling(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Matches ENV instead of the project id output from google_project resource
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "project", project),
				),
			},
			{
				// Unset all ENVs for project
				PreConfig: func() {
					for _, v := range envvar.ProjectEnvVars {
						t.Setenv(v, "")
					}
				},
				Config: testAccSdkProvider_projectUnknownHandling(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Matches project id output from google_project resource
					resource.TestMatchResourceAttr("data.google_provider_config_sdk.default", "project", regexp.MustCompile("tf-test-[0-9a-z]{16}")),
				),
			},
		},
	})
}

// testAccSdkProvider_projectInProviderBlock allows setting the project argument in a provider block.
func testAccSdkProvider_projectInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	project = "%{project}"
}
data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_projectInEnvsOnly allows testing when the project argument
// is only supplied via ENVs
func testAccSdkProvider_projectInEnvsOnly() string {
	return `
data "google_provider_config_sdk" "default" {}
`
}

// testAccSdkProvider_projectUnknownHandling is specifically for testing how an unknown value is used.
func testAccSdkProvider_projectUnknownHandling(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "random_id" "project_name" {
  byte_length = 8
}

provider "google" {
	alias = "alternate"
}

resource "google_project" "project" {
  provider        = google.alternate
  name            = "Test Acc Project"
  project_id      = "tf-test-${random_id.project_name.hex}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account_id}"
  deletion_policy = "DELETE"
}

// Note that this is the unaliased provider, and is used in the data source below
provider "google" {
	project = google_project.project.project_id
}

data "google_provider_config_sdk" "default" {
}
`, context)
}
