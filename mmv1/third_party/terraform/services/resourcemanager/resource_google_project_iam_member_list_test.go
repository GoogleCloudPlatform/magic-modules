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

func TestAccProjectIamMemberList_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// List resources require Terraform >= 1.14.0 (terraform query / .tfquery.hcl).
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIamMemberCreate(project, role, member),
			},

			{
				Query:  true,
				Config: testAccProjectIamMemberListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_project_iam_member.test", 1),
				},
			},
		},
	})
}

func TestAccProjectIamMemberList_filter(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
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
				Config: testAccProjectIamMemberCreate(project, role, member),
			},

			{
				Query:  true,
				Config: testAccProjectIamMemberListQueryWithFilters(project, role, member),

				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_project_iam_member.test", 1),
				},
			},
		},
	})
}

// Normal apply-mode config: only managed resources here (NO list blocks).
func testAccProjectIamMemberCreate(project, role, member string) string {
	return fmt.Sprintf(`
resource "google_project_iam_member" "test" {
  project = %q
  role    = %q
  member  = %q
}
`, project, role, member)
}

// Query-mode config: ONLY list blocks here. This gets treated like a .tfquery.hcl file
// when TestStep.Query = true (terraform query).
func testAccProjectIamMemberListQuery(project string) string {
	return fmt.Sprintf(`
list "google_project_iam_member" "test" {
  provider = google

  # include_resource allows result.resource.* fields to be present in query output
  include_resource = true

  config {
    project = %q
  }
}
`, project)
}

// Query-mode config with optional filters. Keep ONLY if your list schema supports them.
func testAccProjectIamMemberListQueryWithFilters(project, role, member string) string {
	return fmt.Sprintf(`
list "google_project_iam_member" "test" {
  provider = google
  include_resource = true

  config {
    project = %q
    role    = %q
    member  = %q
  }
}
`, project, role, member)
}
