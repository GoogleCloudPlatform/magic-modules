// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappVolumeQuotaRule_netappVolumeQuotaRuleBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappVolumeQuotaRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappVolumeQuotaRule_netappVolumeQuotaRuleFull(context),
			},
			{
				ResourceName:            "google_netapp_volume_quota_rule.test_quotaRule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "volume_name"},
			},
			{
				Config: testAccNetappVolumeQuotaRule_netappVolumeQuotaRuleFull_update(context),
			},
			{
				ResourceName:            "google_netapp_volume_quota_rule.test_quotaRule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "volume_name"},
			},
		},
	})
}

func testAccNetappVolumeQuotaRule_netappVolumeQuotaRuleFull(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-test-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = 2048
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  location = google_netapp_storage_pool.default.location
  name = "tf-test-test-volume%{random_suffix}"
  capacity_gib = 100
  share_name = "tf-test-test-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
}

resource "google_netapp_volume_quota_rule" "test_quotaRule" {
  depends_on = [google_netapp_volume.default]
  location = google_netapp_volume.default.location
  volume_name = google_netapp_volume.default.name
  name = "testvolumequotaRule%{random_suffix}"
  description = "This is a test description"
  type = 'DEFAULT_USER_QUOTA'
  disk_limit_mib = 50
  target = "user1"
  labels = {
	key= "test"
	value= "quota_rule"
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccNetappVolumeQuotaRule_netappVolumeQuotaRuleFull_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-test-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = 2048
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  location = google_netapp_storage_pool.default.location
  name = "tf-test-test-volume%{random_suffix}"
  capacity_gib = 100
  share_name = "tf-test-test-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
}

resource "google_netapp_volume_quota_rule" "test_quotaRule" {
  depends_on = [google_netapp_volume.default]
  location = google_netapp_volume.default.location
  volume_name = google_netapp_volume.default.name
  name = "testvolumequotaRule%{random_suffix}"
  description = "This is a updated description"
  type = 'DEFAULT_USER_QUOTA'
  disk_limit_mib = 50
  target = "user1"
  labels = {
	key= "test"
	value= "quota_rule_update"
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}
