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

func (t *testing.T) {
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
				Config: testAccProjectIamMemberListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_project_iam_member.test", 1),
				},
			},
		},
	})
}

// test with optional filters
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

func testAccProjectIamMemberCreate(project, role, member string) string {
	return fmt.Sprintf(`
resource "google_project_iam_member" "test" {
  project = %q
  role    = %q
  member  = %q
}
`, project, role, member)
}

func testAccProjectIamMemberListQuery(project string) string {
	return fmt.Sprintf(`
list "google_project_iam_member" "test" {
  provider = google

  include_resource = true

  config {
    project = %q
  }
}
`, project)
}

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
