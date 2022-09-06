package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleHealthcareProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_healthcare_project_service_account.gcs_account"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleHealthcareProjectServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_address"),
				),
			},
		},
	})
}

const testAccCheckGoogleHealthcareProjectServiceAccount_basic = `
data "google_healthcare_project_service_account" "gcs_account" {
}
`
