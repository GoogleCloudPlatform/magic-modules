package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcTrigger_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate1(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-name%{random_suffix}"
  location = "europe-west1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "europe-west1"
    }
  }
  labels = {
    foo = "bar"
  }
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service%{random_suffix}"
  location = "europe-west1"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-name%{random_suffix}"
  location = "europe-west1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "europe-west1"
    }
  }
  transport {
    pubsub {
      topic = google_pubsub_topic.foo.id
    }
  }
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service%{random_suffix}"
  location = "europe-west1"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service" "default2" {
  name     = "tf-test-eventarc-service%{random_suffix}2"
  location = "europe-north1"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-name%{random_suffix}"
  location = "europe-west1"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default2.name
      region  = "europe-north1"
    }
  }
  transport {
    pubsub {
      topic = google_pubsub_topic.foo.id
    }
  }
  labels = {
    foo = "bar"
  }
  service_account = google_service_account.eventarc-sa.email
}

resource "google_service_account" "eventarc-sa" {
  account_id   = "tf-test-sa%{random_suffix}"
  display_name = "Test Service Account"
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service%{random_suffix}"
  location = "europe-west1"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service" "default2" {
  name     = "tf-test-eventarc-service%{random_suffix}2"
  location = "europe-north1"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

`, context)
}

func TestAccEventarcTrigger_AlternateForm(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()

	context := map[string]interface{}{
		"region":        region,
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_AlternateForm(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// This tests the long-form for name, location, and project, and the short-form for transport.pubsub.topic
				ImportStateVerifyIgnore: []string{"name", "location", "project", "transport.0.pubsub.0.topic"},
			},
		},
	})
}

func testAccEventarcTrigger_AlternateForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service%{random_suffix}"
  location = "%{region}"
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        ports {
          container_port = 8080
        }
      }
      container_concurrency = 50
      timeout_seconds       = 100
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_eventarc_trigger" "primary" {
  name     = "projects/%{project_name}/locations/%{region}/triggers/tf-test-trigger%{random_suffix}"
  project  = "projects/%{project_name}"
  location = "long/form/%{region}"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region  = "%{region}"
    }
  }
  transport {
    pubsub {
      topic = "tf-test-topic%{random_suffix}"
    }
  }
  depends_on = [google_cloud_run_service.default, google_pubsub_topic.foo]
}
`, context)
}

func TestAccEventarcTrigger_ShortFormWorkflow(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"region":          envvar.GetTestRegionFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_ShortFormWorkflow(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_ShortFormWorkflow(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-name%{random_suffix}"
  location = "%{region}"
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    workflow = "tf-test-eventarc-workflow%{random_suffix}"
  }
  service_account = "%{service_account}"
  depends_on      = [google_workflows_workflow.example]
}

resource "google_workflows_workflow" "example" {
  name                = "tf-test-eventarc-workflow%{random_suffix}"
  deletion_protection = false
  region              = "%{region}"
  source_contents     = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: $${sys.get_env("url")}
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https://en.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, context)
}
