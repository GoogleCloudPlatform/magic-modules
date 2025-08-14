package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryPackage_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryPackageConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_artifact_registry_package.this", "name", "projects/go-containerregistry/locations/us/repositories/gcr.io/packages/gcrane"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryPackageConfig = `
data "google_artifact_registry_package" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
  name          = "gcrane"
}
`
