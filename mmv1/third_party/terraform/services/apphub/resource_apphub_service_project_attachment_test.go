package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApphubServiceProjectAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentBasicExample(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateBasic(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
		},
	})
}

func testAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_service_project_attachment" "example" {
	service_project_attachment_id = google_project.service_project_new.project_id
}

resource "google_project" "service_project_new" {
	project_id ="tf-test-project-2%{random_suffix}"
	name = "Service Project New"
}
`, context)
}

func TestAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateFull(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApphubServiceProjectAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentFullExample(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateFull(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
		},
	})
}

func testAccApphubServiceProjectAttachment_serviceProjectAttachmentUpdateFull(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_service_project_attachment" "example2" {
	service_project_attachment_id = google_project.service_project_full_new.project_id
	service_project = google_project.service_project_full_new.project_id
}

resource "google_project" "service_project_full_new" {
	project_id ="tf-test-project-2%{random_suffix}"
	name = "Service Project Full New"
}
`, context)
}
