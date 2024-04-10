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

	org := envvar.GetTestOrgFromEnv(t)
	members := "user:gterraformtest7@gmail.com"
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectIamCustomRoleDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamMemberRemove_basic(random_suffix, org, members),
			},
		},
	})
}

func testAccCheckGoogleProjectIamMemberRemove_basic(random_suffix, org, members string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "tf-test-%s"
  name       = "tf-test-%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "bar" {
  project = google_project.project.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/editor"
}

resource "time_sleep" "wait_20s" {
  depends_on = [google_project_iam_binding.bar]
  create_duration = "20s"
}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.project.project_id
  member  = "%s"
  depends_on = [time_sleep.wait_20s]
}
`, random_suffix, random_suffix, org, members)
}
