package observability_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccObservabilityProjectSettings_datasource(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"location":        "us",
	}
	dataResourceName := "data.google_observability_project_settings.settings"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityProjectSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "name"),
					resource.TestCheckResourceAttrSet(dataResourceName, "service_account_id"),
					resource.TestCheckResourceAttr(dataResourceName, "location", context["location"].(string)),
					resource.TestCheckResourceAttr(dataResourceName, "project", context["project_name"].(string)),
				),
			},
		},
	})
}
func testAccObservabilityProjectSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "default" {
		project_id      = "%{project_name}"
		name            = "%{project_name}"
		org_id          = "%{org_id}"
		billing_account = "%{billing_account}"
		deletion_policy = "DELETE"
	}
	resource "google_project_service" "observability_service" {
		project = google_project.default.project_id
		service = "observability.googleapis.com"
		disable_on_destroy = false
	}
	resource "time_sleep" "wait_for_project" {
		create_duration = "60s"
		depends_on = [google_project_service.observability_service]
	}

	data "google_observability_project_settings" "settings" {
		project  = google_project.default.project_id
		location = "%{location}"
		depends_on = [time_sleep.wait_for_project]
	}
`, context)
}
