package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleClientConfig_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_client_config.current"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "zone"),
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleClientConfig_omitLocation(t *testing.T) {
	t.Setenv("GOOGLE_REGION", "")
	t.Setenv("GOOGLE_ZONE", "")

	resourceName := "data.google_client_config.current"

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

// Test checks how the data source behaves when credentials are set in the provider block,
// including when the credentials are valid and not valid
func TestAccDataSourceGoogleClientConfig_credentialsInProviderConfig(t *testing.T) {
	// t.Parallel() - Cannot use as we change ENVs

	goodCreds := envvar.GetTestCredsFromEnv()
	t.Setenv("GOOGLE_CREDENTIALS", "")

	badCreds := acctest.GenerateEscapedFakeCredentialsJson("test") // Need to double escape quotes

	resourceName := "data.google_client_config.current"

	acctest.VcrTest(t, resource.TestCase{
		// PreCheck cannot be set as usual, as test has GOOGLE_CREDENTIALS unset
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_credentialsInProviderConfig(goodCreds),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
			{
				Config: testAccCheckGoogleClientConfig_credentialsInProviderConfig(badCreds),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`

func testAccCheckGoogleClientConfig_credentialsInProviderConfig(creds string) string {
	return fmt.Sprintf(`
provider "google" {
  credentials = "%s"
}

data "google_client_config" "current" {}`, creds)
}
