package artifactregistry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryMavenArtifacts_basic(t *testing.T) {
	t.Parallel()

	// At the moment there are no public Maven artifacts available in Artifact Registry.
	// This test is skipped to avoid unnecessary failures.
	// As soon as there are public artifacts available, this test can be enabled by removing the skip and adjusting the configuration accordingly.
	t.Skip("No public Maven artifacts available in Artifact Registry")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryMavenArtifactsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "project"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "location"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "repository_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "maven_artifacts.0.artifact_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "maven_artifacts.0.group_id"),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_maven_artifacts.test", "maven_artifacts.0.name"),
				),
			},
		},
	})
}

const testAccDataSourceArtifactRegistryMavenArtifactsConfig = `
data "google_artifact_registry_maven_artifacts" "test" {
  project       = "example-project"
  location      = "us"
  repository_id = "example-repo"
}
`
