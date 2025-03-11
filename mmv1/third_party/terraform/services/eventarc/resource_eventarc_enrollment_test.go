package eventarc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccEventarcEnrollment_update(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"region":          region,
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcEnrollmentDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcEnrollment_full(context),
			},
			{
				ResourceName:            "google_eventarc_enrollment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "enrollment_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccEventarcEnrollment_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_enrollment.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_enrollment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "enrollment_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccEventarcEnrollment_unset(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_enrollment.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_enrollment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "enrollment_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcEnrollment_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_create_project" {
  create_duration = "60s"
  depends_on      = [google_project.project]
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
  depends_on = [time_sleep.wait_create_project]
}

resource "google_project_service" "pubsub" {
  project    = google_project.project.project_id
  service    = "pubsub.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "google_project_service" "eventarc" {
  project    = google_project.project.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.pubsub]
}

resource "time_sleep" "wait_enable_service" {
  create_duration = "60s"
  depends_on      = [google_project_service.eventarc]
}

resource "google_project_service_identity" "eventarc_sa" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_enable_service]
}

resource "time_sleep" "wait_create_sa" {
  create_duration = "60s"
  depends_on      = [google_project_service_identity.eventarc_sa]
}

resource "google_eventarc_enrollment" "primary" {
  location      = "%{region}"
  enrollment_id = "tf-test-enrollment%{random_suffix}"
  project       = google_project.project.project_id
  display_name  = "basic enrollment"
  message_bus   = google_eventarc_message_bus.message_bus.id
  destination   = google_eventarc_pipeline.pipeline.id
  cel_match     = "message.type == 'google.cloud.dataflow.job.v1beta3.statusChanged'"
  labels = {
    test_label = "test-eventarc-label"
  }
  annotations = {
    test_annotation = "test-eventarc-annotation"
  }
  depends_on = [time_sleep.wait_create_sa]
}

resource "google_compute_network" "psc" {
  name                    = "tf-test-network%{random_suffix}"
  project                 = google_project.project.project_id
  auto_create_subnetworks = false
  depends_on              = [time_sleep.wait_enable_service]
}

resource "google_compute_subnetwork" "psc" {
  name          = "tf-test-subnet%{random_suffix}"
  region        = "%{region}"
  project       = google_project.project.project_id
  network       = google_compute_network.psc.id
  ip_cidr_range = "10.77.0.0/20"
}

resource "google_compute_network_attachment" "psc" {
  name                  = "tf-test-na%{random_suffix}"
  region                = "%{region}"
  project               = google_project.project.project_id
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [
    google_compute_subnetwork.psc.self_link
  ]
}

resource "time_sleep" "wait_create_na" {
  create_duration = "60s"
  depends_on      = [google_compute_network_attachment.psc]
}

resource "google_pubsub_topic" "pipeline_topic" {
  name       = "tf-test-topic%{random_suffix}"
  depends_on = [time_sleep.wait_enable_service]
}

resource "google_eventarc_pipeline" "pipeline" {
  location    = "%{region}"
  pipeline_id = "tf-test-pipeline%{random_suffix}"
  project     = google_project.project.project_id
  destinations {
    topic = google_pubsub_topic.pipeline_topic.id
    network_config {
      network_attachment = google_compute_network_attachment.psc.id
    }
  }
  depends_on = [time_sleep.wait_create_sa, time_sleep.wait_create_na]
}

resource "google_eventarc_message_bus" "message_bus" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
  project        = google_project.project.project_id
  depends_on     = [time_sleep.wait_create_sa]
}
`, context)
}

func testAccEventarcEnrollment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
}

resource "google_project_service" "pubsub" {
  project    = google_project.project.project_id
  service    = "pubsub.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "google_project_service" "eventarc" {
  project    = google_project.project.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.pubsub]
}

resource "google_project_service_identity" "eventarc_sa" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project.project_id
  depends_on = [google_project_service.eventarc]
}

resource "google_eventarc_enrollment" "primary" {
  location      = "%{region}"
  enrollment_id = "tf-test-enrollment%{random_suffix}"
  project       = google_project.project.project_id
  display_name  = "basic updated enrollment"
  message_bus   = google_eventarc_message_bus.message_bus_update.id
  destination   = google_eventarc_pipeline.pipeline_update.id
  cel_match     = "true"
  labels = {
    updated_label = "updated-test-eventarc-label"
  }
  annotations = {
    updated_test_annotation = "updated-test-eventarc-annotation"
  }
  # TODO(tommyreddad) As of time of writing, enrollments can't be updated
  # if their pipeline has been deleted. So use this workaround until the
  # underlying issue in the Eventarc API is fixed.
  depends_on    = [google_eventarc_pipeline.pipeline, google_project_service_identity.eventarc_sa]
}

resource "google_compute_network" "psc" {
  name                    = "tf-test-network%{random_suffix}"
  project                 = google_project.project.project_id
  auto_create_subnetworks = false
  depends_on              = [google_project_service.compute]
}

resource "google_compute_subnetwork" "psc" {
  name          = "tf-test-subnet%{random_suffix}"
  region        = "%{region}"
  project       = google_project.project.project_id
  network       = google_compute_network.psc.id
  ip_cidr_range = "10.77.0.0/20"
}

resource "google_compute_network_attachment" "psc" {
  name                  = "tf-test-na%{random_suffix}"
  region                = "%{region}"
  project               = google_project.project.project_id
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [
    google_compute_subnetwork.psc.self_link
  ]
}

resource "google_pubsub_topic" "pipeline_topic" {
  name       = "tf-test-topic%{random_suffix}"
  depends_on = [google_project_service.pubsub]
}

resource "google_eventarc_pipeline" "pipeline" {
  location    = "%{region}"
  pipeline_id = "tf-test-pipeline%{random_suffix}"
  project     = google_project.project.project_id
  destinations {
    topic = google_pubsub_topic.pipeline_topic.id
    network_config {
      network_attachment = google_compute_network_attachment.psc.id
    }
  }
  depends_on = [google_project_service_identity.eventarc_sa]
}

resource "google_pubsub_topic" "pipeline_update_topic" {
  name       = "tf-test-topic2%{random_suffix}"
  depends_on = [google_project_service.pubsub]
}

resource "google_eventarc_pipeline" "pipeline_update" {
  location    = "%{region}"
  pipeline_id = "tf-test-pipeline2%{random_suffix}"
  project     = google_project.project.project_id
  destinations {
    topic = google_pubsub_topic.pipeline_update_topic.id
    network_config {
      network_attachment = google_compute_network_attachment.psc.id
    }
  }
  depends_on = [google_project_service_identity.eventarc_sa]
}

resource "google_project" "project_update" {
  project_id      = "tf-test2%{random_suffix}"
  name            = "tf-test2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_create_project_update" {
  create_duration = "60s"
  depends_on      = [google_project.project_update]
}

resource "google_project_service" "eventarc_update" {
  project    = google_project.project_update.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [time_sleep.wait_create_project_update]
}

resource "time_sleep" "wait_enable_service_update" {
  create_duration = "60s"
  depends_on      = [google_project_service.eventarc_update]
}

resource "google_project_service_identity" "eventarc_sa_update" {
  project    = google_project.project_update.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [time_sleep.wait_enable_service_update]
}

resource "time_sleep" "wait_create_sa_update" {
  create_duration = "60s"
  depends_on      = [google_project_service_identity.eventarc_sa_update]
}

resource "google_eventarc_message_bus" "message_bus_update" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus2%{random_suffix}"
  project        = google_project.project_update.project_id
  depends_on     = [time_sleep.wait_create_sa_update]
}
`, context)
}

func testAccEventarcEnrollment_unset(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
}

resource "google_project_service" "pubsub" {
  project    = google_project.project.project_id
  service    = "pubsub.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "google_project_service" "eventarc" {
  project    = google_project.project.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.pubsub]
}

resource "google_project_service_identity" "eventarc_sa" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project.project_id
  depends_on = [google_project_service.eventarc]
}

resource "google_eventarc_enrollment" "primary" {
  location      = "%{region}"
  enrollment_id = "tf-test-enrollment%{random_suffix}"
  project       = google_project.project.project_id
  message_bus   = google_eventarc_message_bus.message_bus_update.id
  destination   = google_eventarc_pipeline.pipeline_update.id
  cel_match     = "true"
  depends_on    = [google_project_service_identity.eventarc_sa]
}

resource "google_compute_network" "psc" {
  name                    = "tf-test-network%{random_suffix}"
  project                 = google_project.project.project_id
  auto_create_subnetworks = false
  depends_on              = [google_project_service.compute]
}

resource "google_compute_subnetwork" "psc" {
  name          = "tf-test-subnet%{random_suffix}"
  region        = "%{region}"
  project       = google_project.project.project_id
  network       = google_compute_network.psc.id
  ip_cidr_range = "10.77.0.0/20"
}

resource "google_compute_network_attachment" "psc" {
  name                  = "tf-test-na%{random_suffix}"
  region                = "%{region}"
  project               = google_project.project.project_id
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [
    google_compute_subnetwork.psc.self_link
  ]
}

resource "google_pubsub_topic" "pipeline_update_topic" {
  name       = "tf-test-topic2%{random_suffix}"
  depends_on = [google_project_service.pubsub]
}

resource "google_eventarc_pipeline" "pipeline_update" {
  location    = "%{region}"
  pipeline_id = "tf-test-pipeline2%{random_suffix}"
  project     = google_project.project.project_id
  destinations {
    topic = google_pubsub_topic.pipeline_update_topic.id
    network_config {
      network_attachment = google_compute_network_attachment.psc.id
    }
  }
  depends_on = [google_project_service_identity.eventarc_sa]
}

resource "google_project" "project_update" {
  project_id      = "tf-test2%{random_suffix}"
  name            = "tf-test2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "eventarc_update" {
  project    = google_project.project_update.project_id
  service    = "eventarc.googleapis.com"
}

resource "google_project_service_identity" "eventarc_sa_update" {
  project    = google_project.project_update.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.eventarc_update]
}

resource "google_eventarc_message_bus" "message_bus_update" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus2%{random_suffix}"
  project        = google_project.project_update.project_id
  depends_on     = [google_project_service_identity.eventarc_sa_update]
}
`, context)
}

func testAccCheckEventarcEnrollmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_enrollment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{EventarcBasePath}}projects/{{project}}/locations/{{location}}/enrollments/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("EventarcEnrollment still exists at %s", url)
			}
		}

		return nil
	}
}
