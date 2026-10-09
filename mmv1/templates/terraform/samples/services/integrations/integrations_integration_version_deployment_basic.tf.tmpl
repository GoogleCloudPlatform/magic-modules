resource "google_integrations_integration_version" "version" {
  location = "us-east4"
  integration = "%{integration_name}"
  description = "Integration version created via terraform"

  trigger_configs {
    label = "test_trigger"
    trigger_type = "API"
    trigger_number = "1"
    trigger_id = "api_trigger/test_trigger"
    properties = {
      "Trigger name" = "test_trigger"
    }
    start_tasks {
      task_id = "1"
    }
  }

  task_configs {
    task = "GenericRestV2Task"
    task_id = "1"
    external_task_type = "NORMAL_TASK"
    task_execution_strategy = "WHEN_ALL_SUCCEED"
    parameters {
      key = "url"
      masked = false
      value {
        boolean_value = false
        double_value = 0
        string_value = "https://example.com"
      }
    }
    parameters {
      key = "httpMethod"
      masked = false
      value {
        boolean_value = false
        double_value = 0
        string_value = "GET"
      }
    }
  }
}

resource "google_integrations_integration_version_deployment" "deployment_basic" {
  location = "us-east4"
  integration = "%{integration_name}"
  version = google_integrations_integration_version.version.version
}
