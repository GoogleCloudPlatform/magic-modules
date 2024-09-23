package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFwProvider_impersonate_service_account_delegates(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		//     There are no environment variables for this field
		"impersonate_service_account_delegates can be set in config": testAccSdkProvider_impersonate_service_account_delegates_setInConfig,

		// Schema-level validation
		"when impersonate_service_account_delegates is set to an empty list in the config the value IS ignored": testAccSdkProvider_impersonate_service_account_delegates_emptyListUsage,

		// Usage
		"impersonate_service_account_delegates controls which service account is used for actions": testAccSdkProvider_impersonate_service_account_delegates_usage,
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

func testAccSdkProvider_impersonate_service_account_delegates_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	delegates := []string{
		"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
		"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
	}
	delegatesString := fmt.Sprintf(`["%s","%s"]`, delegates[0], delegates[1])

	// There are no ENVs for this provider argument

	context := map[string]interface{}{
		"random_suffix":                         acctest.RandString(t, 10),
		"impersonate_service_account_delegates": delegatesString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account_delegates.#", "2"),
				),
			},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_delegates_emptyListUsage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"random_suffix":                         acctest.RandString(t, 10),
		"impersonate_service_account_delegates": "[]",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account_delegates.#", "0"),
				),
			},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_delegates_usage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context),
			},
			{
				// This needs to be split into a second step as impersonate_service_account_delegates does
				// not tolerate unknown values
				Config:      testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_2(context),
				ExpectError: regexp.MustCompile("Error creating Topic: googleapi: Error 403: User not authorized"),
			},
		},
	})
}

// testAccSdkProvider_impersonate_service_account_delegates_testProvisioning allows setting the impersonate_service_account_delegates argument in a provider block
// and testing its impact on provisioning a resource
func testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account_delegates = %{impersonate_service_account_delegates}
}

data "google_provider_config_sdk" "default" {}

resource "google_pubsub_topic" "example" {
  name = "tf-test-%{random_suffix}"
}
`, context)
}

func testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context map[string]interface{}) string {
	return acctest.Nprintf(`
// This will succeed due to the Terraform identity having necessary permissions
resource "google_pubsub_topic" "ok" {
  name = "tf-test-%{random_suffix}-ok"
}

//  Create a first service account and ensure the Terraform identity can make tokens for it
resource "google_service_account" "default_1" {
  account_id   = "tf-test-%{random_suffix}-1"
  display_name = "Acceptance test impersonated service account"
}

data "google_client_openid_userinfo" "me" {
}

resource "google_service_account_iam_member" "token_1" {
  service_account_id = google_service_account.default_1.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${data.google_client_openid_userinfo.me.email}"
}

//  Create a second service account and ensure the first service account can make tokens for it
resource "google_service_account" "default_2" {
  account_id   = "tf-test-%{random_suffix}-2"
  display_name = "Acceptance test impersonated service account"
}

resource "google_service_account_iam_member" "token_2" {
  service_account_id = google_service_account.default_2.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.default_1.email}"
}
`, context)
}

func testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_2(context map[string]interface{}) string {
	return testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context) + acctest.Nprintf(`

// Use impersonate_service_account_delegates
provider "google" {
  alias = "impersonation"
  impersonate_service_account = google_service_account.default_2.email
  impersonate_service_account_delegates = [
    google_service_account.default_1.email,
    google_service_account.default_2.email,
  ]
}

// This will fail due to the impersonated service account not having any permissions
resource "google_pubsub_topic" "fail" {
  provider = google.impersonation
  name = "tf-test-%{random_suffix}-fail"
}
`, context)
}
