package composer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerUserWorkloadsSecret_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"env_name":    fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t)),
		"secret_name": fmt.Sprintf("%s-%d", testComposerUserWorkloadsSecretPrefix, acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComposerUserWorkloadsSecret_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_composer_user_workloads_secret.test",
						"google_composer_user_workloads_secret.test"),
				),
			},
		},
	})
}

func testAccDataSourceComposerUserWorkloadsSecret_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_composer_environment" "test" {
	name   = "%{env_name}"
	config {
		software_config {
			image_version = "composer-3-airflow-2"
		}
	}
}
resource "google_composer_user_workloads_secret" "test" {
  environment = google_composer_environment.test.name
  name = "%{secret_name}"
  data = {
    username: base64encode("username"),
    password: base64encode("password"),
  }
}
data "google_composer_user_workloads_secret" "test" {
	name        = google_composer_user_workloads_secret.test.name
	environment = google_composer_environment.test.name
}
`, context)
}
