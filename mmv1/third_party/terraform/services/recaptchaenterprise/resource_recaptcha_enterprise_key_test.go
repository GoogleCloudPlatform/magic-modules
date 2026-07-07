package recaptchaenterprise_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccRecaptchaEnterpriseKey_AndroidAppStore(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_AndroidAppStore(context),
			},
			{
				ResourceName:            "google_recaptcha_enterprise_key.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccRecaptchaEnterpriseKey_AndroidAppStoreUpdate(context),
			},
			{
				ResourceName:            "google_recaptcha_enterprise_key.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccRecaptchaEnterpriseKey_IosAppleDeveloperId(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_IosAppleDeveloperId(context),
			},
			{
				ResourceName:            "google_recaptcha_enterprise_key.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "ios_settings.0.apple_developer_id.0.private_key"},
			},
			{
				Config: testAccRecaptchaEnterpriseKey_IosAppleDeveloperIdUpdate(context),
			},
			{
				ResourceName:            "google_recaptcha_enterprise_key.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "ios_settings.0.apple_developer_id.0.private_key"},
			},
		},
	})
}

func testAccRecaptchaEnterpriseKey_AndroidAppStore(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "tf-test-android-%{random_suffix}"
  project      = "%{project_name}"

  android_settings {
    allow_all_package_names                   = true
    allowed_package_names                     = []
    support_non_google_app_store_distribution = true
  }

  labels = {
    test = "android"
  }
}
`, context)
}

func testAccRecaptchaEnterpriseKey_AndroidAppStoreUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "tf-test-android-upd-%{random_suffix}"
  project      = "%{project_name}"

  android_settings {
    allow_all_package_names                   = false
    allowed_package_names                     = ["com.example.app"]
    support_non_google_app_store_distribution = false
  }

  labels = {
    test = "android-update"
  }
}
`, context)
}

func testAccRecaptchaEnterpriseKey_IosAppleDeveloperId(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "tf-test-ios-%{random_suffix}"
  project      = "%{project_name}"

  ios_settings {
    allow_all_bundle_ids = true
    allowed_bundle_ids   = []
    apple_developer_id {
      key_id      = "1234567890"
      private_key = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCg=="
      team_id     = "0987654321"
    }
  }

  labels = {
    test = "ios"
  }
}
`, context)
}

func testAccRecaptchaEnterpriseKey_IosAppleDeveloperIdUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "tf-test-ios-upd-%{random_suffix}"
  project      = "%{project_name}"

  ios_settings {
    allow_all_bundle_ids = false
    allowed_bundle_ids   = ["com.example.ios"]
    apple_developer_id {
      key_id      = "0987654321"
      private_key = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCg=="
      team_id     = "1234567890"
    }
  }

  labels = {
    test = "ios-update"
  }
}
`, context)
}
