package alloydb_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// The cluster creation should succeed with minimal number of arguments
func TestAccAlloydbCluster_secondaryClusterMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`

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

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// This test passes if secondary cluster can be created with existing primary cluster and primary instance
func TestAccAlloydbCluster_secondaryClusterWithPrimaryClusterAndInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterWithPrimaryClusterAndInstanceCreate(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterWithPrimaryClusterAndInstanceRemoveSecondary(context),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterWithPrimaryClusterAndInstanceCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`

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

	machine_config {
	  cpu_count = 2
	}
	depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

func testAccAlloydbCluster_secondaryClusterWithPrimaryClusterAndInstanceRemoveSecondary(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_instance" "default" {
	cluster       = google_alloydb_cluster.default.name
	instance_id   = "tf-test-alloydb-instance%{random_suffix}"
	instance_type = "PRIMARY"

	machine_config {
	  cpu_count = 2
	}
	depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

// Validation test to ensure proper error is raised if create secondary cluster is called without any secondary_config field
func TestAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Can not create secondary cluster without secondary_config field."),
			},
		},
	})
}
func testAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  depends_on = [google_alloydb_cluster.default]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// Validation test to ensure proper error is raised if secondary_config field is provided but no cluster_type is specified.
func TestAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Add {cluster_type: \"SECONDARY\"} if attempting to create a secondary cluster, otherwise remove the secondary_config."),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// Validation test to ensure proper error is raised if secondary_config field is provided but cluster_type is primary
func TestAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Add {cluster_type: \"SECONDARY\"} if attempting to create a secondary cluster, otherwise remove the secondary_config."),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// This test passes if secondary cluster can be updated
func TestAccAlloydbCluster_secondaryClusterUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterUpdate(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`

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

  labels = {
    foo = "bar"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// This test passes if secondary cluster can be deleted
func TestAccAlloydbCluster_secondaryClusterDelete(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterDelete(context),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterDelete(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

// This test passes if secondary cluster can be promoted and the original primary clusters can be deleted after that
func TestAccAlloydbCluster_secondaryClusterPromote(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceCreateWithMandatoryFields(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteDeleteOriginalPrimary(context),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromote(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  depends_on = [google_service_networking_connection.vpc_connection]

  lifecycle {
    ignore_changes = [instance_type]
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
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

func testAccAlloydbCluster_secondaryClusterPromoteDeleteOriginalPrimary(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  depends_on = [google_service_networking_connection.vpc_connection]

  lifecycle {
    ignore_changes = [instance_type]
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }
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
