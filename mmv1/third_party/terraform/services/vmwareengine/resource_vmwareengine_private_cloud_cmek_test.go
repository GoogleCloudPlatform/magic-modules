package vmwareengine_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVmwareenginePrivateCloud_cmek_validation(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":               "me-west1",
		"random_suffix":        acctest.RandString(t, 10),
		"vmwareengine_project": os.Getenv("GOOGLE_VMWAREENGINE_PROJECT"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Step 1: CMEK without kms_key_name -> Should fail plan
			{
				Config:      testVmwareenginePrivateCloudCmekMissingKeyConfig(context),
				ExpectError: regexp.MustCompile("encryption_config.kms_key_name must be set when encryption_config.type is CMEK"),
			},
			// Step 2: GMEK with kms_key_name -> Should fail plan
			{
				Config:      testVmwareenginePrivateCloudGmekWithKeyConfig(context),
				ExpectError: regexp.MustCompile("encryption_config.kms_key_name cannot be set when encryption_config.type is GMEK"),
			},
		},
	})
}

func testVmwareenginePrivateCloudCmekMissingKeyConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project     = "%{vmwareengine_project}"
  name        = "tf-test-nw-%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project     = "%{vmwareengine_project}"
  name        = "tf-test-pc-%{random_suffix}"
  location    = "%{region}-b"
  type        = "TIME_LIMITED"

  deletion_delay_hours = 0
  send_deletion_delay_hours_if_zero = true

  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }

  management_cluster {
    cluster_id = "sample-cluster-%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
      custom_core_count = 32
    }
  }

  encryption_config {
    type = "CMEK"
    # kms_key_name is missing
  }
}
`, context)
}

func testVmwareenginePrivateCloudGmekWithKeyConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project     = "%{vmwareengine_project}"
  name        = "tf-test-nw-%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project     = "%{vmwareengine_project}"
  name        = "tf-test-pc-%{random_suffix}"
  location    = "%{region}-b"
  type        = "TIME_LIMITED"

  deletion_delay_hours = 0
  send_deletion_delay_hours_if_zero = true

  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }

  management_cluster {
    cluster_id = "sample-cluster-%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
      custom_core_count = 32
    }
  }

  encryption_config {
    type         = "GMEK"
    kms_key_name = "projects/my-project/locations/us-central1/keyRings/my-ring/cryptoKeys/my-key"
  }
}
`, context)
}
