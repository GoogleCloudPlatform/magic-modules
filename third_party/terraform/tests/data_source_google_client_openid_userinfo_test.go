package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleClientOpenIDUserinfo_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientOpenIDUserinfo_basic,
				// While this _should_ pass, we have no way to provide custom default scopes to the test runner
				// so we won't actually have the necessary scopes to ensure this is working.
				ExpectError: regexp.MustCompile("Invalid Credentials"),
			},
		},
	})
}

const testAccCheckGoogleClientOpenIDUserinfo_basic = `
data "google_client_openid_userinfo" "me" { }
`
