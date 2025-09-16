package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDialogflowVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowVersion_full1(context),
			},
			{
				ResourceName:      "google_dialogflow_version.foobar",
				RefreshState:true,
			},
			{
				Config: testAccDialogflowVersion_full2(context),
			},
			{
				ResourceName:      "google_dialogflow_version.foobar",
				RefreshState:true,
			},
		},
	})
}


func testAccDialogflowVersion_full1(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "agent_project" {
		name = "tf-test-dialogflow-%{random_suffix}"
		project_id = "tf-test-dialogflow-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
		deletion_policy = "DELETE"
	}

	resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
		disable_dependent_services = false
	}

	resource "google_service_account" "dialogflow_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}

	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflow_service_account.email}"
	}

	resource "google_dialogflow_agent" "agent" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-agent-%{random_suffix}"
		default_language_code = "en"
		time_zone = "America/New_York"
		depends_on = [google_project_iam_member.agent_create]
	}

	resource "google_dialogflow_version" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		description = "tf-test-description-%{random_suffix}"
	}
	`, context)
}

func testAccDialogflowVersion_full2(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "agent_project" {
		name = "tf-test-dialogflow-%{random_suffix}"
		project_id = "tf-test-dialogflow-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
		deletion_policy = "DELETE"
	}

	resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
		disable_dependent_services = false
	}

	resource "google_service_account" "dialogflow_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}
	
	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflow_service_account.email}"
	}

	resource "google_dialogflow_agent" "agent" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-agent-%{random_suffix}"
		default_language_code = "en"
		time_zone = "America/New_York"
		depends_on = [google_project_iam_member.agent_create]
	}

	resource "google_dialogflow_version" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		description = "tf-test-version-%{random_suffix}2"
	}
	`, context)
}
