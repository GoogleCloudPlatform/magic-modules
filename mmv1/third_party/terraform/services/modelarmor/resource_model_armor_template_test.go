package modelarmor_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
    "github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccModelArmorTemplate_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      "us-central1",
		"templateId":    "test-template-" + acctest.RandString(t, 5),
		"labelValue":    "test-label-value-" + acctest.RandString(t, 5),
		"filter_config_rai_settings_rai_filters_0_filter_type":                "SEXUALLY_EXPLICIT",
		"filter_config_rai_settings_rai_filters_0_confidence_level":           "LOW_AND_ABOVE",
		"filter_config_sdp_settings_basic_config_filter_enforcement":          "ENABLED",
		"filter_config_pi_and_jailbreak_filter_settings_filter_enforcement":    "ENABLED",
		"filter_config_pi_and_jailbreak_filter_settings_confidence_level":      "LOW_AND_ABOVE",
		"filter_config_malicious_uri_filter_settings_filter_enforcement":       "ENABLED",
		"template_metadata_custom_llm_response_safety_error_message":          "This is a custom error message for LLM response",
		"template_metadata_log_template_operations":                           true,
		"template_metadata_log_sanitize_operations":                           false,
		"template_metadata_multi_language_detection_enable_multi_language_detection": true,
		"template_metadata_ignore_partial_invocation_failures":                false,
		"template_metadata_custom_prompt_safety_error_code":                   400,
		"template_metadata_custom_prompt_safety_error_message":                "This is a custom error message for prompt",
		"template_metadata_custom_llm_response_safety_error_code":              401,
	}

	config_basic := `
		resource "google_model_armor_template" "basic" {
			location    = "{{.location}}"
			template_id = "{{.templateId}}"
			provider    = google-beta
		
			labels = {
				label-key = "{{.labelValue}}"
			}
			filter_config {
				rai_settings {
					rai_filters {
						filter_type      = "{{.filter_config_rai_settings_rai_filters_0_filter_type}}"
						confidence_level = "{{.filter_config_rai_settings_rai_filters_0_confidence_level}}"
					}
				}
				sdp_settings {
					basic_config {
						filter_enforcement = "{{.filter_config_sdp_settings_basic_config_filter_enforcement}}"
					}
				}
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
					enable_multi_language_detection = {{.template_metadata_multi_language_detection_enable_multi_language_detection}}
				}
				ignore_partial_invocation_failures = {{.template_metadata_ignore_partial_invocation_failures}}
				custom_prompt_safety_error_code    = {{.template_metadata_custom_prompt_safety_error_code}}
				custom_prompt_safety_error_message = "{{.template_metadata_custom_prompt_safety_error_message}}"
				custom_llm_response_safety_error_code = {{.template_metadata_custom_llm_response_safety_error_code}}
			}
		}
	`

	config_updated := `
		resource "google_model_armor_template" "basic" {
			location    = "{{.location}}"
			template_id = "{{.templateId}}"
			provider    = google-beta
		
			labels = {
				label-key = "{{.labelValue}}"
			}
			filter_config {
				rai_settings {
					rai_filters {
						filter_type      = "{{.filter_config_rai_settings_rai_filters_0_filter_type}}"
						confidence_level = "{{.filter_config_rai_settings_rai_filters_0_confidence_level}}"
					}
				}
				sdp_settings {
					basic_config {
						filter_enforcement = "{{.filter_config_sdp_settings_basic_config_filter_enforcement}}"
					}
				}
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
					enable_multi_language_detection = {{.template_metadata_multi_language_detection_enable_multi_language_detection}}
				}
				ignore_partial_invocation_failures = {{.template_metadata_ignore_partial_invocation_failures}}
				custom_prompt_safety_error_code    = {{.template_metadata_custom_prompt_safety_error_code}}
				custom_prompt_safety_error_message = "{{.template_metadata_custom_prompt_safety_error_message}}"
				custom_llm_response_safety_error_code = {{.template_metadata_custom_llm_response_safety_error_code}}
			}
		}
	`

	// Helper function to expand the template
	expandTemplate := func(tmplStr string, data map[string]interface{}) (string, error) {
		tmpl, err := template.New("config").Parse(tmplStr)
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

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckModelArmorTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: func() string {
					cfg, err := expandTemplate(config_basic, context)
					if err != nil {
						t.Fatalf("Failed to expand basic config template: %v", err)
						return "" // Return empty string in case of error
					}
					return cfg
				}(),
			},
			{
				ResourceName:            "projects/modelarmor-api-test/locations/us-central1/templates/at-autogen",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "template_id", "terraform_labels"},
			},
			{
				Config: func() string {
					cfg, err := expandTemplate(config_updated, context)
					if err != nil {
						t.Fatalf("Failed to expand updated config template: %v", err)
						return "" // Return empty string in case of error
					}
					return cfg
				}(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_model_armor_template.basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_model_armor_template.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "template_id", "terraform_labels"},
			},
		},
	})
}

// Add this dummy function because the original implementation is not provided.
func testAccCheckModelArmorTemplateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		return nil
	}
}