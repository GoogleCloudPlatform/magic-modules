package servicedirectory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/servicedirectory"
)

func TestAccDataSourceServiceDirectoryNamespaces_basic(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	namespaceId := fmt.Sprintf("tf-test-namespaces-%s", suffix)
	// The label value is alphanumeric, like the labels.owner=sd example in the API reference.
	labelValue := fmt.Sprintf("tftest%s", suffix)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckServiceDirectoryNamespaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceDirectoryNamespacesConfig(namespaceId, labelValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_service_directory_namespaces.one", "namespaces.#", "1"),
					resource.TestCheckResourceAttr("data.google_service_directory_namespaces.one", "namespaces.0.namespace_id", namespaceId),
					resource.TestCheckResourceAttr("data.google_service_directory_namespaces.one", "namespaces.0.labels.testid", labelValue),
					resource.TestCheckResourceAttrPair("data.google_service_directory_namespaces.one", "namespaces.0.name", "google_service_directory_namespace.foobar", "name"),
					// The shared test project holds other namespaces, so the unfiltered read has no count check.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_service_directory_namespaces.all", "namespaces.*", map[string]string{"namespace_id": namespaceId}),
				),
			},
		},
	})
}

func testAccDataSourceServiceDirectoryNamespacesConfig(namespaceId, labelValue string) string {
	return fmt.Sprintf(`
resource "google_service_directory_namespace" "foobar" {
  namespace_id = "%s"
  location     = "us-central1"

  labels = {
    testid = "%s"
  }
}

data "google_service_directory_namespaces" "one" {
  location = google_service_directory_namespace.foobar.location
  filter   = "labels.testid=${google_service_directory_namespace.foobar.labels.testid}"
}

data "google_service_directory_namespaces" "all" {
  location = google_service_directory_namespace.foobar.location
}
`, namespaceId, labelValue)
}
