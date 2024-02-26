package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccApphubServiceProject_serviceProjectUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApphubServiceProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubServiceProject_serviceProjectBasicExample(context),
			},
			{
				ResourceName:            "google_apphub_service_project.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
			{
				Config: testAccApphubServiceProject_serviceProjectUpdate(context),
			},
			{
				ResourceName:            "google_apphub_service_project.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
		},
	})
}

func testAccApphubServiceProject_serviceProjectUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_service_project" "example" {
	service_project_attachment_id = google_project.service_project.project_id
}

resource "google_project" "service_project" {
	project_id ="tf-test-project-2%{random_suffix}"
	name = "Service Project New"
}
`, context)
}
