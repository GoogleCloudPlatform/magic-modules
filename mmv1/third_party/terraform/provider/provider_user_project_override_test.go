package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccSdkProvider_user_project_override is a series of acc tests asserting how the plugin-framework provider handles credentials arguments
// It is PF specific because the HCL used uses a PF-implemented data source
// It is a counterpart to TestAccFwProvider_user_project_override
func TestAccSdkProvider_user_project_override(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                testAccSdkProvider_user_project_override_configPrecedenceOverEnvironmentVariables,
		"when user_project_override is unset in the config, environment variables are used": testAccSdkProvider_user_project_override_precedenceOrderEnvironmentVariables,

		// Schema-level validation
		"when user_project_override is set in the config the value can be a boolean (true/false) or a string (true/false/1/0)": testAccSdkProvider_user_project_override_booleansInConfigOnly,
		"when user_project_override is set via environment variables any of these values can be used: true/false/1/0":          testAccSdkProvider_user_project_override_envStringsAccepted,

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

func testAccSdkProvider_user_project_override_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	override := "true"
	providerOverride := false

	// ensure all possible region env vars set; show they aren't used
	t.Setenv("USER_PROJECT_OVERRIDE", override)

	context := map[string]interface{}{
		"user_project_override": providerOverride,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", fmt.Sprintf("%v", providerOverride)),
				),
			},
		},
	})
}

func testAccSdkProvider_user_project_override_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for region, and they are in order of precedence.
		USER_PROJECT_OVERRIDE
	*/

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "") // unset
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// defaults to false when not set via config or ENVs
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "true")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
		},
	})
}

func testAccSdkProvider_user_project_override_booleansInConfigOnly(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context_true := map[string]interface{}{
		"user_project_override": true,
	}
	context_false := map[string]interface{}{
		"user_project_override": false,
	}

	context_1 := map[string]interface{}{
		"user_project_override": "1",
	}
	context_0 := map[string]interface{}{
		"user_project_override": "0",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context_true),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context_false),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_true),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_false),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_1),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_0),
				// No error expected
			},
		},
	})
}

func testAccSdkProvider_user_project_override_envStringsAccepted(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "true")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "1")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "false")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "0")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
		},
	})
}

// TestAccSdkProvider_user_project_overrideInProviderBlock allows setting the user_project_override argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the plugin-framework
func testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context map[string]interface{}) string {
	v := acctest.Nprintf(`
provider "google" {
	user_project_override = %{user_project_override}
}

data "google_provider_config_sdk" "default" {}

output "user_project_override" {
  value = data.google_provider_config_sdk.default.user_project_override
  sensitive = true
}
`, context)
	return v
}

func testAccSdkProvider_user_project_overrideInProviderBlock_string(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	user_project_override = "%{user_project_override}"
}

data "google_provider_config_sdk" "default" {}

output "user_project_override" {
  value = data.google_provider_config_sdk.default.user_project_override
  sensitive = true
}
`, context)
}

// testAccSdkProvider_user_project_overrideInEnvsOnly allows testing when the user_project_override argument
// is only supplied via ENVs
func testAccSdkProvider_user_project_overrideInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}

output "user_project_override" {
  value = data.google_provider_config_sdk.default.user_project_override
  sensitive = true
}
`, context)
}

