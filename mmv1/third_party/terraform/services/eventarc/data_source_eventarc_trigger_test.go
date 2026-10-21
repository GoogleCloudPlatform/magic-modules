package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/cloudrun"
	_ "github.com/hashicorp/terraform-provider-google/google/services/eventarc"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
)

func TestAccDataSourceGoogleEventarcTrigger_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleEventarcTrigger_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_eventarc_trigger.foo", "google_eventarc_trigger.primary"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleEventarcTrigger_optionalProjectLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleEventarcTrigger_optionalProjectLocation(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_eventarc_trigger.foo", "google_eventarc_trigger.primary"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleEventarcTrigger_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-trigger-%{random_suffix}"
  location = "%{region}"
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
  labels = {
    foo = "bar"
  }
  transport {
    pubsub {
      topic = google_pubsub_topic.foo.id
    }
  }
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service-%{random_suffix}"
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

data "google_eventarc_trigger" "foo" {
  name     = google_eventarc_trigger.primary.name
  location = google_eventarc_trigger.primary.location
  project  = google_eventarc_trigger.primary.project
}
`, context)
}

func testAccDataSourceGoogleEventarcTrigger_optionalProjectLocation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
  name     = "tf-test-trigger-%{random_suffix}"
  location = "%{region}"
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
      topic = google_pubsub_topic.foo.id
    }
  }
}

resource "google_pubsub_topic" "foo" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-eventarc-service-%{random_suffix}"
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

data "google_eventarc_trigger" "foo" {
  name = google_eventarc_trigger.primary.name
}
`, context)
}
