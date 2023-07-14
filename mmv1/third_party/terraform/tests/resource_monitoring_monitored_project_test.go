package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccMonitoringMonitoredProject_monitoringMonitoredProjectBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.73.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccMonitoringMonitoredProject_monitoringMonitoredProjectBasic(context),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccMonitoringMonitoredProject_monitoringMonitoredProjectBasic(context),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_monitoring_monitored_project.primary",
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"metrics_scope"},
			},
		},
	})
}

func testAccMonitoringMonitoredProject_monitoringMonitoredProjectBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_id}"
  name          = google_project.basic.project_id
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}-display"
  org_id     = "%{org_id}"
}
`, context)
}

func TestAccMonitoringMonitoredProject_monitoringMonitoredProjectLongName(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.72.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccMonitoringMonitoredProject_monitoringMonitoredProjectLongName(context),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccMonitoringMonitoredProject_monitoringMonitoredProjectLongName(context),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_monitoring_monitored_project.primary",
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"metrics_scope"},
			},
		},
	})
}

func testAccMonitoringMonitoredProject_monitoringMonitoredProjectLongName(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "locations/global/metricsScopes/%{project_id}"
  name          = "locations/global/metricsScopes/%{project_id}/projects/${google_project.basic.project_id}"
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}-display"
  org_id     = "%{org_id}"
}
`, context)
}
