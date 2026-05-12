package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccProjectIamMemberlist_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv(t)
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks:   projectIamMemberTerraformVersionChecks(),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIamMemberListBasic(project, role, member),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("listed_role", role),
					resource.TestCheckOutput("listed_memeber", member),
				),
			},
		},
	})
}

func TestAccProjectIamMemberList_filter(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv(t)
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks:   projectIamMemberTerraformVersionChecks(),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIamMemberListFilter(project, role, member),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("filtered_count", "1"),
					resource.TestCheckOutput("filtered_role", role),
					resource.TestCheckOutput("filtered_memeber", member),
				),
			},
		},
	})
}

func testAccProjectIamMemberListBasic(project, role, member string) string {
	return fmt.Sprintf(`	
resource "google_project_iam_member" "test" {
	project = %[1]q
	role    = %[2]q
	member  = %[3]q
}

list "google_project_iam_member" "test" {
	provider = google
	config {
		project = %[1]q
	}
	depends_on = [google_project_iam_member.test]
}
	output "listed_role" {
	value = one([
		for e in List.google_project_iam_member.test.results : r.resource.role
		if r.resource.role == %[2]q && r.resource.member == %[3]q
		])
	}
	
	output "listed_member" {
	value = one([
		for e in List.google_project_iam_member.test.results : r.resource.member
		if r.resource.role == %[2]q && r.resource.member == %[3]q
		])
	}

	`, project, role, member)
}

func testAccProjectIamMemberListFilter(project, role, member string) string {
	return fmt.Sprintf(`
resource "google_project_iam_member" "test" {
	project = %[1]q
	role    = %[2]q
	member  = %[3]q
}

list "google_project_iam_member" "test" {
	provider = google
	config {
		project = %[1]q
		role = %[2]q
		member = %[3]q
	}
	depends_on = [google_project_iam_member.test]
}

output "filtered_role" {
	value = one(list.google_project_iam_member.test.results).resource.role
	}

output "filtered_member" {
	value = one(list.google_project_iam_member.test.results).resource.member
	}
	`, project, role, member)
}
