package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceGoogleComputeInterconnectAttachments_basic(t *testing.T) {
	t.Parallel()

	// This is not the default test region, so the list with no region argument can only find the attachment by listing every region.
	const region = "us-east1"
	project := envvar.GetTestProjectFromEnv()
	suffix := acctest.RandString(t, 10)
	attachmentName := fmt.Sprintf("tf-test-attachments-ds-%s", suffix)
	// The shared test project holds other attachments, so there is no count check, and both lists must hold this one.
	want := map[string]string{
		"name":                     attachmentName,
		"region":                   region,
		"type":                     "PARTNER",
		"edge_availability_domain": "AVAILABILITY_DOMAIN_1",
		"admin_enabled":            "true",
		"mtu":                      "1500",
		"labels.mykey":             "myvalue",
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.google_compute_interconnect_attachments.region", "id", fmt.Sprintf("projects/%s/regions/%s/interconnectAttachments", project, region)),
		resource.TestCheckResourceAttr("data.google_compute_interconnect_attachments.all", "id", fmt.Sprintf("projects/%s/aggregated/interconnectAttachments", project)),
	}
	for _, dataSource := range []string{"data.google_compute_interconnect_attachments.region", "data.google_compute_interconnect_attachments.all"} {
		checks = append(checks,
			resource.TestCheckTypeSetElemNestedAttrs(dataSource, "interconnect_attachments.*", want),
			resource.TestCheckTypeSetElemAttrPair(dataSource, "interconnect_attachments.*.self_link", "google_compute_interconnect_attachment.foobar", "self_link"),
			resource.TestCheckTypeSetElemAttrPair(dataSource, "interconnect_attachments.*.router", "google_compute_router.foobar", "self_link"),
			resource.TestCheckTypeSetElemAttrPair(dataSource, "interconnect_attachments.*.pairing_key", "google_compute_interconnect_attachment.foobar", "pairing_key"),
		)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInterconnectAttachmentsConfig(suffix, attachmentName, region),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// testAccDataSourceGoogleComputeInterconnectAttachmentsConfig creates a partner attachment, which needs no physical interconnect.
func testAccDataSourceGoogleComputeInterconnectAttachmentsConfig(suffix, attachmentName, region string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-attachments-ds-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_router" "foobar" {
  name    = "tf-test-attachments-ds-%[1]s"
  network = google_compute_network.foobar.name
  region  = "%[3]s"
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "foobar" {
  name                     = "%[2]s"
  region                   = "%[3]s"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.foobar.id
  mtu                      = 1500
  labels                   = { mykey = "myvalue" }
}

data "google_compute_interconnect_attachments" "region" {
  project = google_compute_interconnect_attachment.foobar.project
  region  = google_compute_interconnect_attachment.foobar.region
}

data "google_compute_interconnect_attachments" "all" {
  project = google_compute_interconnect_attachment.foobar.project
}
`, suffix, attachmentName, region)
}
