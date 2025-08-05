package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryDockerImages_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryDockerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "project"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "location"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "repository_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "docker_images.0.image_name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "docker_images.0.name"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_docker_images.this", "docker_images.0.self_link"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
// Currently, gcr.io does not provide a imageSizeBytes or buildTime field in the JSON response
const testAccDataSourceArtifactRegistryDockerImagesConfig = `
data "google_artifact_registry_docker_images" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
}
`
