package apphub_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestDataSourceApphubDiscoveredWorkload_basic(t *testing.T) {
    t.Parallel()
    
    context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

    acctest.VcrTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.AccTestPreCheck(t) },
        ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
        Steps: []resource.TestStep{
            {
                Config: testDataSourceApphubDiscoveredWorkload_basic(context),
            },
        },
    })
}

func testDataSourceApphubDiscoveredWorkload_basic(context map[string]interface{}) string {
    return acctest.Nprintf(`
data "google_project" "host_project" {}

resource "google_project_service" "apphub" {
  project = data.google_project.host_project.project_id
  service = "apphub.googleapis.com"
  disable_on_destroy = false
}

data "google_apphub_discovered_workload" "catalog-workload" {
  provider = google
  location = "us-east1"
  workload_uri = "my-workload-uri"
}
`, context)
}

