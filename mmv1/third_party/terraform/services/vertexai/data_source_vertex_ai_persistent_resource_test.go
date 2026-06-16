package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/vertexai"
)

func TestAccDataSourceVertexAIPersistentResource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIPersistentResourceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVertexAIPersistentResource_basic(context),
				Check: acctest.CheckDataSourceStateMatchesResourceState(
					"data.google_vertex_ai_persistent_resource.foo",
					"google_vertex_ai_persistent_resource.persistent_resource",
				),
			},
		},
	})
}

func testAccDataSourceVertexAIPersistentResource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_persistent_resource" "persistent_resource" {
  name = "tf-test-persistent-resource%{random_suffix}"
  location = "us-central1"
  resource_pools {
    id = "pool-1"
    machine_spec {
      machine_type = "n1-standard-4"
    }
    replica_count = 1
    disk_spec {
      boot_disk_type = "pd-ssd"
      boot_disk_size_gb = 100
    }
  }
}

data "google_vertex_ai_persistent_resource" "foo" {
  name = google_vertex_ai_persistent_resource.persistent_resource.name
  location = google_vertex_ai_persistent_resource.persistent_resource.location
  project = google_vertex_ai_persistent_resource.persistent_resource.project
}
`, context)
}
