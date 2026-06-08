package resourcemanager_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFolderIamMemberList_basic(t *testing.T) {
	t.Parallel()

	folder := os.Getenv("GOOGLE_FOLDER")
	if folder == "" {
		t.Skip("GOOGLE_FOLDER must be set for this test")
	}

	if !strings.HasPrefix(folder, "folders/") {
		folder = "folders/" + folder
	}
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
				Config: testAccFolderIamMemberCreate(folder, role, member),
			},

			{
				Query:  true,
				Config: testAccFolderIamMemberListQuery(folder),
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

	folder := os.Getenv("GOOGLE_FOLDER")
	if folder == "" {
		t.Skip("GOOGLE_FOLDER must be set for this test")
	}

	if !strings.HasPrefix(folder, "folders/") {
		folder = "folders/" + folder
	}
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
				Config: testAccFolderIamMemberCreate(folder, role, member),
			},

			{
				Query:  true,
				Config: testAccFolderIamMemberListQueryWithFilters(folder, role, member),

				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_folder_iam_member.test", 1),
				},
			},
		},
	})
}

func testAccFolderIamMemberCreate(folder, role, member string) string {
	return fmt.Sprintf(`
resource "google_folder_iam_member" "test" {
	folder = "%s"
	role = "%s"
	member = "%s"
}
`, folder, role, member)
}

func testAccFolderIamMemberListQuery(folder string) string {
	return fmt.Sprintf(`

list "google_folder_iam_member" "test" {
	provider = google
	include_resource = true
	config {
		folder = "%s"
	}

}
`, folder)
}

func testAccFolderIamMemberListQueryWithFilters(folder, role, member string) string {
	return fmt.Sprintf(`

list "google_folder_iam_member" "test" {
	provider = google
	include_resource = true
	config {
		folder = "%s"
		role = "%s"
		member = "%s"
	}

}
`, folder, role, member)
}
