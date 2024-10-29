package gemini_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"
)

func TestAccGeminiRepositoryGroupIamBindingGenerated(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", "")
	resourcePath := "projects/juliamat-sandbox/locations/us-central1/connections/testtest/gitRepositoryLinks/JumiDeluxe-testtest" // TODO: Change
	repositoryGroupId := acctest.BoostrapSharedRepositoryGroup(t, "basic", location, "", codeRepositoryIndexId, resourcePath)
	context := map[string]interface{}{
		"random_suffix":            acctest.RandString(t, 10),
		"role":                     "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index_id": codeRepositoryIndexId,
		"repository_group_id":      repositoryGroupId,
		"location":                 location,
		"project":                  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroupIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccGeminiRepositoryGroupIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGeminiRepositoryGroupIamMemberGenerated(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", "")
	resourcePath := "projects/juliamat-sandbox/locations/us-central1/connections/testtest/gitRepositoryLinks/JumiDeluxe-testtest" // TODO: Change
	repositoryGroupId := acctest.BoostrapSharedRepositoryGroup(t, "basic", location, "", codeRepositoryIndexId, resourcePath)
	context := map[string]interface{}{
		"random_suffix":            acctest.RandString(t, 10),
		"role":                     "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index_id": codeRepositoryIndexId,
		"repository_group_id":      repositoryGroupId,
		"location":                 location,
		"project":                  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccGeminiRepositoryGroupIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGeminiRepositoryGroupIamPolicyGenerated(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", "")
	resourcePath := "projects/juliamat-sandbox/locations/us-central1/connections/testtest/gitRepositoryLinks/JumiDeluxe-testtest" // TODO: Change
	repositoryGroupId := acctest.BoostrapSharedRepositoryGroup(t, "basic", location, "", codeRepositoryIndexId, resourcePath)
	context := map[string]interface{}{
		"random_suffix":            acctest.RandString(t, 10),
		"role":                     "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index_id": codeRepositoryIndexId,
		"repository_group_id":      repositoryGroupId,
		"location":                 location,
		"project":                  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroupIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_gemini_repository_group_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGeminiRepositoryGroupIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGeminiRepositoryGroupIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_member" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccGeminiRepositoryGroupIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  depends_on = [
    google_gemini_repository_group_iam_policy.foo
  ]
}
`, context)
}

func testAccGeminiRepositoryGroupIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_policy" "foo" {
}

resource "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccGeminiRepositoryGroupIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_binding" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccGeminiRepositoryGroupIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_binding" "foo" {
  project = "%{project}"
  location = "%{location}"
  coderepositoryindex = "%{code_repository_index_id}"
  repository_group_id = "%{repository_group_id}"
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
