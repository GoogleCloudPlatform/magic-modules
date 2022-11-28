package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleFirebaseAppleApp(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":   getTestProjectFromEnv(),
		"bundle_id":    "apple.app.12345",
		"display_name": "Display Name AppleApp DataSource",
		"app_store_id": 12345,
		"team_id":      1234567890,
	}

	resourceName := "data.google_firebase_apple_app.my_app"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleFirebaseAppleApp(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState(resourceName,
						"google_firebase_apple_app.my_app"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleFirebaseAppleApp(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firebase_apple_app" "my_app" {
  provider = google-beta
  project = "%{project_id}"
  bundle_id = "%{bundle_id}"
  display_name = "%{display_name}"
  app_store_id = "%{app_store_id}"
  team_id = "%{team_id}"
}

data "google_firebase_apple_app" "my_app" {
  app_id = google_firebase_apple_app.my_app.app_id
}
`, context)
}
