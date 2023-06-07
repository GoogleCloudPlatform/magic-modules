package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleCoreBillingProjectInfo_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": acctest.GetTestBillingAccountFromEnv(t),
		"org_id":          acctest.GetTestOrgFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCoreBillingProjectInfoDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCoreBillingProjectInfo_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_billing_project_info.default", "google_billing_project_info.default"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCoreBillingProjectInfo_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "default" {
  project_id = "tf-test%{random_suffix}"
  name       = "tf-test%{random_suffix}"
  org_id     = "%{org_id}"
  lifecycle {
    ignore_changes = [billing_account]
  }
}

resource "google_billing_project_info" "default" {
  project         = google_project.default.project_id
  billing_account = "%{billing_account}"
}

data "google_billing_project_info" "default" {
  project    = google_project.default.project_id
  depends_on = [google_billing_project_info.default]
}
`, context)
}
