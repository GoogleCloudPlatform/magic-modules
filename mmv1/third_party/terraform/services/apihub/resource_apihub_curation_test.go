package apihub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccApihubCuration_apihubCurationBasic_Update(t *testing.T) {
	// This is added for reference, but the test needs to be skipped as it needs API hub instance as a prerequisite
	// But the support for that resources is not yet complete.
	t.Skip()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccApihubCuration_apihubCuration_basic(context),
			},
			{
				ResourceName:            "google_apihub_curation.apihub_curation_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"curation_id", "location"},
			},
			{
				Config: testAccApihubCuration_apihubCuration_update(context),
			},
			{
				ResourceName:            "google_apihub_curation.apihub_curation_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"curation_id", "location"},
			},
		},
	})
}

func testAccApihubCuration_apihubCuration_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apihub_curation" "apihub_curation_basic" {
  location = "us-central1"
  curation_id = "test%{random_suffix}"
  display_name = "Test Curation"
  description = "This is a sample curation resource managed by Terraform."
  endpoint {
    application_integration_endpoint_details {
      trigger_id = "api_trigger/curation_API_1"
      uri = "https://integrations.googleapis.com/v1/projects/1082615593856/locations/us-central1/integrations/curation:execute"
    }
  }

}


`, context)
}

func testAccApihubCuration_apihubCuration_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apihub_curation" "apihub_curation_basic" {
  location = "us-central1"
  curation_id = "test%{random_suffix}"
  display_name = "Test Curation Updated"
  description = "This is a sample updated curation resource managed by Terraform."
  endpoint {
    application_integration_endpoint_details {
      trigger_id = "api_trigger/curation_API_1"
      uri = "https://integrations.googleapis.com/v1/projects/1082615593856/locations/us-central1/integrations/curation:execute"
    }
  }

}


`, context)
}
