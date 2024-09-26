// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccSdkProvider_request_timeout is a series of acc tests asserting how the SDK provider handles request_timeout arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_request_timeout
func TestAccSdkProvider_request_timeout(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"request_timeout can be set in config in different formats": testAccSdkProvider_request_timeout_setInConfig,
		//no ENVs to test

		// Schema-level validation
		"when request_timeout is set to an empty string in the config the value fails validation, as it is not a duration": testAccSdkProvider_request_timeout_emptyStringValidation,

		// Usage
		// We cannot test the impact of this field in an acc test
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

func testAccSdkProvider_request_timeout_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	providerTimeout1 := "3m0s"
	providerTimeout2 := "3m"
	expectedValue := "3m0s"

	context1 := map[string]interface{}{
		"request_timeout": providerTimeout1,
	}
	context2 := map[string]interface{}{
		"request_timeout": providerTimeout2,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_request_timeout_inProviderBlock(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_timeout", expectedValue),
				),
			},
			{
				Config: testAccSdkProvider_request_timeout_inProviderBlock(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_timeout", expectedValue),
				),
			},
		},
	})
}

func testAccSdkProvider_request_timeout_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"request_timeout": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_request_timeout_inProviderBlock(context),
				ExpectError: regexp.MustCompile("invalid duration"),
			},
		},
	})
}

// testAccSdkProvider_request_timeout_inProviderBlock allows setting the request_timeout argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_request_timeout_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	request_timeout = "%{request_timeout}"
}

data "google_provider_config_sdk" "default" {}
`, context)
}
