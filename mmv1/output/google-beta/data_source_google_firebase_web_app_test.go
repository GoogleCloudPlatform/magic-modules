package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleFirebaseWebApp(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":   GetTestProjectFromEnv(),
		"display_name": "tf_test Display Name WebApp DataSource",
	}

	resourceName := "data.google_firebase_web_app.my_app"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleFirebaseWebApp(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						resourceName,
						"google_firebase_web_app.my_app",
						map[string]struct{}{
							"deletion_policy": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceGoogleFirebaseWebApp(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firebase_web_app" "my_app" {
  project = "%{project_id}"
  display_name = "%{display_name}"
  deletion_policy = "DELETE"
}

data "google_firebase_web_app" "my_app" {
  app_id = google_firebase_web_app.my_app.app_id
}

data "google_firebase_web_app" "my_app_project" {
  project = "%{project_id}"
  app_id = google_firebase_web_app.my_app.app_id
}
`, context)
}
