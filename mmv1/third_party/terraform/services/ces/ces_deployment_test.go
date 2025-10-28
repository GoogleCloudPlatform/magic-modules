package ces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck" // Add this import
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCESDeployment_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESDeployment_cesDeploymentBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_deployment.my-deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "app_version", "location"},
			},
			{
				Config: testAccCESDeployment_cesDeploymentBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_deployment.my-deployment", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_deployment.my-deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "app_version", "location"},
			},
		},
	})
}

func testAccCESDeployment_cesDeploymentBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_deployment" "my-deployment" {
    location     = "us"
    display_name = "tf-test-my-deployment%{random_suffix}"
    app          = google_ces_app.my-app.name
    app_version  = "projects/example-project/locations/us/apps/example-app/versions/example-version"
    channel_profile {
        channel_type = "API"
        disable_barge_in_control = true
        disable_dtmf = true
        persona_property {
            persona = "CHATTY"
        }
        profile_id = "temp_profile_id"
        web_widget_config {
            modality = "CHAT_AND_VOICE"
            theme = "DARK"
            web_widget_title = "temp_webwidget_title"
        }
    }
}
`, context)
}

func testAccCESDeployment_cesDeploymentBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_deployment" "my-deployment" {
    location     = "us"
    display_name = "tf-test-my-deployment%{random_suffix}"
    app          = google_ces_app.my-app.name
    app_version  = "projects/example-project/locations/us/apps/example-app/versions/example-version"
    channel_profile {
        channel_type = "WEB_UI"
        disable_barge_in_control = true
        disable_dtmf = true
        persona_property {
            persona = "CONCISE"
        }
        profile_id = "temp_profile_id"
        web_widget_config {
            modality = "CHAT_ONLY"
            theme = "LIGHT"
            web_widget_title = "temp_webwidget_title"
        }
    }
}
`, context)
}
