package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryTags_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "project"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "location"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "repository_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "package_name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_tags.this", "tags.0.version"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryTagsConfig = `
data "google_artifact_registry_tags" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
  package_name  = "gcrane"
  # Filter doesn't work with gcr.io
  # filter        = "name=\"projects/go-containerregistry/locations/us/repositories/gcr.io/packages/gcrane/tags/latest\""
}
`
