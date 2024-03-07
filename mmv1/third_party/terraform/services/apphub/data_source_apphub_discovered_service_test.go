package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestDataSourceApphubDiscoveredService_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApphubDiscoveredService_basic(context),
			},
		},
	})
}

func testDataSourceApphubDiscoveredService_basic(context map[string]interface{}) string {
	return acctest.Nprintf(
		`
    data "google_project" "host_project" {}

    resource "google_project_service" "apphub" {
      project = data.google_project.host_project.project_id
      service = "apphub.googleapis.com"
      disable_on_destroy = false
    }
    
    data "google_apphub_discovered_service" "catalog-service" {
      provider = google
      location = "us-east1"
      service_uri = "my-service-uri"
    }

`, context)
}
