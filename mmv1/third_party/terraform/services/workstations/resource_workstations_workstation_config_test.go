package workstations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/tags"
	_ "github.com/hashicorp/terraform-provider-google/google/services/workstations"
)

func TestAccWorkstationsWorkstationConfig_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"key_short_name":           "tf-test-key-" + acctest.RandString(t, 10),
		"org_id":                   envvar.GetTestOrgFromEnv(t),
		"value_short_name":         "tf-test-value-" + acctest.RandString(t, 10),
		"workstation_cluster_name": "tf-test-workstation-cluster" + randomSuffix,
		"workstation_config_name":  "tf-test-workstation-config" + randomSuffix,
		"random_suffix":            randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationConfig_workstationConfigBasicExample(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationConfig_update(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationConfig_updateStep3(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkstationsWorkstationConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "tag_key1" {
  parent     = "organizations/%{org_id}"
  short_name = "%{key_short_name}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent     = google_tags_tag_key.tag_key1.id
  short_name = "%{value_short_name}"
}

resource "google_compute_network" "default" {
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-west4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id   = "%{workstation_config_name}"
  workstation_cluster_id  = google_workstations_workstation_cluster.default.workstation_cluster_id
  location                = "us-west4"
  display_name            = "Display Name N %{random_suffix}"
  disable_tcp_connections = true

  readiness_checks {
    path = "/"
    port = 80
  }

  allowed_ports {
    first = 80
    last  = 80
  }

  container {
    image       = "us-central1-docker.pkg.dev/cloud-workstations-images/predefined/code-oss:latest"
    working_dir = "/test"
    env = {
      NAME = "FOO"
      BABE = "bar"
    }
  }

  host {
    gce_instance {
      machine_type                 = "n1-standard-4"
      boot_disk_size_gb            = 35
      disable_public_ip_addresses  = true
      enable_nested_virtualization = true
      instance_metadata = {
        "startup-key"  = "startup-value"
        "custom-field" = "custom-value"
      }
    }
  }

  labels = {
    foo = "bar"
  }

  max_usable_workstations = 2

  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

func testAccWorkstationsWorkstationConfig_updateStep3(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "tag_key1" {
  parent     = "organizations/%{org_id}"
  short_name = "%{key_short_name}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent     = google_tags_tag_key.tag_key1.id
  short_name = "%{value_short_name}"
}

resource "google_compute_network" "default" {
  name                    = "%{workstation_cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%{workstation_cluster_name}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "%{workstation_cluster_name}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-west4"

  labels = {
    "label" = "key"
  }

  annotations = {
    label-one = "value-one"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id   = "%{workstation_config_name}"
  workstation_cluster_id  = google_workstations_workstation_cluster.default.workstation_cluster_id
  location                = "us-west4"
  display_name            = "Display Name N2 %{random_suffix}"
  idle_timeout            = "600s"
  running_timeout         = "21600s"
  replica_zones           = ["us-west4-a", "us-west4-b"]
  disable_tcp_connections = false

  annotations = {
    label-one = "value-one"
  }

  labels = {
    "label" = "key"
  }

  max_usable_workstations = 1

  readiness_checks {
    path = "/healthz"
    port = 8080
  }

  allowed_ports {
    first = 80
    last  = 80
  }

  allowed_ports {
    first = 22
    last  = 22
  }

  allowed_ports {
    first = 1024
    last  = 65535
  }

  container {
    image       = "us-central1-docker.pkg.dev/cloud-workstations-images/predefined/code-oss:latest"
    working_dir = null
  }

  host {
    gce_instance {
      machine_type                = "e2-standard-4"
      boot_disk_size_gb           = 35
      disable_public_ip_addresses = true
      disable_ssh                 = false
      vm_tags = {
        (google_tags_tag_key.tag_key1.id) = google_tags_tag_value.tag_value1.id
      }
      instance_metadata = {
        "startup-key" = "startup-value-changed"
        "second-key"  = "second-value"
      }
    }
  }
}
`, context)
}

func TestAccWorkstationsWorkstationConfig_updateHostDetails(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationConfig_updateHostDetailsDefault(context),
			},
			{
				ResourceName:      "google_workstations_workstation_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				// enable_audit_agent is ignore_read in WorkstationConfig.yaml because the API does not return it
				ImportStateVerifyIgnore: []string{"enable_audit_agent", "etag"},
			},
			{
				Config: testAccWorkstationsWorkstationConfig_updateHostDetailsUpdated(context),
			},
			{
				ResourceName:      "google_workstations_workstation_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				// enable_audit_agent is ignore_read in WorkstationConfig.yaml because the API does not return it
				ImportStateVerifyIgnore: []string{"enable_audit_agent", "etag"},
			},
			{
				Config: testAccWorkstationsWorkstationConfig_updateHostDetailsUnsetInstanceConfigs(context),
			},
			{
				ResourceName:      "google_workstations_workstation_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				// enable_audit_agent is ignore_read in WorkstationConfig.yaml because the API does not return it
				ImportStateVerifyIgnore: []string{"enable_audit_agent", "etag"},
			},
		},
	})
}

func testAccWorkstationsWorkstationConfig_updateHostDetailsDefault(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "europe-west4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "europe-west4"
}

resource "google_service_account" "default" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Service Account"
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "europe-west4"

  enable_audit_agent = true

  host {
    gce_instance {
      machine_type                = "n1-standard-2"
      boot_disk_size_gb           = 35
      pool_size                   = 0
      service_account             = google_service_account.default.email
      service_account_scopes      = ["https://www.googleapis.com/auth/cloud-platform"]
      disable_public_ip_addresses = false

      accelerators {
        type  = "nvidia-tesla-t4"
        count = 1
      }

      shielded_instance_config {
        enable_secure_boot          = false
        enable_vtpm                 = false
        enable_integrity_monitoring = false
      }

      confidential_instance_config {
        enable_confidential_compute = false
      }
    }
  }
}
`, context)
}

func testAccWorkstationsWorkstationConfig_updateHostDetailsUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "europe-west4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "europe-west4"
}

# No longer explicitly used in google_workstations_workstation_config resource block below, but the
# service account needs to keep existing to allow the field to default from the API without error
resource "google_service_account" "default" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Service Account"
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "europe-west4"

  enable_audit_agent = true

  host {
    gce_instance {
      machine_type                = "n2d-standard-2"
      boot_disk_size_gb           = 35
      pool_size                   = 1
      service_account_scopes      = ["https://www.googleapis.com/auth/cloud-platform"]
      disable_public_ip_addresses = true
      tags                        = ["foo", "bar"]

      shielded_instance_config {
        enable_secure_boot          = true
        enable_vtpm                 = true
        enable_integrity_monitoring = true
      }

      confidential_instance_config {
        enable_confidential_compute = true
      }

      boost_configs {
        id           = "boost-1"
        machine_type = "n2d-standard-2"
      }
    }
  }
}
`, context)
}

func testAccWorkstationsWorkstationConfig_updateHostDetailsUnsetInstanceConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "europe-west4"
  network       = google_compute_network.default.name
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "europe-west4"
}

# No longer explicitly used in google_workstations_workstation_config resource block below, but the
# service account needs to keep existing to allow the field to default from the API without error
resource "google_service_account" "default" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Service Account"
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "europe-west4"

  enable_audit_agent = false

  host {
    gce_instance {
      machine_type                = "n2d-standard-2"
      boot_disk_size_gb           = 35
      pool_size                   = 1
      service_account_scopes      = ["https://www.googleapis.com/auth/cloud-platform"]
      disable_public_ip_addresses = true
      tags                        = ["foo", "bar"]

      shielded_instance_config {}
      confidential_instance_config {}
    }
  }
}
`, context)
}

func TestAccWorkstationsWorkstationConfig_directories(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationConfig_directories(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "terraform_labels"},
			},
			{
				ResourceName:            "google_workstations_workstation_config.source_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag"},
			},
			{
				ResourceName:            "google_workstations_workstation_config.ephemeral_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "terraform_labels"},
			},
			{
				ResourceName:            "google_workstations_workstation_config.ephemeral_image",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkstationsWorkstationConfig_directoriesUpdated(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "terraform_labels"},
			},
			{
				ResourceName:            "google_workstations_workstation_config.source_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag"},
			},
		},
	})
}

func testAccWorkstationsWorkstationConfig_directories(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west4"
  network       = google_compute_network.default.name
}

resource "google_compute_disk" "test_source_disk" {
  name = "tf-test-workstation-source-disk%{random_suffix}"
  size = 10
  type = "pd-ssd"
  zone = "us-west4-a"
}

resource "google_compute_snapshot" "test_source_snapshot" {
  name        = "tf-test-workstation-source-snapshot%{random_suffix}"
  source_disk = google_compute_disk.test_source_disk.name
  zone        = "us-west4-a"
}

resource "google_compute_image" "test_source_image" {
  name              = "tf-test-workstation-source-image%{random_suffix}"
  source_disk       = google_compute_disk.test_source_disk.self_link
  storage_locations = ["us-west4"]
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-west4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  persistent_directories {
    mount_path = "/home"
  }

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "source_snapshot" {
  workstation_config_id  = "tf-test-ws-cfg-snap%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  persistent_directories {
    mount_path = "/home"

    gce_pd {
      source_snapshot = google_compute_snapshot.test_source_snapshot.id
      reclaim_policy  = "DELETE"
    }
  }
}

resource "google_workstations_workstation_config" "ephemeral_snapshot" {
  workstation_config_id  = "tf-test-ws-cfg-eph-snap%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  ephemeral_directories {
    mount_path = "/cache"
    gce_pd {
      source_snapshot = google_compute_snapshot.test_source_snapshot.id
      read_only       = true
    }
  }

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "ephemeral_image" {
  workstation_config_id  = "tf-test-ws-cfg-eph-img%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  ephemeral_directories {
    mount_path = "/cache"
    gce_pd {
      disk_type    = "pd-standard"
      source_image = google_compute_image.test_source_image.id
    }
  }

  labels = {
    foo = "bar"
  }
}
`, context)
}

func testAccWorkstationsWorkstationConfig_directoriesUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west4"
  network       = google_compute_network.default.name
}

resource "google_compute_disk" "test_source_disk" {
  name = "tf-test-workstation-source-disk%{random_suffix}"
  size = 10
  type = "pd-ssd"
  zone = "us-west4-a"
}

resource "google_compute_snapshot" "test_source_snapshot" {
  name        = "tf-test-workstation-source-snapshot%{random_suffix}"
  source_disk = google_compute_disk.test_source_disk.name
  zone        = "us-west4-a"
}

resource "google_compute_snapshot" "test_source_snapshot2" {
  name        = "tf-test-workstation-source-snapshot2%{random_suffix}"
  source_disk = google_compute_disk.test_source_disk.name
  zone        = "us-west4-a"
}

resource "google_compute_image" "test_source_image" {
  name              = "tf-test-workstation-source-image%{random_suffix}"
  source_disk       = google_compute_disk.test_source_disk.self_link
  storage_locations = ["us-west4"]
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-west4"

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  persistent_directories {
    mount_path = "/home"

    gce_pd {
      size_gb        = 200
      fs_type        = "ext4"
      disk_type      = "pd-standard"
      reclaim_policy = "DELETE"
    }
  }

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "source_snapshot" {
  workstation_config_id  = "tf-test-ws-cfg-snap%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  persistent_directories {
    mount_path = "/home"

    gce_pd {
      source_snapshot = google_compute_snapshot.test_source_snapshot2.id
      reclaim_policy  = "RETAIN"
    }
  }
}

resource "google_workstations_workstation_config" "ephemeral_snapshot" {
  workstation_config_id  = "tf-test-ws-cfg-eph-snap%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  ephemeral_directories {
    mount_path = "/cache"
    gce_pd {
      source_snapshot = google_compute_snapshot.test_source_snapshot.id
      read_only       = true
    }
  }

  labels = {
    foo = "bar"
  }
}

resource "google_workstations_workstation_config" "ephemeral_image" {
  workstation_config_id  = "tf-test-ws-cfg-eph-img%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-west4"

  ephemeral_directories {
    mount_path = "/cache"
    gce_pd {
      disk_type    = "pd-standard"
      source_image = google_compute_image.test_source_image.id
    }
  }

  labels = {
    foo = "bar"
  }
}
`, context)
}

func TestAccWorkstationsWorkstationConfig_hyperdisk(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkstationsWorkstationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkstationsWorkstationConfig_hyperdisk(context),
			},
			{
				ResourceName:            "google_workstations_workstation_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag"},
			},
			{
				ResourceName:            "google_workstations_workstation_config.source_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag"},
			},
		},
	})
}

func testAccWorkstationsWorkstationConfig_hyperdisk(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-workstation-cluster%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-workstation-cluster%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-east5"
  network       = google_compute_network.default.name
}

resource "google_compute_disk" "my_source_disk" {
  name = "tf-test-workstation-source-disk%{random_suffix}"
  size = 10
  type = "pd-ssd"
  zone = "us-east5-a"
}

resource "google_compute_snapshot" "my_source_snapshot" {
  name        = "tf-test-workstation-source-snapshot%{random_suffix}"
  source_disk = google_compute_disk.my_source_disk.name
  zone        = "us-east5-a"
}

resource "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "tf-test-workstation-cluster%{random_suffix}"
  network                = google_compute_network.default.id
  subnetwork             = google_compute_subnetwork.default.id
  location               = "us-east5"
}

resource "google_workstations_workstation_config" "default" {
  workstation_config_id  = "tf-test-workstation-config%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east5"

  host {
    gce_instance {
      machine_type = "c3-standard-22"
    }
  }

  persistent_directories {
    mount_path = "/home"
    gce_hd {
      size_gb         = 200
      reclaim_policy  = "DELETE"
      archive_timeout = "3600s"
    }
  }
}

resource "google_workstations_workstation_config" "source_snapshot" {
  workstation_config_id  = "tf-test-ws-cfg-hd-snap%{random_suffix}"
  workstation_cluster_id = google_workstations_workstation_cluster.default.workstation_cluster_id
  location               = "us-east5"

  host {
    gce_instance {
      machine_type                = "c3-standard-22"
      boot_disk_size_gb           = 35
      disable_public_ip_addresses = true
    }
  }

  persistent_directories {
    mount_path = "/home"
    gce_hd {
      source_snapshot = google_compute_snapshot.my_source_snapshot.id
      reclaim_policy  = "DELETE"
      archive_timeout = "3600s"
    }
  }
}
`, context)
}
