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

	resourceName := "data.google_artifact_registry_docker_image.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryDockerImageConfig,
				Check: resource.ComposeTestCheckFunc(
					// Data source using a tag
					checkTaggedDataSources(resourceName+"Tag", "latest"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "image_size_bytes"),
					validateTimeStamps(resourceName+"Tag"),

					// Data source using a digest
					checkDigestDataSources(
						resourceName+"Digest",
						"projects/cloudrun/locations/us/repositories/container/dockerImages/hello@sha256:77cb9fbc6a667b8bfdbeca4c49e7703d825746eba53b736f0318bb7712828821",
						"us-docker.pkg.dev/cloudrun/container/hello@sha256:77cb9fbc6a667b8bfdbeca4c49e7703d825746eba53b736f0318bb7712828821",
					),
					resource.TestCheckResourceAttrSet(resourceName+"Digest", "image_size_bytes"),
					validateTimeStamps(resourceName+"Digest"),

					// url safe docker name using a tag
					checkTaggedDataSources(resourceName+"UrlTag", "latest"),

					// url safe docker name using a digest
					checkDigestDataSources(
						resourceName+"UrlDigest",
						"projects/go-containerregistry/locations/us/repositories/gcr.io/dockerImages/krane%2Fdebug@sha256:26903bf659994649af0b8ccb2675b76318b2bc3b2c85feea9a1f9d5b98eff363",
						"us-docker.pkg.dev/go-containerregistry/gcr.io/krane/debug@sha256:26903bf659994649af0b8ccb2675b76318b2bc3b2c85feea9a1f9d5b98eff363",
					),

					// Data source using no tag or digest
					resource.TestCheckResourceAttrSet(resourceName+"None", "project"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "repository"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "region"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "image"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "name"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "self_link"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
// Currently, gcr.io does not provide a imageSizeBytes or buildTime field in the JSON response
const testAccDataSourceArtifactRegistryDockerImageConfig = `
data "google_artifact_registry_docker_image" "testTag" {
	project    = "cloudrun"
	repository = "container"
	region     = "us"
	image      = "hello"
	tag        = "latest"
}

data "google_artifact_registry_docker_image" "testDigest" {
	project    = "cloudrun"
	repository = "container"
	region     = "us"
	image      = "hello"
	digest     = "sha256:77cb9fbc6a667b8bfdbeca4c49e7703d825746eba53b736f0318bb7712828821"
}

data "google_artifact_registry_docker_image" "testUrlTag" {
	project    = "go-containerregistry"
	repository = "gcr.io"
	region     = "us"
	image      = "krane/debug"
	tag        = "latest"
}

data "google_artifact_registry_docker_image" "testUrlDigest" {
	project    = "go-containerregistry"
	repository = "gcr.io"
	region     = "us"
	image      = "krane/debug"
	digest     = "sha256:26903bf659994649af0b8ccb2675b76318b2bc3b2c85feea9a1f9d5b98eff363"
}

data "google_artifact_registry_docker_image" "testNone" {
	project    = "go-containerregistry"
	repository = "gcr.io"
	region     = "us"
	image      = "crane"
}
`

func checkTaggedDataSources(resourceName string, expectedTag string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "project"),
		resource.TestCheckResourceAttrSet(resourceName, "repository"),
		resource.TestCheckResourceAttrSet(resourceName, "region"),
		resource.TestCheckResourceAttrSet(resourceName, "image"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "self_link"),
		resource.TestCheckTypeSetElemAttr(resourceName, "tags.*", expectedTag),
		resource.TestCheckResourceAttrSet(resourceName, "media_type"),
	)
}

func checkDigestDataSources(resourceName string, expectedName string, expectedSelfLink string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "project"),
		resource.TestCheckResourceAttrSet(resourceName, "repository"),
		resource.TestCheckResourceAttrSet(resourceName, "region"),
		resource.TestCheckResourceAttrSet(resourceName, "image"),
		resource.TestCheckResourceAttr(resourceName, "name", expectedName),
		resource.TestCheckResourceAttr(resourceName, "self_link", expectedSelfLink),
		resource.TestCheckResourceAttrSet(resourceName, "media_type"),
	)
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

func isRFC3339(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
