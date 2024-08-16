package resourcemanager_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
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

// Test checks how the data source behaves when invalid credentials are used
// This test sets them as an ENV
func TestAccDataSourceGoogleClientConfig_invalidCredentials(t *testing.T) {
	// t.Parallel() - Cannot use as we change ENVs

	badCreds := acctest.GenerateFakeCredentialsJson("test")
	t.Setenv("GOOGLE_CREDENTIALS", badCreds)

	acctest.VcrTest(t, resource.TestCase{
		// PreCheck cannot be set as usual, as test has GOOGLE_CREDENTIALS unset
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleClientConfig_basic,
				ExpectError: regexp.MustCompile("Error setting access_token"),
			},
		},
	})
}

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`
