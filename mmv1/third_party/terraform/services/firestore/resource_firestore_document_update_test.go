package firestore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/tpgresource"

	"github.com/hashicorp/terraform-provider-google-beta/google/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google/envvar"
)

func TestAccFirestoreDocument_firestoreDocumentUpdateExampleEmptyField(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckFirestoreDocumentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDocument_firestoreDocumentBasicInitialConfiguration(context),
			},
			{
				ResourceName:            "google_firestore_document.mydoc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection", "database", "document_id"},
			},
			// New test steps to address the empty fields diff issue
			{
				// This step creates a document with empty fields
				Config: testAccFirestoreDocument_firestoreDocumentUpdateWithEmptyFields(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_firestore_document.empty_doc", "name"),
					resource.TestCheckResourceAttr("google_firestore_document.empty_doc", "fields", "{}"),
				),
			},
			{
				// This step asserts that a plan on the empty document shows no diff
				Config:             testAccFirestoreDocument_firestoreDocumentUpdateWithEmptyFields(context), // Apply the same config again
				PlanOnly:           true,                                                                     // runs terraform plan
				ExpectNonEmptyPlan: false,                                                                    // nodiff expected
			},
		},
	})

}

func testAccFirestoreDocument_firestoreDocumentBasicInitialConfiguration(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id = "tf-test-project-id%{random_suffix}"
  name       = "tf-test-project-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.project]

  create_duration = "60s"
}

resource "google_project_service" "firestore" {
  project = google_project.project.project_id
  service = "firestore.googleapis.com"

  # Needed for CI tests for permissions to propagate, should not be needed for actual usage
  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_firestore_database" "database" {
  project     = google_project.project.project_id
  name        = "(default)"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  depends_on = [google_project_service.firestore]
}

resource "google_firestore_document" "mydoc" {
  project     = google_project.project.project_id
  database    = google_firestore_database.database.name
  collection  = "somenewcollection"
  document_id = "tf-test-my-doc-id%{random_suffix}"
  fields      = "{\"something\":{\"mapValue\":{\"fields\":{\"akey\":{\"stringValue\":\"avalue\"}}}}}"
}
`, context)
}

func testAccFirestoreDocument_firestoreDocumentUpdateWithEmptyFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id = "tf-test-project-id%{random_suffix}"
  name       = "tf-test-project-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.project]

  create_duration = "60s"
}

resource "google_project_service" "firestore" {
  project = google_project.project.project_id
  service = "firestore.googleapis.com"

  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_firestore_database" "database" {
  project     = google_project.project.project_id
  name        = "(default)"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  depends_on = [google_project_service.firestore]
}

resource "google_firestore_document" "empty_doc" {
  project     = google_project.project.project_id
  database    = google_firestore_database.database.name
  collection  = "emptycollection"
  document_id = "tf-test-empty-doc-id%{random_suffix}"
  fields      = jsonencode({}) # This is the key: an empty JSON object
}
`, context)
}

func testAccCheckFirestoreDocumentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_firestore_document" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{FirestoreBasePath}}{{name}}")
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
				return fmt.Errorf("FirestoreDocument still exists at %s", url)
			}
		}

		return nil
	}
}
