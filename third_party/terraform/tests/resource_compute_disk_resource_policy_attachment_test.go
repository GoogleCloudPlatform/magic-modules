package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeDiskResourcePolicyAttachment_basic(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	policyName := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(10))
	policyName2 := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskResourcePolicyAttachmentExists(
						"google_compute_disk_resource_policy_attachment.foobar", "google_compute_resource_policy.foobar", policyName),
				),
			},
			{
				Config: testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskResourcePolicyAttachmentExists(
						"google_compute_disk_resource_policy_attachment.foobar", "google_compute_resource_policy.foobar", policyName2),
				),
			},
		},
	})
}

func testAccCheckComputeDiskResourcePolicyAttachmentExists(attachResName, policyResName, policyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attachRes, ok := s.RootModule().Resources[attachResName]
		if !ok {
			return fmt.Errorf("Not found: %s", attachResName)
		}

		if attachRes.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		policyRes, ok := s.RootModule().Resources[policyResName]
		if !ok {
			return fmt.Errorf("Not found: %s", policyResName)
		}

		if policyRes.Primary.Attributes["name"] != policyName {
			return fmt.Errorf("Resource Policy is incorrect")
		}
		if attachRes.Primary.Attributes["name"] != policyRes.Primary.Attributes["self_link"] {
			return fmt.Errorf("Resource Policy Attachment is incorrect")
		}

		return nil
	}
}

func testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	labels = {
		my-label = "my-label-value"
	}
}

resource "google_compute_resource_policy" "foobar" {
	name = "%s"
	region = "us-central1"
	snapshot_schedule_policy {
		schedule {
			daily_schedule {
				days_in_cycle = 1
				start_time = "04:00"
			}
		}
	}
}

resource "google_compute_disk_resource_policy_attachment" "foobar" {
	name = google_compute_resource_policy.foobar.self_link
  disk = google_compute_disk.foobar.name
	zone = "us-central1-a"
}`, diskName, policyName)
}
