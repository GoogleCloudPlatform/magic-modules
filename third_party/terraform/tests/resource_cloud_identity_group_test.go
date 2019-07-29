package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudIdentityGroup_basic(t *testing.T) {
	t.Parallel()

	name := "tftest-cloudidentity-" + acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroup(name),
			},
			{
				ResourceName:      "google_cloud_identity_group.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIdentityGroup(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_identity_group" "default" {
  name          = "%s"
  display_name   = "%s"

  labels = {
  	name = "%s"
  	label_key = "Label-Value"
  }
}
`, name, name, name)
}
