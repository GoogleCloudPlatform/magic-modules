package ces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck" // Add this import
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/ces"
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
        noise_suppression_level = "NOISE_SUPPRESSION_LEVEL_UNSPECIFIED"
        persona_property {
            persona = "CHATTY"
        }
        profile_id = "temp_profile_id"
        web_widget_config {
            modality = "CHAT_AND_VOICE"
            theme = "DARK"
            web_widget_title = "temp_webwidget_title"
            security_settings {
                enable_public_access = true
                enable_origin_check = true
                allowed_origins = ["https://example.com", "https://test.com"]
                enable_recaptcha = true
            }
        }
        instagram_config {
            instagram_account_id = "insta-id-1"
        }
        whatsapp_config {
            phone_number = "12345678"
            phone_number_id = "phone-id-1"
            waba_id = "waba-id-1"
        }
    }
    experiment_config {
        version_release {
            state = "STATE_UNSPECIFIED"
            traffic_allocations {
                app_version = "projects/example-project/locations/us/apps/example-app/versions/example-version"
                traffic_percentage = 100
                id = "1"
            }
        }
    }
    instagram_credentials {
        auth_code = "insta-auth-code"
        conversation_profile_id = "insta-profile-id"
    }
    whatsapp_credentials {
        auth_code = "wa-auth-code"
        business_account_id = "wa-business-id"
        conversation_profile_id = "wa-profile-id"
        phone_number = "wa-phone"
        pin = "wa-pin"
        waba_id = "wa-waba-id"
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
        noise_suppression_level = "NOISE_SUPPRESSION_LEVEL_UNSPECIFIED"
        persona_property {
            persona = "CONCISE"
        }
        profile_id = "temp_profile_id"
        web_widget_config {
            modality = "CHAT_ONLY"
            theme = "LIGHT"
            web_widget_title = "temp_webwidget_title"
            security_settings {
                enable_public_access = false
                enable_origin_check = false
                allowed_origins = ["https://updated.com"]
                enable_recaptcha = false
            }
        }
        instagram_config {
            instagram_account_id = "insta-id-2"
        }
        whatsapp_config {
            phone_number = "87654321"
            phone_number_id = "phone-id-2"
            waba_id = "waba-id-2"
        }
    }
    experiment_config {
        version_release {
            state = "STATE_UNSPECIFIED"
            traffic_allocations {
                app_version = "projects/example-project/locations/us/apps/example-app/versions/example-version"
                traffic_percentage = 50
                id = "1"
            }
        }
    }
    instagram_credentials {
        auth_code = "insta-auth-code-updated"
        conversation_profile_id = "insta-profile-id-updated"
    }
    whatsapp_credentials {
        auth_code = "wa-auth-code-updated"
        business_account_id = "wa-business-id-updated"
        conversation_profile_id = "wa-profile-id-updated"
        phone_number = "wa-phone-updated"
        pin = "wa-pin-updated"
        waba_id = "wa-waba-id-updated"
    }
}
`, context)
}
