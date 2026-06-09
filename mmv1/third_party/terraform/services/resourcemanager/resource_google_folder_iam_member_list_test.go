package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFolderIamMemberList_basic(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(t, 10)
	org := envvar.GetTestOrgFromEnv(t)
	parent := "organizations/" + org
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFolderIamMemberCreate(folderDisplayName, parent, role, member),
			},

			{
				Query:  true,
				Config: testAccFolderIamMemberListQuery(parent, folderDisplayName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_folder_iam_member.test", 1),
				},
			},
		},
	})
}

// test with optional filters
func TestAccFolderIamMemberList_filter(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(t, 10)
	org := envvar.GetTestOrgFromEnv(t)
	parent := "organizations/" + org
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFolderIamMemberCreate(folderDisplayName, parent, role, member),
			},

			{
				Query:  true,
				Config: testAccFolderIamMemberListQueryWithFilters(parent, folderDisplayName, role, member),

				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_folder_iam_member.test", 1),
				},
			},
		},
	})
}

func testAccFolderIamMemberCreate(folderDisplayName, parent, role, member string) string {
	return fmt.Sprintf(`
resource "google_folder" "test" {
  display_name = "%s"
  parent  = "%s"
  deletion_protection  = false
}

resource "google_folder_iam_member" "test" {
	folder = google_folder.test.name
	role = "%s"
	member = "%s"
}
`, folderDisplayName, parent, role, member)
}

func testAccFolderIamMemberListQuery(parent, folderDisplayName string) string {
	return fmt.Sprintf(`
data "google_folders" "under_parent"{
	parent_id = "%s"
}

locals {
	folder_name = one([
		for f in data.google_folders.under_parent.folders :
		f.name if f.display_name == "%s"
	])
}

list "google_folder_iam_member" "test" {
	provider = google
	include_resources = true
	config {
		folder = local.folder_name
	}

}
`, parent, folderDisplayName)
}

func testAccFolderIamMemberListQueryWithFilters(parent, folderDisplayName, role, member string) string {
	return fmt.Sprintf(`
data "google_folders" "under_parent"{
	parent_id = "%s"
}

locals {
	folder_name = one([
		for f in data.google_folders.under_parent.folders :
		f.name if f.display_name == "%s"
	])
}

list "google_folder_iam_member" "test" {
	provider = google
	include_resources = true
	config {
		folder = google_folder.test.name
		role = "%s"
		member = "%s"
	}

}
`, parent, folderDisplayName, role, member)
}
