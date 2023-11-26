package artifactregistry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryDockerImage(t *testing.T) {
	t.Parallel()

	resourceName := "data.artifactregistry_docker_image.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryDockerImageConfig,
				Check: resource.ComposeTestCheckFunc(
					// Data source using a tag
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "project"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "repository"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "region"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "image"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "name"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "self_link"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "tags"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "image_size_bytes"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "media_type"),
					validateTimeStamps(resourceName+"Tag"),

					// Data source using a digest
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "project"),
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "repository"),
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "region"),
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "image"),
					resource.TestCheckResourceAttr(resourceName+"Digest", "name", "projects/go-containerregistry/locations/us/repositories/gcr.io/dockerImages/crane@sha256:0f1cfc0f8c87eb871b4c6f5c4b80f89fa912986369b1e3313a5e808214270bb3"),
					resource.TestCheckResourceAttr(resourceName+"Digest", "self_link", "us-docker.pkg.dev/go-containerregistry/gcr.io/crane@sha256:0f1cfc0f8c87eb871b4c6f5c4b80f89fa912986369b1e3313a5e808214270bb3"),
					// tags may become an empty list in the future
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "image_size_bytes"),
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "media_type"),
					validateTimeStamps(resourceName+"Digest"),
				),
			},
		},
	})
}

// Test the data source against the public AR repo
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
const testAccDataSourceArtifactRegistryDockerImageConfig = `
data "google_artifact_registry_docker_image" "testTag" {
  project    = "go-containerregistry"
  repository = "gcr.io"
  region     = "us"
  image      = "crane"
  tag        = "latest"
}

data "google_artifact_registry_docker_image" "testDigest" {
	project    = "go-containerregistry"
	repository = "gcr.io"
	region     = "us"
	image      = "crane"
	digest     = "sha256:0f1cfc0f8c87eb871b4c6f5c4b80f89fa912986369b1e3313a5e808214270bb3"
  }
`

func isRFC3339(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

func validateTimeStamps(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// check that the timestamps are RFC3339
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		if !isRFC3339(ds.Primary.Attributes["upload_time"]) {
			return fmt.Errorf("upload_time is not RFC3339: %s", ds.Primary.Attributes["upload_time"])
		}
		if !isRFC3339(ds.Primary.Attributes["build_time"]) {
			return fmt.Errorf("build_time is not RFC3339: %s", ds.Primary.Attributes["build_time"])
		}
		if !isRFC3339(ds.Primary.Attributes["update_time"]) {
			return fmt.Errorf("update_time is not RFC3339: %s", ds.Primary.Attributes["update_time"])
		}

		return nil
	}
}
