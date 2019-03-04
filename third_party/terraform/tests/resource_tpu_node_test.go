package google

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTpuNode_tpuNodeBUpdateTensorFlowVersion(t *testing.T) {
	t.Parallel()

	nodeId := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTpuNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTpuNode_tpuNodeTensorFlow(nodeId, "1.11"),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccTpuNode_tpuNodeTensorFlow(nodeId, "1.12"),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

// WARNING: cidr_block must not overlap with other existing TPU blocks
// Make sure if you change this value that it does not overlap with the
// autogenerated examples.
func testAccTpuNode_tpuNodeTensorFlow(nodeId, tensorFlowVer string) string {
	return fmt.Sprintf(`
resource "google_tpu_node" "tpu" {
  name           = "%s"
  zone           = "us-central1-b"

  accelerator_type   = "v3-8"
  tensorflow_version = "%s"
  cidr_block         = "10.1.0.0/29"
}
`, nodeId, tensorFlowVer)
}
