package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryVersion_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryVersionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_artifact_registry_version.this", "name", "projects/go-containerregistry/locations/us/repositories/gcr.io/packages/gcrane/versions/sha256:c0cf52c2bd8c636bbf701c6c74c5ff819447d384dc957d52a52a668de63e8f5d"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryVersionConfig = `
data "google_artifact_registry_version" "this" {
  project       = "go-containerregistry"
  location      = "us"
  repository_id = "gcr.io"
  package_name  = "gcrane"
  version_name  = "sha256:c0cf52c2bd8c636bbf701c6c74c5ff819447d384dc957d52a52a668de63e8f5d"
}
`
