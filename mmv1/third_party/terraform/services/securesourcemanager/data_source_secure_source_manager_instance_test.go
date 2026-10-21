package securesourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/securesourcemanager"
)

func TestAccSecureSourceManagerInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecureSourceManagerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecureSourceManagerInstanceDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// deletion_policy is a client-side virtual field not returned by the API; the test resource sets it to "DELETE" for cleanup.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_secure_source_manager_instance.default",
						"google_secure_source_manager_instance.default",
						[]string{"deletion_policy"},
					),
					// Verify omitted optional location and project fall back to provider defaults.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_secure_source_manager_instance.without_optional_scope",
						"google_secure_source_manager_instance.default",
						[]string{"deletion_policy"},
					),
				),
			},
		},
	})
}

func testAccSecureSourceManagerInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "default" {
  location        = "us-central1"
  instance_id     = "tf-test-my-instance%{random_suffix}"
  deletion_policy = "DELETE"

  labels = {
    "foo" = "bar"
  }
}

data "google_secure_source_manager_instance" "default" {
  instance_id = google_secure_source_manager_instance.default.instance_id
  location    = google_secure_source_manager_instance.default.location
  project     = google_secure_source_manager_instance.default.project
}

data "google_secure_source_manager_instance" "without_optional_scope" {
  instance_id = google_secure_source_manager_instance.default.instance_id
}
`, context)
}
