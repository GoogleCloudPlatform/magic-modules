package vectorsearch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVectorSearchCollection_update(t *testing.T) {
	t.Parallel()

	runId := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"collection_id": "tf-test-update-" + runId,
		"key_ring_id":   "tf-test-update-" + runId,
		"crypto_key_id": "tf-test-update-" + runId,
	}

	// To store details of the resource to check if it's been replaced
	resourceDetails := make(map[string]string)
	resourceName := "google_vector_search_collection.example-collection"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// 1. Create with basic config (no CMEK)
			{
				Config: testAccVectorSearchCollection_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "My Awesome Collection"),
					resource.TestCheckResourceAttr(resourceName, "encryption_spec.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vector_schema.#", "1"),
					storeResourceDetails(&resourceDetails, resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection_id", "labels", "location", "terraform_labels"},
			},
			// 2. Update mutable fields - NO ForceNew
			{
				Config: testAccVectorSearchCollection_updated_mutable(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "My Updated Awesome Collection"),
					resource.TestCheckResourceAttr(resourceName, "description", "This collection stores important data - updated."),
					resource.TestCheckResourceAttr(resourceName, "labels.env", "prod"),
					resource.TestCheckResourceAttr(resourceName, "vector_schema.#", "2"),
					checkResourceNotInternallyReplaced(&resourceDetails, resourceName),
				),
			},
			// 3. Add CMEK - should force new resource (detect via create_time)
			{
				Config: testAccVectorSearchCollection_updated_mutable_cmek(context, "key1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "My Updated Awesome Collection"), // Should retain from config
					resource.TestCheckResourceAttr(resourceName, "encryption_spec.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_spec.0.crypto_key_name"),
					checkResourceInternallyReplaced(&resourceDetails, resourceName),
					storeResourceDetails(&resourceDetails, resourceName),
				),
			},
			// 4. Update display_name on CMEK resource - NO ForceNew
			{
				Config: testAccVectorSearchCollection_updated_mutable_cmek_rename(context, "key1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "CMEK Collection Renamed"),
					resource.TestCheckResourceAttr(resourceName, "encryption_spec.#", "1"),
					checkResourceNotInternallyReplaced(&resourceDetails, resourceName),
				),
			},
			// 5. Change CMEK key - should force new resource (detect via create_time)
			{
				Config: testAccVectorSearchCollection_updated_mutable_cmek(context, "key2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "encryption_spec.#", "1"),
					checkResourceInternallyReplaced(&resourceDetails, resourceName),
					storeResourceDetails(&resourceDetails, resourceName),
				),
			},
			// 6. Remove CMEK - should force new resource (detect via create_time)
			{
				Config: testAccVectorSearchCollection_updated_mutable(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "encryption_spec.#", "0"),
					checkResourceInternallyReplaced(&resourceDetails, resourceName),
					storeResourceDetails(&resourceDetails, resourceName),
				),
			},
		},
	})
}

// Helper function to store the resource ID and create_time
func storeResourceDetails(details *map[string]string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		(*details)["id"] = rs.Primary.ID
		val, ok := rs.Primary.Attributes["create_time"]
		if !ok {
			return fmt.Errorf("Attribute 'create_time' not found in state for %s", resourceName)
		}
		(*details)["create_time"] = val
		return nil
	}
}

// Helper function to check if the resource was internally replaced by comparing create_time
func checkResourceInternallyReplaced(oldDetails *map[string]string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		newID := rs.Primary.ID
		newCreateTime, ok := rs.Primary.Attributes["create_time"]
		if !ok {
			return fmt.Errorf("Attribute 'create_time' not found in state for %s", resourceName)
		}

		oldID := (*oldDetails)["id"]
		oldCreateTime := (*oldDetails)["create_time"]

		if newID == oldID && newCreateTime == oldCreateTime {
			return fmt.Errorf("Resource %s was not internally replaced, ID remained: %s, create_time remained: %s", resourceName, newID, newCreateTime)
		}
		if newCreateTime == oldCreateTime {
			return fmt.Errorf("Resource %s ID changed from %s to %s, but create_time unexpectedly remained: %s", resourceName, oldID, newID, newCreateTime)
		}
		return nil
	}
}

// Helper function to check if the resource was NOT internally replaced
func checkResourceNotInternallyReplaced(oldDetails *map[string]string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		newID := rs.Primary.ID
		newCreateTime, ok := rs.Primary.Attributes["create_time"]
		if !ok {
			return fmt.Errorf("Attribute 'create_time' not found in state for %s", resourceName)
		}

		oldID := (*oldDetails)["id"]
		oldCreateTime := (*oldDetails)["create_time"]

		if newID != oldID {
			return fmt.Errorf("Resource %s ID unexpectedly changed from %s to %s", resourceName, oldID, newID)
		}
		if newCreateTime != oldCreateTime {
			return fmt.Errorf("Resource %s was unexpectedly internally replaced, ID remained: %s, but create_time changed from %s to %s", resourceName, newID, oldCreateTime, newCreateTime)
		}
		return nil
	}
}

const dataSchemaBasic = `<<EOF
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string"
    },
    "plot": {
      "type": "string"
    }
  }
}
EOF`

const dataSchemaUpdated = `<<EOF
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string"
    },
    "plot": {
      "type": "string"
    },
    "year": {
      "type": "integer"
    }
  }
}
EOF`

const vectorSchemaBasic = `
  vector_schema {
    field_name = "text_embedding"
    dense_vector {
      dimensions = 768
      vertex_embedding_config {
        model_id   = "textembedding-gecko@003"
        task_type  = "RETRIEVAL_DOCUMENT"
        text_template = "Title: {title} ---- Plot: {plot}"
      }
    }
  }`

const vectorSchemaUpdated = `
  vector_schema {
    field_name = "text_embedding"
    dense_vector {
      dimensions = 768
      vertex_embedding_config {
        model_id   = "textembedding-gecko@003"
        task_type  = "RETRIEVAL_DOCUMENT"
        text_template = "Title: {title} ---- Plot: {plot}"
      }
    }
  }

  vector_schema {
    field_name = "sparse_embedding"
    sparse_vector {}
  }`

func testAccVectorSearchCollection_basic(context map[string]interface{}) string {
	return acctest.Nprintf(fmt.Sprintf(`
resource "google_vector_search_collection" "example-collection" {
  location      = "us-central1"
  collection_id = "%%{collection_id}"

  display_name = "My Awesome Collection"
  description  = "This collection stores important data."

  labels = {
    env  = "dev"
    team = "my-team"
  }

  data_schema = %s
  %s
  // No encryption_spec
}
`, dataSchemaBasic, vectorSchemaBasic), context)
}

func testAccVectorSearchCollection_updated_mutable(context map[string]interface{}) string {
	return acctest.Nprintf(fmt.Sprintf(`
resource "google_vector_search_collection" "example-collection" {
  location      = "us-central1"
  collection_id = "%%{collection_id}"

  display_name = "My Updated Awesome Collection"
  description  = "This collection stores important data - updated."

  labels = {
    env  = "prod" // Changed
    // team label removed
  }

  data_schema = %s
  %s
  // No encryption_spec
}
`, dataSchemaUpdated, vectorSchemaUpdated), context)
}

func cmekResources(context map[string]interface{}, keyName string) string {
	return acctest.Nprintf(fmt.Sprintf(`
data "google_project" "project" {}

resource "google_kms_key_ring" "key_ring" {
  name     = "%%{key_ring_id}"
  location = "us-central1"
	project = data.google_project.project.project_id
}

resource "google_kms_crypto_key" "%s" {
  name     = "%%{crypto_key_id}-%s"
  key_ring = google_kms_key_ring.key_ring.id
}

resource "google_kms_crypto_key_iam_member" "crypto_key_member_vs_sa_%s" {
  crypto_key_id = google_kms_crypto_key.%s.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-vectorsearch.iam.gserviceaccount.com"
}
`, keyName, keyName, keyName, keyName), context)
}

func testAccVectorSearchCollection_updated_mutable_cmek(context map[string]interface{}, keyName string) string {
	return fmt.Sprintf(`
%s

resource "google_vector_search_collection" "example-collection" {
  location      = "us-central1"
  collection_id = "%s"

  display_name = "My Updated Awesome Collection"
  description  = "This collection stores important data - updated."

  encryption_spec {
    crypto_key_name = google_kms_crypto_key.%s.id
  }

  labels = {
    env  = "prod"
  }

  data_schema = %s
  %s

  depends_on = [google_kms_crypto_key_iam_member.crypto_key_member_vs_sa_%s]
}
`, cmekResources(context, keyName), context["collection_id"], keyName, dataSchemaUpdated, vectorSchemaUpdated, keyName)
}

func testAccVectorSearchCollection_updated_mutable_cmek_rename(context map[string]interface{}, keyName string) string {
	return fmt.Sprintf(`
%s

resource "google_vector_search_collection" "example-collection" {
  location      = "us-central1"
  collection_id = "%s"

  display_name = "CMEK Collection Renamed" // Changed
  description  = "This collection stores important data - updated."

  encryption_spec {
    crypto_key_name = google_kms_crypto_key.%s.id
  }

  labels = {
    env  = "prod"
  }

  data_schema = %s
  %s

  depends_on = [google_kms_crypto_key_iam_member.crypto_key_member_vs_sa_%s]
}
`, cmekResources(context, keyName), context["collection_id"], keyName, dataSchemaUpdated, vectorSchemaUpdated, keyName)
}
