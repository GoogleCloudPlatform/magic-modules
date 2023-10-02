package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbUser_updateRoles_BuiltIn(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"alloydb_cluster_name":         "tf-test-user-cluster-" + acctest.RandString(t, 10),
		"alloydb_cluster_pass":         "cluster_secret3",
		"alloydb_cluster_resource_id":  "test-cluster3",
		"alloydb_instance_name":        "tf-test-instance3",
		"alloydb_instance_resource_id": "test-instance3",
		"alloydb_user_name":            "user_3",
		"alloydb_user_pass":            "user_3_pass",
		"random_suffix":                acctest.RandString(t, 10),
		"network_name":                 acctest.BootstrapSharedTestNetwork(t, "alloydb-user-updaterole-builtin"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAlloydbUser_updateRoles_BuiltIn(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAlloydbUser_updateRoles_BuiltIn(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "%{alloydb_instance_resource_id}" {
  cluster       = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "%{alloydb_cluster_resource_id}" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id

  initial_user {
    password = "%{alloydb_cluster_pass}"
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          = "%{alloydb_cluster_name}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"

  password = "%{alloydb_user_pass}"
  database_roles = []
  depends_on = [google_alloydb_instance.%{alloydb_instance_resource_id}]
}`, context)
}

func TestAccAlloydbUser_updatePassword_BuiltIn(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"alloydb_cluster_name":         "tf-test-user-cluster-" + acctest.RandString(t, 10),
		"alloydb_cluster_pass":         "cluster_secret3",
		"alloydb_cluster_resource_id":  "test-cluster3",
		"alloydb_instance_name":        "tf-test-instance3",
		"alloydb_instance_resource_id": "test-instance3",
		"alloydb_user_name":            "user_3",
		"alloydb_user_pass":            "user_3_pass",
		"random_suffix":                acctest.RandString(t, 10),
		"network_name":                 acctest.BootstrapSharedTestNetwork(t, "alloydb-user-updatepass-builtin"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAlloydbUser_updatePass_BuiltIn(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAlloydbUser_updatePass_BuiltIn(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "%{alloydb_instance_resource_id}" {
  cluster       = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "%{alloydb_cluster_resource_id}" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id

  initial_user {
    password = "%{alloydb_cluster_pass}"
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          = "%{alloydb_cluster_name}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"

  password = "%{alloydb_user_pass}-foo"
  database_roles = ["alloydbsuperuser"]
  depends_on = [google_alloydb_instance.%{alloydb_instance_resource_id}]
}`, context)
}

func TestAccAlloydbUser_updateRoles_IAM(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"alloydb_cluster_name":         "tf-test-user-cluster-" + acctest.RandString(t, 10),
		"alloydb_cluster_pass":         "cluster_secret3",
		"alloydb_cluster_resource_id":  "test-cluster3",
		"alloydb_instance_name":        "tf-test-instance3",
		"alloydb_instance_resource_id": "test-instance3",
		"alloydb_user_name":            "user_3@foo.com",
		"alloydb_user_pass":            "user_3_pass",
		"random_suffix":                acctest.RandString(t, 10),
		"network_name":                 acctest.BootstrapSharedTestNetwork(t, "alloydb-user-updaterole-iam"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserIamExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccAlloydbUser_updateRoles_Iam(context),
			},
			{
				ResourceName:            "google_alloydb_user.user2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAlloydbUser_updateRoles_Iam(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "%{alloydb_instance_resource_id}" {
  cluster       = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"
  depends_on = [google_service_networking_connection.vpc_connection]
}
resource "google_alloydb_cluster" "%{alloydb_cluster_resource_id}" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
  initial_user {
    password = "%{alloydb_cluster_pass}"
  }
}
data "google_project" "project" {}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_compute_global_address" "private_ip_alloc" {
  name          = "%{alloydb_cluster_name}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}
resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
resource "google_alloydb_user" "user2" {
  cluster = google_alloydb_cluster.%{alloydb_cluster_resource_id}.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_IAM_USER"
  database_roles = ["alloydbiamuser", "alloydbsuperuser"]
  depends_on = [google_alloydb_instance.%{alloydb_instance_resource_id}]
}`, context)
}
