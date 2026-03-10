package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryVersions_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryVersionsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "project"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "location"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "repository_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "package_name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "versions.0.name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "versions.0.create_time"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_versions.this", "versions.0.update_time"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryVersionsConfig = `
data "google_artifact_registry_versions" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
  package_name  = "gcrane"
  filter        = "name=\"projects/go-containerregistry/locations/us/repositories/gcr.io/packages/gcrane/versions/*:b*\""
  view          = "FULL"
}
`
