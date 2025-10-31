package ces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCESGuardrail_cesGuardrailBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESGuardrailDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESGuardrail_cesGuardrailBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
			{
				Config: testAccCESGuardrail_cesGuardrailBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_guardrail.ces_guardrail_basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
		},
	})
}

func testAccCESGuardrail_cesGuardrailBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_basic" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    respond_immediately  {
        responses {
            text = "Text"
            disabled = false
        }
    }
  }
  enabled = true
  model_safety  {
    safety_settings {
        category = "HARM_CATEGORY_HATE_SPEECH"
        threshold = "BLOCK_NONE"
    }
  }
}
`, context)
}

func testAccCESGuardrail_cesGuardrailBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_basic" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description updated"
  action {
    respond_immediately  {
        responses {
            text = "Text updated"
            disabled = true
        }
    }
  }
  enabled = false
  model_safety  {
    safety_settings {
        category = "HARM_CATEGORY_HATE_SPEECH"
        threshold = "BLOCK_NONE"
    }
  }
}
`, context)
}

func TestAccCESGuardrail_cesGuardrailTransferAgentContentFilterExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESGuardrailDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESGuardrail_cesGuardrailTransferAgentContentFilterExample_full(context),
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_transfer_agent_content_filter",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
			{
				Config: testAccCESGuardrail_cesGuardrailTransferAgentContentFilterExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_guardrail.ces_guardrail_transfer_agent_content_filter", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_transfer_agent_content_filter",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
		},
	})
}

func testAccCESGuardrail_cesGuardrailTransferAgentContentFilterExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_transfer_agent_content_filter" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description updated"
  action {
    transfer_agent {
        agent = "projects/${google_ces_app.ces_app_for_guardrail.project}/locations/us/apps/${google_ces_app.ces_app_for_guardrail.app_id}/agents/fake-agent"
    }
  }
  enabled = true
  content_filter {
    banned_contents = ["example"]
    banned_contents_in_user_input = ["example"]
    banned_contents_in_agent_response = ["example"]
    match_type = "SIMPLE_STRING_MATCH"
    disregard_diacritics = false
  }
}
`, context)
}

func testAccCESGuardrail_cesGuardrailTransferAgentContentFilterExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_transfer_agent_content_filter" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    transfer_agent {
        agent = "projects/${google_ces_app.ces_app_for_guardrail.project}/locations/us/apps/${google_ces_app.ces_app_for_guardrail.app_id}/agents/fake-agent-updated"
    }
  }
  enabled = true
  content_filter {
    banned_contents = ["example_updated"]
    banned_contents_in_user_input = ["example_updated"]
    banned_contents_in_agent_response = ["example_updated"]
    match_type = "SIMPLE_STRING_MATCH"
    disregard_diacritics = false
  }
}
`, context)
}

func TestAccCESGuardrail_cesGuardrailGenerativeAnswerLlmPromptSecurityExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESGuardrailDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESGuardrail_cesGuardrailGenerativeAnswerLlmPromptSecurityExample_full(context),
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_generative_answer_llm_prompt_security",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
			{
				Config: testAccCESGuardrail_cesGuardrailGenerativeAnswerLlmPromptSecurityExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_guardrail.ces_guardrail_generative_answer_llm_prompt_security", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_generative_answer_llm_prompt_security",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
		},
	})
}

func testAccCESGuardrail_cesGuardrailGenerativeAnswerLlmPromptSecurityExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_generative_answer_llm_prompt_security" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt"
    }
  }
  enabled = true
  llm_prompt_security {
    custom_policy {
      max_conversation_messages = 10
      model_settings {
        model = "gemini-2.5-flash"
        temperature = 50
      }
      prompt = "example_prompt"
      policy_scope = "USER_QUERY"
      fail_open = true
      allow_short_utterance = true
    }
  }
}
`, context)
}

func testAccCESGuardrail_cesGuardrailGenerativeAnswerLlmPromptSecurityExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_generative_answer_llm_prompt_security" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt_updated"
    }
  }
  enabled = true
  llm_prompt_security {
    custom_policy {
      max_conversation_messages = 9
      model_settings {
        model = "gemini-2.0-flash"
        temperature = 49
      }
      prompt = "example_prompt_updated"
      policy_scope = "USER_QUERY"
      fail_open = false
      allow_short_utterance = false
    }
  }
}
`, context)
}

func TestAccCESGuardrail_cesGuardrailCodeCallbackExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESGuardrailDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESGuardrail_cesGuardrailCodeCallbackExample_full(context),
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_code_callback",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
			{
				Config: testAccCESGuardrail_cesGuardrailCodeCallbackExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_guardrail.ces_guardrail_code_callback", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_code_callback",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
		},
	})
}

func testAccCESGuardrail_cesGuardrailCodeCallbackExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_code_callback" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt"
    }
  }
  enabled = true
  code_callback {
    before_agent_callback {
        description = "Example callback"
        disabled    = false
        python_code = "def callback(context):\n    return {'override': true}"
    }
    after_agent_callback {
        description = "Example callback"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': true}"
    }
    before_model_callback {
        description = "Example callback"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': true}"
    }
    after_model_callback {
        description = "Example callback"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': true}"
    }
  }
}
`, context)
}

func testAccCESGuardrail_cesGuardrailCodeCallbackExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_code_callback" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt"
    }
  }
  enabled = true
  code_callback {
    before_agent_callback {
        description = "Example callback updated"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': False}"
    }
    after_agent_callback {
        description = "Example callback updated"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': False}"
    }
    before_model_callback {
        description = "Example callback updated"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': False}"
    }
    after_model_callback {
        description = "Example callback updated"
        disabled    = true
        python_code = "def callback(context):\n    return {'override': False}"
    }
  }
}
`, context)
}

func TestAccCESGuardrail_cesGuardrailLlmPolicyExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESGuardrailDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESGuardrail_cesGuardrailLlmPolicyExample_full(context),
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_llm_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
			{
				Config: testAccCESGuardrail_cesGuardrailLlmPolicyExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_guardrail.ces_guardrail_llm_policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_guardrail.ces_guardrail_llm_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id", "guardrail_id"},
			},
		},
	})
}

func testAccCESGuardrail_cesGuardrailLlmPolicyExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_llm_policy" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt"
    }
  }
  enabled = true
  llm_policy {
    max_conversation_messages = 10
    model_settings {
        model = "gemini-2.5-flash"
        temperature = 50
    }
    prompt = "example_prompt"
    policy_scope = "USER_QUERY"
    fail_open = true
    allow_short_utterance = true
  }
}
`, context)
}

func testAccCESGuardrail_cesGuardrailLlmPolicyExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "ces_app_for_guardrail" {
  app_id = "tf-test-app-id%{random_suffix}"
  location = "us"
  description = "App used as parent for CES Guardrail example"
  display_name = "tf-test-my-app%{random_suffix}"

  language_settings {
    default_language_code    = "en-US"
    supported_language_codes = ["es-ES", "fr-FR"]
    enable_multilingual_support = true
    fallback_action          = "escalate"
  }
  time_zone_settings {
    time_zone = "America/Los_Angeles"
  }
}

resource "google_ces_guardrail" "ces_guardrail_llm_policy" {
  guardrail_id = "tf-test-guardrail-id%{random_suffix}"
  location     = google_ces_app.ces_app_for_guardrail.location
  app          = google_ces_app.ces_app_for_guardrail.app_id
  display_name = "tf-test-my-guardrail%{random_suffix}"
  description  = "Guardrail description"
  action {
    generative_answer {
        prompt = "example_prompt"
    }
  }
  enabled = true
  llm_policy {
    max_conversation_messages = 8
    model_settings {
        model = "gemini-2.0-flash"
        temperature = 45
    }
    prompt = "example_prompt_updated"
    policy_scope = "USER_QUERY"
    fail_open = false
    allow_short_utterance = false
  }
}
`, context)
}
