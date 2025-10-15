package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryPackages_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryPackagesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_packages.this", "project"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_packages.this", "location"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_packages.this", "repository_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_packages.this", "packages.0.name"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryPackagesConfig = `
data "google_artifact_registry_packages" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
  filter        = "name=\"projects/go-containerregistry/locations/us/repositories/gcr.io/packages/gcrane\""
}
`
