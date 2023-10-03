package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// The instance creation should succeed with minimal number of arguments
func TestAccAlloydbInstance_secondaryInstanceCreateWithMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceCreateWithMandatoryFields(context),
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceRemoveSecondaryCluster(context),
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceCreateWithMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  depends_on = [google_service_networking_connection.vpc_connection, google_alloydb_instance.default]

  lifecycle {
    ignore_changes = [instance_type]
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
  }
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-instance%{random_suffix}"
}

resource "google_compute_global_address" "private_ip_alloc" {
	name          =  "tf-test-alloydb-cluster%{random_suffix}"
	address_type  = "INTERNAL"
	purpose       = "VPC_PEERING"
	prefix_length = 16
	network       = google_compute_network.default.id
  }

resource "google_service_networking_connection" "vpc_connection" {
	network                 = google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

`, context)
}

func testAccAlloydbInstance_secondaryInstanceRemoveSecondaryCluster(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-instance%{random_suffix}"
}

resource "google_compute_global_address" "private_ip_alloc" {
	name          =  "tf-test-alloydb-cluster%{random_suffix}"
	address_type  = "INTERNAL"
	purpose       = "VPC_PEERING"
	prefix_length = 16
	network       = google_compute_network.default.id
  }

  resource "google_service_networking_connection" "vpc_connection" {
	network                 = google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
  }

`, context)
}

// This test passes if secondary instance can be updated
func TestAccAlloydbInstance_secondaryInstanceUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydb-secondary-instance-update"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceCreateWithMandatoryFields(context),
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceUpdate(context),
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceRemoveSecondaryCluster(context),
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  depends_on = [google_service_networking_connection.vpc_connection, google_alloydb_instance.default]

  // Default machine_config.cpu_count = 2
  machine_config {
    cpu_count = 4
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
  }
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-instance%{random_suffix}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = google_compute_network.default.id
  }

  resource "google_service_networking_connection" "vpc_connection" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
  }
  
`, context)
}

// This test passes if we are able to create a secondary instance with an associated read-pool instance
func TestAccAlloydbInstance_secondaryInstanceWithReadPool(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceWithReadPool(context),
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceRemoveSecondaryCluster(context),
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceWithReadPool(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  depends_on = [google_service_networking_connection.vpc_connection, google_alloydb_instance.default]

  lifecycle {
    ignore_changes = [instance_type]
  }
}

resource "google_alloydb_instance" "read_pool" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-read-pool%{random_suffix}"
  instance_type = "READ_POOL"
  read_pool_config {
    node_count = 4
  }
  depends_on = [google_service_networking_connection.vpc_connection, google_alloydb_instance.secondary]

}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
  }
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-instance%{random_suffix}"
}

resource "google_compute_global_address" "private_ip_alloc" {
	name          =  "tf-test-alloydb-cluster%{random_suffix}"
	address_type  = "INTERNAL"
	purpose       = "VPC_PEERING"
	prefix_length = 16
	network       = google_compute_network.default.id
  }

  resource "google_service_networking_connection" "vpc_connection" {
	network                 = google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
  }

`, context)
}
