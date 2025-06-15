package modelarmor_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccModelArmorTemplate_createAndUpdate(t *testing.T) {
	t.Parallel()

	templateId := "tf-test-" + acctest.RandString(t, 10)

	// Step 1: Create without any SDP settings
	initialContext := map[string]interface{}{
		"location":         "us-central1",
		"templateId":       templateId,
		"label_test_label": "env-testing-" + acctest.RandString(t, 5),

		"filter_config_rai_settings_rai_filters_0_filter_type":      "HATE_SPEECH",
		"filter_config_rai_settings_rai_filters_0_confidence_level": "MEDIUM_AND_ABOVE",

		// No sdp_settings_config_type means no sdp_settings block rendered
		"filter_config_pi_and_jailbreak_filter_settings_filter_enforcement": "ENABLED",
		"filter_config_pi_and_jailbreak_filter_settings_confidence_level":   "HIGH",

		"filter_config_malicious_uri_filter_settings_filter_enforcement": "ENABLED",

		"template_metadata_custom_llm_response_safety_error_message":                 "This is a custom error message for LLM response",
		"template_metadata_log_template_operations":                                  true,
		"template_metadata_log_sanitize_operations":                                  true,
		"template_metadata_multi_language_detection_enable_multi_language_detection": true,
		"template_metadata_ignore_partial_invocation_failures":                       true,
		"template_metadata_custom_prompt_safety_error_code":                          400,
		"template_metadata_custom_prompt_safety_error_message":                       "This is a custom error message for prompt",
		"template_metadata_custom_llm_response_safety_error_code":                    401,
	}

	// Step 2: Update to include basic SDP settings
	updatedContext := map[string]interface{}{
		"location":         initialContext["location"],
		"templateId":       initialContext["templateId"],
		"label_test_label": "env-updated-" + acctest.RandString(t, 5),

		"filter_config_rai_settings_rai_filters_0_filter_type":      "DANGEROUS",
		"filter_config_rai_settings_rai_filters_0_confidence_level": "LOW_AND_ABOVE",

		"sdp_settings_config_type":                                   "basic_config", // Control sdp_settings block
		"filter_config_sdp_settings_basic_config_filter_enforcement": "ENABLED",
		// advanced keys not needed for basic
		"filter_config_sdp_settings_advanced_config_inspect_template":    "",
		"filter_config_sdp_settings_advanced_config_deidentify_template": "",


		"filter_config_pi_and_jailbreak_filter_settings_filter_enforcement": "DISABLED",
		"filter_config_pi_and_jailbreak_filter_settings_confidence_level":   "MEDIUM_AND_ABOVE",

		"filter_config_malicious_uri_filter_settings_filter_enforcement": "DISABLED",

		"template_metadata_custom_llm_response_safety_error_message":                 "Updated LLM error message",
		"template_metadata_log_template_operations":                                  false,
		"template_metadata_log_sanitize_operations":                                  false,
		"template_metadata_multi_language_detection_enable_multi_language_detection": false,
		"template_metadata_ignore_partial_invocation_failures":                       false,
		"template_metadata_custom_prompt_safety_error_code":                          404,
		"template_metadata_custom_prompt_safety_error_message":                       "Updated prompt error message",
		"template_metadata_custom_llm_response_safety_error_code":                    500,
	}

	// Single Terraform configuration template
	const config_template = `
resource "google_model_armor_template" "test-resource" {
  location    = "{{.location}}"
  template_id = "{{.templateId}}"

  labels = {
    "test-label" = "{{.label_test_label}}"
  }

  filter_config {
    rai_settings {
      rai_filters {
        filter_type      = "{{.filter_config_rai_settings_rai_filters_0_filter_type}}"
        confidence_level = "{{.filter_config_rai_settings_rai_filters_0_confidence_level}}"
      }
    }

    {{with .sdp_settings_config_type}} {{if ne . ""}}
    sdp_settings {
      {{if eq . "basic_config"}}
      basic_config {
          filter_enforcement = "{{$.filter_config_sdp_settings_basic_config_filter_enforcement}}"
      }
      {{else if eq . "advanced_config"}}
      advanced_config {
        inspect_template     = "{{$.filter_config_sdp_settings_advanced_config_inspect_template}}"
        deidentify_template  = "{{$.filter_config_sdp_settings_advanced_config_deidentify_template}}"
      }
      {{end}}
    }
    {{end}}{{end}}
    pi_and_jailbreak_filter_settings {
      filter_enforcement = "{{.filter_config_pi_and_jailbreak_filter_settings_filter_enforcement}}"
      confidence_level   = "{{.filter_config_pi_and_jailbreak_filter_settings_confidence_level}}"
    }
    malicious_uri_filter_settings {
      filter_enforcement = "{{.filter_config_malicious_uri_filter_settings_filter_enforcement}}"
    }
  }
  template_metadata {
    custom_llm_response_safety_error_message = "{{.template_metadata_custom_llm_response_safety_error_message}}"
    log_template_operations                  = {{.template_metadata_log_template_operations}}
    log_sanitize_operations                  = {{.template_metadata_log_sanitize_operations}}
    multi_language_detection {
      enable_multi_language_detection        = {{.template_metadata_multi_language_detection_enable_multi_language_detection}}
    }
    ignore_partial_invocation_failures       = {{.template_metadata_ignore_partial_invocation_failures}}
    custom_prompt_safety_error_code          = {{.template_metadata_custom_prompt_safety_error_code}}
    custom_prompt_safety_error_message       = "{{.template_metadata_custom_prompt_safety_error_message}}"
    custom_llm_response_safety_error_code    = {{.template_metadata_custom_llm_response_safety_error_code}}
  }
}
`

	// Helper function to expand the template
	expandTemplate := func(data map[string]interface{}) (string, error) {
		tmpl, err := template.New("config").Parse(config_template)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// Helper to build checks for Step 1 (no SDP)
	step1Checks := func(ctx map[string]interface{}) []resource.TestCheckFunc {
		return []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "labels.test-label", ctx["label_test_label"].(string)),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.sdp_settings.#", "0"), // Ensure no sdp_settings block
		}
	}

	// Helper to build checks for Step 2 (SDP basic)
	step2Checks := func(ctx map[string]interface{}) []resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "labels.test-label", ctx["label_test_label"].(string)),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.rai_settings.0.rai_filters.0.filter_type", ctx["filter_config_rai_settings_rai_filters_0_filter_type"].(string)),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "template_metadata.0.log_sanitize_operations", "false"),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.sdp_settings.#", "1"),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.sdp_settings.0.basic_config.#", "1"),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.sdp_settings.0.basic_config.0.filter_enforcement", ctx["filter_config_sdp_settings_basic_config_filter_enforcement"].(string)),
			resource.TestCheckResourceAttr("google_model_armor_template.test-resource", "filter_config.0.sdp_settings.0.advanced_config.#", "0"),
		}
		return checks
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckModelArmorTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{ // Step 1: Create WITHOUT SDP
				Config: func() string {
					cfg, err := expandTemplate(initialContext)
					if err != nil {
						t.Fatalf("Failed to expand initial config template: %v", err)
						return ""
					}
					return cfg
				}(),
				Check: resource.ComposeTestCheckFunc(step1Checks(initialContext)...),
			},
			{ // Step 2: Update WITH SDP basic
				Config: func() string {
					cfg, err := expandTemplate(updatedContext)
					if err != nil {
						t.Fatalf("Failed to expand updated config template: %v", err)
						return ""
					}
					return cfg
				}(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_model_armor_template.test-resource", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(step2Checks(updatedContext)...),
			},
		},
	})
}
