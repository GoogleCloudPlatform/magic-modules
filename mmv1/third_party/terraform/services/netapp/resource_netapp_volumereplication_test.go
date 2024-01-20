// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappvolumereplication_netappVolumeReplicationCreateExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappvolumereplicationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappvolumereplication_netappVolumeReplicationCreateExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_volumereplication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappvolumereplication_netappVolumeReplicationCreateExample_update(context),
			},
			{
				ResourceName:            "google_netapp_volumereplication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappvolumereplication_netappVolumeReplicationCreateExample_stop(context),
			},
			{
				ResourceName:            "google_netapp_volumereplication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappvolumereplication_netappVolumeReplicationCreateExample_resume(context),
			},
			{
				ResourceName:            "google_netapp_volumereplication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "labels", "terraform_labels"},
			},
		},
	})
}

// Basic replication
func testAccNetappvolumereplication_netappVolumeReplicationCreateExample_basic(context map[string]interface{}) string {
	var result string = acctest.Nprintf(`

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
}

resource "google_netapp_volumereplication" "my-replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "EVERY_10_MINUTES"
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same makes
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
  }
}
`, context)
	// Give mirror some time to reach mirror_state==MIRRORED state
	time.Sleep(120 * time.Second)
	return result
}

// Update parameters
func testAccNetappvolumereplication_netappVolumeReplicationCreateExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
}

resource "google_netapp_volumereplication" "my-replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "HOURLY"
  description 		   = "This is a replication resource"
  labels {
	"foo": "bar",
  }
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same makes
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
  }
  replication_enabled = true
  delete_destination_volume = true
  force_stopping = true
}
`, context)
}

// Stop replication
func testAccNetappvolumereplication_netappVolumeReplicationCreateExample_stop(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
}

resource "google_netapp_volumereplication" "my-replication" {
	depends_on           = [google_netapp_volume.source_volume]
	location             = google_netapp_volume.source_volume.location
	volume_name          = google_netapp_volume.source_volume.name
	name                 = "tf-test-test-replication%{random_suffix}"
	replication_schedule = "HOURLY"
	description 		   = "This is a replication resource"
	labels {
	  "foo": "bar",
	}
	destination_volume_parameters {
	  storage_pool = google_netapp_storage_pool.destination_pool.id
	  volume_id    = "tf-test-destination-volume%{random_suffix}"
	  # Keeping the share_name of source and destination the same makes
	  # simplifies implementing client failover concepts
	  share_name  = "tf-test-source-volume%{random_suffix}"
	  description = "This is a replicated volume"
	}
	replication_enabled = false
	delete_destination_volume = true
	force_stopping = true
`, context)
}

// resume replication
func testAccNetappvolumereplication_netappVolumeReplicationCreateExample_resume(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
}

resource "google_netapp_volumereplication" "my-replication" {
	depends_on           = [google_netapp_volume.source_volume]
	location             = google_netapp_volume.source_volume.location
	volume_name          = google_netapp_volume.source_volume.name
	name                 = "tf-test-test-replication%{random_suffix}"
	replication_schedule = "HOURLY"
	description 		   = "This is a replication resource"
	labels {
	  "foo": "bar",
	}
	destination_volume_parameters {
	  storage_pool = google_netapp_storage_pool.destination_pool.id
	  volume_id    = "tf-test-destination-volume%{random_suffix}"
	  # Keeping the share_name of source and destination the same makes
	  # simplifies implementing client failover concepts
	  share_name  = "tf-test-source-volume%{random_suffix}"
	  description = "This is a replicated volume"
	}
	replication_enabled = true
	delete_destination_volume = true
	force_stopping = true
`, context)
}
