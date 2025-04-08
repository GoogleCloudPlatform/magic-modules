package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleProjectIamCustomRoles_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamCustomRolesConfig(project, roleId),
				Check: resource.ComposeTestCheckFunc(
					// We can't guarantee no project won't have our custom role as first element, so we'll check set-ness rather than correctness
					resource.TestCheckResourceAttrSet("data.google_project_iam_custom_roles.this", "roles.0.id"),
					resource.TestCheckResourceAttrSet("data.google_project_iam_custom_roles.this", "roles.0.name"),
					resource.TestCheckResourceAttrSet("data.google_project_iam_custom_roles.this", "roles.0.role_id"),
					resource.TestCheckResourceAttrSet("data.google_project_iam_custom_roles.this", "roles.0.stage"),
					resource.TestCheckResourceAttrSet("data.google_project_iam_custom_roles.this", "roles.0.title"),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectIamCustomRolesConfig(project string, roleId string) string {
	return fmt.Sprintf(`
locals {
  project = "%s"
  role_id = "%s"
}
resource "google_project_iam_custom_role" "this" {
  project = local.project
  role_id = local.role_id
  title   = "Terraform Test"
  permissions = [
	"iam.roles.create",
	"iam.roles.delete",
    "iam.roles.list",
  ]
}
data "google_project_iam_custom_roles" "this" {
  project = google_project_iam_custom_role.this.project
}
`, project, roleId)
}
