package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccUniverseDomain(t *testing.T) {
	// Skip VCR since this test can only run in specific test project.
	t.Skip()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUniverseDomain_basic_disk(),
			},
		},
	})
}

func testAccUniverseDomain_basic_disk() string {
	return fmt.Sprintf(`
provider "google" {
  universe_domain = "test"
}
	  
resource "google_compute_instance_template" "instance_template" {
  name = "demo-it"
  machine_type = "n1-standard-1"

// boot disk
  disk {
	disk_size_gb = 20
  }

  network_interface {
	network = "default"
  }
}
`)
}
