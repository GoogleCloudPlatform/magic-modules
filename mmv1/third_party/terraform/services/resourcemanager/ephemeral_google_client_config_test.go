package resourcemanager_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccEphemeralGoogleClientConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Note: For ephemeral resources, we can't directly check attributes
				// since they don't persist in state. Instead, we verify the configuration
				// compiles and runs without error.
				Config: testAccCheckEphemeralGoogleClientConfig_basic,
			},
		},
	})
}

func TestAccEphemeralGoogleClientConfig_omitLocation(t *testing.T) {
	t.Setenv("GOOGLE_REGION", "")
	t.Setenv("GOOGLE_ZONE", "")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Note: For ephemeral resources, we can't directly check attributes
				// since they don't persist in state. Instead, we verify the configuration
				// compiles and runs without error.
				Config: testAccCheckEphemeralGoogleClientConfig_basic,
			},
		},
	})
}

func TestAccEphemeralGoogleClientConfig_invalidCredentials(t *testing.T) {
	badCreds := acctest.GenerateFakeCredentialsJson("test")
	t.Setenv("GOOGLE_CREDENTIALS", badCreds)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckEphemeralGoogleClientConfig_basic,
				ExpectError: regexp.MustCompile("Error setting access_token"),
			},
		},
	})
}

func TestAccEphemeralGoogleClientConfig_usedInProvider(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEphemeralGoogleClientConfig_usedInProvider,
				Check: resource.ComposeTestCheckFunc(
					// Verify that the ephemeral resource can be used to configure a provider
					// and that the provider works correctly
					resource.TestCheckResourceAttrSet("data.google_client_openid_userinfo.me", "email"),
				),
			},
		},
	})
}

const testAccCheckEphemeralGoogleClientConfig_basic = `
provider "google" {
  default_labels = {
    default_key = "default_value"
  }
}

ephemeral "google_client_config" "current" { }
`

const testAccCheckEphemeralGoogleClientConfig_usedInProvider = `
provider "google" {
  default_labels = {
    default_key = "default_value"
  }
}

ephemeral "google_client_config" "current" { }

provider "google" {
  alias        = "ephemeral_configured"
  access_token = ephemeral.google_client_config.current.access_token
  project      = ephemeral.google_client_config.current.project
  region       = ephemeral.google_client_config.current.region
  zone         = ephemeral.google_client_config.current.zone
}

data "google_client_openid_userinfo" "me" {
  provider = google.ephemeral_configured
}
`
