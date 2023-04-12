package google_test

import (
	google "internal/terraform-provider-google"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessApprovalOrganizationServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": google.GetTestOrgFromEnv(t),
	}

	resourceName := "data.google_access_approval_organization_service_account.aa_account"

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context map[string]interface{}) string {
	return google.Nprintf(`
data "google_access_approval_organization_service_account" "aa_account" {
  organization_id = "%{org_id}"
}
`, context)
}
