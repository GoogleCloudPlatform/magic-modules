package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceSpannerInstance_basic(t *testing.T) {
	// Randomness from spanner instance
	google.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": google.RandString(t, 10),
	}

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpannerInstanceBasic(context),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceState("data.google_spanner_instance.foo", "google_spanner_instance.bar"),
				),
			},
		},
	})
}

func testAccDataSourceSpannerInstanceBasic(context map[string]interface{}) string {
	return google.Nprintf(`
resource "google_spanner_instance" "bar" {
	config       = "regional-us-central1"
	display_name = "Test Spanner Instance"
	num_nodes    = 2
	labels = {
		"foo" = "bar"
	}
}

data "google_spanner_instance" "foo" {
	name = google_spanner_instance.bar.name
}
`, context)
}
