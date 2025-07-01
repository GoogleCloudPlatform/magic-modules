package modelarmorglobal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccModelArmorGlobalFloorsetting_basic(t *testing.T) {

	basicContext := map[string]interface{}{
		"location": "global",
		"parent":   "projects/modelarmor-api-test",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccModelArmorGlobalFloorsetting_basicContext(basicContext),
			},
			{
				ResourceName:      "google_model_armor_floorsetting.floorsetting-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccModelArmorGlobalFloorsetting_basicContext(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_model_armor_floorsetting" "floorsetting-basic" {
  parent      = "%{parent}"
  location    = "%{location}"

  filter_config {

  }
}
`, context)
}

func TestAccModelArmorGlobalFloorSetting_update(t *testing.T) {

	context := map[string]interface{}{
		"location": "global",
		"parent":   "projects/modelarmor-api-test",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccModelArmorGlobalFloorsetting_initial(context),
			},
			{
				ResourceName:            "google_model_armor_floorsetting.test-resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent"},
			},
			{
				Config: testAccModelArmorGlobalFloorsetting_updated(context),
			},
			{
				ResourceName:            "google_model_armor_floorsetting.test-resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent"},
			},
		},
	})
}

func testAccModelArmorGlobalFloorsetting_initial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_model_armor_floorsetting" "test-resource" {
  location    = "%{location}"
  parent      = "%{parent}"

  filter_config {
    rai_settings {
      rai_filters {
        filter_type      = "DANGEROUS"
        confidence_level = "LOW_AND_ABOVE"
      }
    }
    sdp_settings {
      basic_config {
        filter_enforcement = "ENABLED"
      }
    }
    pi_and_jailbreak_filter_settings {
      filter_enforcement = "ENABLED"
      confidence_level   = "MEDIUM_AND_ABOVE"
    }
    malicious_uri_filter_settings {
      filter_enforcement = "ENABLED"
    }
  }

  enable_floor_setting_enforcement = true
  
  integrated_services =  [ "AI_PLATFORM" ]

  ai_platform_floor_setting {
    inspect_only            = true
    enable_cloud_logging    = true
  }
  
  floor_setting_metadata {
    multi_language_detection {
      enable_multi_language_detection = true
    }
  }
}
`, context)
}

func testAccModelArmorGlobalFloorsetting_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_model_armor_floorsetting" "test-resource" {
  location    = "%{location}"
  parent      = "%{parent}"

  filter_config {
    rai_settings {
      rai_filters {
        filter_type      = "SEXUALLY_EXPLICIT"
        confidence_level = "HIGH"
      }
    }
    sdp_settings {
      advanced_config {
        inspect_template    = "projects/modelarmor-api-test/locations/global/inspectTemplates/modelarmor-tf-test"
        deidentify_template = "projects/modelarmor-api-test/locations/us-central1/deidentifyTemplates/modelarmor-tf-test"
      }
    }
    pi_and_jailbreak_filter_settings {
      filter_enforcement = "ENABLED"
      confidence_level   = "MEDIUM_AND_ABOVE"
    }
    malicious_uri_filter_settings {
      filter_enforcement = "ENABLED"
    }
  }

  enable_floor_setting_enforcement = true

  ai_platform_floor_setting {
    inspect_only            = false
    enable_cloud_logging    = false
  }
  
  floor_setting_metadata {
    multi_language_detection {
      enable_multi_language_detection = false
    }
  }
}
`, context)
}
