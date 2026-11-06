package discoveryengine_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDiscoveryEngineDocumentImport_lifecycleAndPurgeOnDestroy(t *testing.T) {
	t.Parallel()

	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
			Role:   "roles/storage.admin",
		},
	})

	randomSuffix := acctest.RandString(t, 10)
	dataStoreID := "tf-test-doc-import-ds-" + randomSuffix
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"data_store_id": dataStoreID,
		"bucket_name":   "tf-test-doc-import-bkt-" + randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineDocumentImportDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Step 1: Import 2 documents with purge_on_destroy = false and verify
				// both Terraform state counts and live ListDocuments count == 2.
				Config: testAccDiscoveryEngineDocumentImport_withImport(context, "tf-test-first-import", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "project", context["project"].(string)),
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "success_count", "2"),
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "failure_count", "0"),
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "total_count", "2"),
					resource.TestCheckResourceAttrSet("google_discovery_engine_document_import.test", "operation"),
					testAccCheckDiscoveryEngineBranchDocumentCount(t, "global", dataStoreID, 2),
				),
			},
			{
				// Step 2: Destroy only the DocumentImport resource while keeping the
				// DataStore alive; purge_on_destroy = false must leave the 2 documents intact.
				Config: testAccDiscoveryEngineDocumentImport_dataStoreOnly(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiscoveryEngineBranchDocumentCount(t, "global", dataStoreID, 2),
				),
			},
			{
				// Step 3: Create a second DocumentImport on the same branch (starting with
				// purge_on_destroy = false).
				Config: testAccDiscoveryEngineDocumentImport_withImport(context, "tf-test-second-import", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "success_count", "2"),
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "failure_count", "0"),
					testAccCheckDiscoveryEngineBranchDocumentCount(t, "global", dataStoreID, 2),
				),
			},
			{
				// Step 4: Update purge_on_destroy from false -> true in-place (must NOT trigger ForceNew).
				Config: testAccDiscoveryEngineDocumentImport_withImport(context, "tf-test-second-import", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_document_import.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_discovery_engine_document_import.test", "purge_on_destroy", "true"),
					testAccCheckDiscoveryEngineBranchDocumentCount(t, "global", dataStoreID, 2),
				),
			},
			{
				// Step 5: Destroy the DocumentImport resource with purge_on_destroy = true while
				// keeping the DataStore alive; PurgeDocuments (force=true) must empty the branch.
				Config: testAccDiscoveryEngineDocumentImport_dataStoreOnly(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiscoveryEngineBranchDocumentCount(t, "global", dataStoreID, 0),
				),
			},
		},
	})
}

func TestAccDiscoveryEngineDocumentImport_createFailsOnDocumentError(t *testing.T) {
	t.Parallel()

	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
			Role:   "roles/storage.admin",
		},
	})

	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"data_store_id": "tf-test-doc-import-err-" + randomSuffix,
		"bucket_name":   "tf-test-doc-import-err-" + randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineDocumentImportDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDiscoveryEngineDocumentImport_invalidDocument(context),
				ExpectError: regexp.MustCompile(`1 document\(s\) failed out of 1 total`),
			},
		},
	})
}

func testAccCheckDiscoveryEngineBranchDocumentCount(t *testing.T, location, dataStoreID string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		project := envvar.GetTestProjectFromEnv()
		billingProject := project
		if config.BillingProject != "" {
			billingProject = config.BillingProject
		}

		baseURL := strings.ReplaceAll(transport_tpg.BaseUrl(discoveryengine.Product, config), "{{location}}", location)
		rawURL := fmt.Sprintf("%sprojects/%s/locations/%s/collections/default_collection/dataStores/%s/branches/default_branch/documents",
			baseURL, project, location, dataStoreID)

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    rawURL,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return fmt.Errorf("error listing documents in branch %s: %s", rawURL, err)
		}

		docs, _ := res["documents"].([]interface{})
		if len(docs) != expected {
			return fmt.Errorf("expected %d document(s) in branch %s, found %d", expected, rawURL, len(docs))
		}
		return nil
	}
}

func testAccDiscoveryEngineDocumentImport_dataStoreOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "%{data_store_id}"
  display_name                 = "tf-test-datastore-doc-import-lifecycle"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}

resource "google_storage_bucket" "import_bucket" {
  name                        = "%{bucket_name}"
  location                    = "US"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "discoveryengine_access" {
  bucket     = google_storage_bucket.import_bucket.name
  role       = "roles/storage.objectViewer"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com"
  depends_on = [google_discovery_engine_data_store.basic]
}

resource "google_storage_bucket_object" "documents" {
  name    = "documents.jsonl"
  bucket  = google_storage_bucket.import_bucket.name
  content = <<-EOT
    {"doc_id":"doc-1","title":"Lifecycle Document 1"}
    {"doc_id":"doc-2","title":"Lifecycle Document 2"}
  EOT
}
`, context)
}

func testAccDiscoveryEngineDocumentImport_withImport(context map[string]interface{}, importID string, purgeOnDestroy bool) string {
	stepContext := map[string]interface{}{
		"data_store_id":    context["data_store_id"],
		"bucket_name":      context["bucket_name"],
		"import_id":        importID,
		"purge_on_destroy": purgeOnDestroy,
	}
	return testAccDiscoveryEngineDocumentImport_dataStoreOnly(stepContext) + acctest.Nprintf(`
resource "google_discovery_engine_document_import" "test" {
  location            = google_discovery_engine_data_store.basic.location
  data_store_id       = google_discovery_engine_data_store.basic.data_store_id
  branch_id           = "default_branch"
  import_id           = "%{import_id}"
  reconciliation_mode = "INCREMENTAL"
  auto_generate_ids   = false
  id_field            = "doc_id"
  purge_on_destroy    = %{purge_on_destroy}

  gcs_source {
    input_uris  = ["gs://${google_storage_bucket.import_bucket.name}/${google_storage_bucket_object.documents.name}"]
    data_schema = "custom"
  }

  depends_on = [google_storage_bucket_iam_member.discoveryengine_access]
}
`, stepContext)
}

func testAccDiscoveryEngineDocumentImport_invalidDocument(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "%{data_store_id}"
  display_name                 = "tf-test-datastore-doc-import-err"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}

resource "google_storage_bucket" "import_bucket" {
  name                        = "%{bucket_name}"
  location                    = "US"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "discoveryengine_access" {
  bucket     = google_storage_bucket.import_bucket.name
  role       = "roles/storage.objectViewer"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com"
  depends_on = [google_discovery_engine_data_store.basic]
}

resource "google_storage_bucket_object" "invalid_documents" {
  name    = "invalid.jsonl"
  bucket  = google_storage_bucket.import_bucket.name
  content = "not-valid-json-document\n"
}

resource "google_discovery_engine_document_import" "invalid" {
  location      = google_discovery_engine_data_store.basic.location
  data_store_id = google_discovery_engine_data_store.basic.data_store_id
  import_id     = "tf-test-invalid-import"

  gcs_source {
    input_uris  = ["gs://${google_storage_bucket.import_bucket.name}/${google_storage_bucket_object.invalid_documents.name}"]
    data_schema = "document"
  }

  depends_on = [google_storage_bucket_iam_member.discoveryengine_access]
}
`, context)
}
