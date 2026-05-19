package vertexai_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIArtifact_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": strings.ToLower(acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIArtifact_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_artifact.artifact",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "region", "metadatastore", "artifact_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIArtifact_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_artifact.artifact",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "region", "metadatastore", "artifact_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIArtifact_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_metadata_store" "store" {
  provider      = google-beta
  name          = "tf-test-store-%{random_suffix}"
  description   = "Store to test Artifact"
  region        = "us-central1"
}

resource "google_vertex_ai_artifact" "artifact" {
  provider       = google-beta
  artifact_id    = "tf-test-artifact-%{random_suffix}"
  display_name   = "tf-test-artifact-display-%{random_suffix}"
  description    = "An artifact description"
  region         = "us-central1"
  metadatastore  = google_vertex_ai_metadata_store.store.name
  schema_title   = "system.Dataset"
  schema_version = "0.0.1"
  uri            = "https://example.com/dataset"
  metadata = {
    key = "value"
  }
  labels = {
    foo = "bar"
  }
}
`, context)
}

func testAccVertexAIArtifact_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_metadata_store" "store" {
  provider      = google-beta
  name          = "tf-test-store-%{random_suffix}"
  description   = "Store to test Artifact"
  region        = "us-central1"
}

resource "google_vertex_ai_artifact" "artifact" {
  provider       = google-beta
  artifact_id    = "tf-test-artifact-%{random_suffix}"
  display_name   = "tf-test-artifact-display-updated-%{random_suffix}"
  description    = "An artifact description updated"
  region         = "us-central1"
  metadatastore  = google_vertex_ai_metadata_store.store.name
  schema_title   = "system.Dataset"
  schema_version = "0.0.1"
  state          = "LIVE"
  uri            = "https://example.com/dataset-updated"
  metadata = {
    key = "value-updated"
  }
  labels = {
    foo = "bar2"
  }
}
`, context)
}
