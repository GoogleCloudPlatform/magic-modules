package resourcemanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectIamMemberRemove_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	role := "roles/editor"
	members:= "default-gce-sa@example.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectIamCustomRoleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamMemberRemove_basic(role, project, members),
			},
		},
	})
}

func testAccCheckGoogleProjectIamMemberRemove_basic(roleId, project, members string) string {
	return fmt.Sprintf(`
resource "google_project_iam_member_remove" "foo" {
  role     = "%s"
  project  = "%s"
  member  = "%s"
}
`, roleId, project, members)
}
