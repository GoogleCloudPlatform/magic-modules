package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatastreamConnectionProfile_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_update(context),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location"},
			},
			{
				Config: testAccDatastreamConnectionProfile_update2(context, true),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamConnectionProfile_update2(context, false),
			},
		},
	})
}

func testAccDatastreamConnectionProfile_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%{random_suffix}"

	gcs_profile {
		bucket    = "my-bucket"
		root_path = "/path"
	}
	lifecycle {
		prevent_destroy = true
	}
}
`, context)
}


func testAccDatastreamConnectionProfile_update2(context map[string]interface{}, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%s"

	gcs_profile {
		bucket    = "my-other-bucket"
		root_path = "/path"
	}
	%s
}
`, context["random_suffix"], lifecycleBlock)
}
