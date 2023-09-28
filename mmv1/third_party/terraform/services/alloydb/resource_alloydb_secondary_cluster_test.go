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

  depends_on = [google_alloydb_cluster.default]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  cluster_type = "PRIMARY"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

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
  cluster_type = "PRIMARY"
}

data "google_project" "project" {}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

`, context)
}

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

  continuous_backup_config {
    enabled = false
  }

  cluster_type = "PRIMARY"

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.default.name
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
