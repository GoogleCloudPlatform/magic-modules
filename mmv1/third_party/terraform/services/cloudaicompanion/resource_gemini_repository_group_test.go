package gemini_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccGeminiRepositoryGroup_update(t *testing.T) {
	t.Parallel()

	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", "us-central1", "", "")
	context := map[string]interface{}{
		"project_id":               os.Getenv("GOOGLE_PROJECT"),
		"code_repository_index_id": codeRepositoryIndexId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroup_basic(context),
			},
			{
				ResourceName:            "google_gemini_repository_group.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"coderepositoryindex", "labels", "location", "repository_group_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiRepositoryGroup_update(context),
			},
			{
				ResourceName:            "google_gemini_repository_group.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"coderepositoryindex", "labels", "location", "repository_group_id", "terraform_labels"},
			},
		},
	})
}

func testAccGeminiRepositoryGroup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "test-repository-group-id1" 
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/tf-test-cloudaicompanion1/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "my_repository"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion1"
  disabled = false

  github_config {
    github_app = "DEVELOPER_CONNECT"
    app_installation_id = 54180648

    authorizer_credential {
      oauth_token_secret_version = "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1"
    }
  }
}
`, context)
}
func testAccGeminiRepositoryGroup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "test-repository-group-id2"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/tf-test-cloudaicompanion2/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1", "label2": "value2"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "my_repository"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion2"
  disabled = false

  github_config {
    github_app = "DEVELOPER_CONNECT"
    app_installation_id = 54180648

    authorizer_credential {
      oauth_token_secret_version = "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1"
    }
  }
}
`, context)
}
