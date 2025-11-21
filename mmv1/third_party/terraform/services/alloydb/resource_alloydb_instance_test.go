package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbInstance_update(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_alloydbInstanceBasic(context),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_update(context),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_alloydbInstanceBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 4
    machine_type = "n2-highmem-4"
  }

  labels = {
	test = "tf-test-alloydb-instance%{random_suffix}"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance with minimal number of fields
func TestAccAlloydbInstance_createInstanceWithMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithMandatoryFields(context),
			},
		},
	})
}

// This test passes if we are able to create a primary instance STOP it and then START it back again
func TestAccAlloydbInstance_stopstart(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1")

	context := map[string]interface{}{
		"random_suffix": suffix,
		"network_name":  networkName,
	}

	contextStop := map[string]interface{}{
		"random_suffix":     suffix,
		"network_name":      networkName,
		"activation_policy": "NEVER",
	}

	contextStart := map[string]interface{}{
		"random_suffix":     suffix,
		"network_name":      networkName,
		"activation_policy": "ALWAYS",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithMandatoryFields(context),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_updateActivationPolicy(contextStop),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "activation_policy", "NEVER"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "state", "STOPPED"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
			{
				Config: testAccAlloydbInstance_updateActivationPolicy(contextStart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "activation_policy", "ALWAYS"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "state", "READY"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_updateActivationPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  activation_policy = "%{activation_policy}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
		password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_createInstanceWithMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance with maximum number of fields
/* func TestAccAlloydbInstance_createInstanceWithMaximumFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithMaximumFields(context),
			},
		},
	})
}

func testAccAlloydbInstance_createInstanceWithMaximumFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  labels = {
    test_label = "test-alloydb-label"
  }
  annotations = {
    test_annotation = "test-alloydb-annotation"
  }
  gce_zone = "us-east1-a"
  database_flags = {
	  "alloydb.enable_auto_explain" = "true"
  }
  availability_type = "REGIONAL"
  machine_config {
	  cpu_count = 4
  }
  query_insights_config {
    query_string_length = 300
    record_application_tags = "false"
    record_client_address = "true"
    query_plans_per_minute = 10
  }
  depends_on = [google_service_networking_connection.vpc_connection]
  lifecycle {
    ignore_changes = [
      gce_zone,
      annotations
    ]
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}*/

// This test passes if we are able to create a primary instance with an associated read-pool instance
func TestAccAlloydbInstance_createPrimaryAndReadPoolInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createPrimaryAndReadPoolInstance(context),
			},
		},
	})
}

func testAccAlloydbInstance_createPrimaryAndReadPoolInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_instance" "read_pool" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}-read"
  instance_type = "READ_POOL"
  read_pool_config {
    node_count = 4
  }
  depends_on = [google_alloydb_instance.primary]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to update a database flag in primary instance
/*func TestAccAlloydbInstance_updateDatabaseFlagInPrimaryInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_autoExplainEnabledInPrimaryInstance(context),
			},
			{
				ResourceName:      "google_alloydb_instance.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbInstance_autoExplainDisabledInPrimaryInstance(context),
			},
		},
	})
}

func testAccAlloydbInstance_autoExplainEnabledInPrimaryInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
	  "alloydb.enable_auto_explain" = "true"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}*/

func testAccAlloydbInstance_autoExplainDisabledInPrimaryInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
	  "alloydb.enable_auto_explain" = "false"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance by specifying network_config.network and network_config.allocated_ip_range
func TestAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(t *testing.T) {
	t.Parallel()

	testId := "alloydb-1"
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  networkName,
		"address_name":  addressName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(context),
			},
		},
	})
}

func testAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network    = data.google_compute_network.default.id
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name =  "%{address_name}"
}
`, context)
}

// This test passes if an instance is able to be created specifying require
// connectors and the ssl mode; if the instance is able to update require
// connectors, and the ssl mode in the client connection config; if the ssl
// mode specified is removed it doesn't not change the ssl mode; and if the
// require connectors is remove it doesn't change require connectors either.
func TestAccAlloydbInstance_clientConnectionConfig(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1")

	context := map[string]interface{}{
		"random_suffix":      suffix,
		"network_name":       networkName,
		"require_connectors": true,
		"ssl_mode":           "ENCRYPTED_ONLY",
	}
	context2 := map[string]interface{}{
		"random_suffix":      suffix,
		"network_name":       networkName,
		"require_connectors": false,
		"ssl_mode":           "ALLOW_UNENCRYPTED_AND_ENCRYPTED",
	}
	context3 := map[string]interface{}{
		"random_suffix":      suffix,
		"network_name":       networkName,
		"require_connectors": false,
	}
	context4 := map[string]interface{}{
		"random_suffix":      suffix,
		"network_name":       networkName,
		"require_connectors": false,
		"ssl_mode":           "ENCRYPTED_ONLY",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_clientConnectionConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.require_connectors", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.ssl_config.0.ssl_mode", "ENCRYPTED_ONLY"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_clientConnectionConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.require_connectors", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.ssl_config.0.ssl_mode", "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_noSSLModeConfig(context3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.require_connectors", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.ssl_config.0.ssl_mode", "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_clientConnectionConfig(context4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.require_connectors", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.ssl_config.0.ssl_mode", "ENCRYPTED_ONLY"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_noSSLModeConfig(context3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.require_connectors", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "client_connection_config.0.ssl_config.0.ssl_mode", "ENCRYPTED_ONLY"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
		},
	})
}

func testAccAlloydbInstance_noSSLModeConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  client_connection_config {
    require_connectors = %{require_connectors}
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_clientConnectionConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  client_connection_config {
    require_connectors = %{require_connectors}
    ssl_config {
      ssl_mode = "%{ssl_mode}"
    }
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
`, context)
}

// This test passes if an instance can be created with public IP enabled,
// and update the authorized external networks.
func TestAccAlloydbInstance_networkConfig(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1")

	context1 := map[string]interface{}{
		"random_suffix":                suffix,
		"network_name":                 networkName,
		"enable_public_ip":             true,
		"enable_outbound_public_ip":    true,
		"authorized_external_networks": "",
	}

	context2 := map[string]interface{}{
		"random_suffix":             suffix,
		"network_name":              networkName,
		"enable_public_ip":          true,
		"enable_outbound_public_ip": false,
		"authorized_external_networks": `
		authorized_external_networks {
			cidr_range = "8.8.8.8/30"
		}
		authorized_external_networks {
			cidr_range = "8.8.4.4/30"
		}
		`,
	}

	context3 := map[string]interface{}{
		"random_suffix":             suffix,
		"network_name":              networkName,
		"enable_public_ip":          true,
		"enable_outbound_public_ip": true,
		"cidr_range":                "8.8.8.8/30",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_networkConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_public_ip", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_outbound_public_ip", "true"),
					resource.TestCheckResourceAttrSet("google_alloydb_instance.default", "outbound_public_ip_addresses.0"), // Ensure it's set
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_networkConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_public_ip", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.authorized_external_networks.0.cidr_range", "8.8.8.8/30"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.authorized_external_networks.1.cidr_range", "8.8.4.4/30"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.authorized_external_networks.#", "2"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_outbound_public_ip", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "outbound_public_ip_addresses.#", "0"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_networkConfigWithAnAuthNetwork(context3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_public_ip", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.authorized_external_networks.0.cidr_range", "8.8.8.8/30"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.authorized_external_networks.#", "1"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "network_config.0.enable_outbound_public_ip", "true"),
					resource.TestCheckResourceAttrSet("google_alloydb_instance.default", "outbound_public_ip_addresses.0"),
				),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
		},
	})
}

func testAccAlloydbInstance_networkConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
    "password.enforce_complexity" = "on"
  }

  network_config {
    enable_public_ip = %{enable_public_ip}
    enable_outbound_public_ip = %{enable_outbound_public_ip}
    %{authorized_external_networks}
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_networkConfigWithAnAuthNetwork(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
    "password.enforce_complexity" = "on"
  }

  network_config {
    enable_public_ip = %{enable_public_ip}
    enable_outbound_public_ip = %{enable_outbound_public_ip}
    authorized_external_networks {
      cidr_range = "%{cidr_range}"
    }
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
`, context)
}

func TestAccAlloydbInstance_updatePscInstanceConfig(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	context1 := map[string]interface{}{
		"random_suffix": random_suffix,
		"psc_enabled":   true,
	}

	context2 := map[string]interface{}{
		"random_suffix": random_suffix,
		"psc_enabled":   false,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_pscInstanceConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "psc_config.0.psc_enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "psc_instance_config.0.allowed_consumer_projects.#", "1"),
				),
			},
			{
				Config: testAccAlloydbInstance_pscInstanceConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "psc_config.0.psc_enabled", "false"),
				),
			},
			{
				Config: testAccAlloydbInstance_updatePscInstanceConfigAllowlist(context1), // context1 has psc_enabled: true
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "psc_config.0.psc_enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "psc_instance_config.0.allowed_consumer_projects.#", "2"),
				),
			},
		},
	})
}

func testAccAlloydbInstance_pscInstanceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
  psc_instance_config {
	allowed_consumer_projects = ["${data.google_project.project.number}"]
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
    psc_enabled = %{psc_enabled}
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}
data "google_project" "project" {}
`, context)
}

func testAccAlloydbInstance_updatePscInstanceConfigAllowlist(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
  psc_instance_config {
	allowed_consumer_projects = ["${data.google_project.project.number}", "1044355742748"]
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
    psc_enabled = %{psc_enabled}
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}
data "google_project" "project" {}
`, context)
}

func TestAccAlloydbInstance_createInstanceWithPscInterfaceConfigs(t *testing.T) {
	t.Parallel()

	networkName := acctest.BootstrapSharedTestNetwork(t, "tf-test-alloydb-network")
	subnetworkName := acctest.BootstrapSubnet(t, "tf-test-alloydb-subnetwork", networkName)

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix":         random_suffix,
		"networkAttachmentName": acctest.BootstrapNetworkAttachment(t, "tf-test-alloydb-create-na", subnetworkName),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_pscInterfaceConfigs(context),
			},
		},
	})
}

func testAccAlloydbInstance_pscInterfaceConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
    machine_type = "n2-highmem-2"
  }
  psc_instance_config {
	allowed_consumer_projects = ["${data.google_project.project.number}"]
	psc_interface_configs {
		network_attachment_resource = "projects/${data.google_project.project.number}/regions/${google_alloydb_cluster.default.location}/networkAttachments/%{networkAttachmentName}"
	}
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
	psc_enabled = true
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}
data "google_project" "project" {}
`, context)
}

func TestAccAlloydbInstance_updateInstanceWithPscInterfaceConfigs(t *testing.T) {
	t.Parallel()

	networkName := acctest.BootstrapSharedTestNetwork(t, "tf-test-alloydb-network")
	subnetworkName := acctest.BootstrapSubnet(t, "tf-test-alloydb-subnetwork", networkName)

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix":         random_suffix,
		"networkAttachmentName": acctest.BootstrapNetworkAttachment(t, "tf-test-alloydb-update-na", subnetworkName),
		"psc_enabled":           true, // Needed for the first step
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_pscInstanceConfig(context),
			},
			{
				Config: testAccAlloydbInstance_pscInterfaceConfigs(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "psc_instance_config.0.psc_interface_configs.#", "1"),
					resource.TestCheckResourceAttrSet("google_alloydb_instance.default", "psc_instance_config.0.psc_interface_configs.0.network_attachment"),
				),
			},
		},
	})
}

func TestAccAlloydbInstance_updatePscAutoConnections(t *testing.T) {
	t.Parallel()

	networkName := acctest.BootstrapSharedTestNetwork(t, "tf-test-alloydb-network-psc")
	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"network_name":  networkName,
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_pscAutoConnections(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "psc_instance_config.0.psc_auto_connections.#", "1"),
				),
			},
			{
				Config: testAccAlloydbInstance_updatePscAutoConnections(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "psc_instance_config.#", "0"),
				),
			},
		},
	})
}

func testAccAlloydbInstance_pscAutoConnections(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
  psc_instance_config {
	psc_auto_connections {
		consumer_project = "${data.google_project.project.project_id}"
		consumer_network = "projects/${data.google_project.project.project_id}/global/networks/%{network_name}"
	}
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
	psc_enabled = true
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}
data "google_project" "project" {}
`, context)
}

func testAccAlloydbInstance_updatePscAutoConnections(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
	psc_enabled = true
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}
data "google_project" "project" {}
`, context)
}

func TestAccAlloydbInstance_createPrimaryAndReadPoolInstanceWithAllocatedIpRangeOverride(t *testing.T) {
	t.Parallel()

	testId := "alloydb-1"
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"address_name":  acctest.BootstrapSharedTestGlobalAddress(t, testId),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, testId),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createPrimaryAndReadPoolInstanceWithAllocatedIpRangeOverride(context),
			},
		},
	})
}

func testAccAlloydbInstance_createPrimaryAndReadPoolInstanceWithAllocatedIpRangeOverride(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_instance" "read_pool" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}-read"
  instance_type = "READ_POOL"
  read_pool_config {
    node_count = 4
  }
  network_config {
	allocated_ip_range_override = data.google_compute_global_address.private_ip_alloc.name
  }
  depends_on = [google_alloydb_instance.primary]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name =  "%{address_name}"
}
`, context)
}

func TestAccAlloydbInstance_Update_ObservabilityConfig(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1")

	context1 := map[string]interface{}{
		"random_suffix":                 random_suffix,
		"network_name":                  networkName,
		"enabled":                       true,
		"preserve_comments":             true,
		"track_wait_events":             true,
		"track_wait_event_types":        true,
		"max_query_string_length":       1024,
		"record_application_tags":       true,
		"query_plans_per_minute":        10,
		"track_client_address":          true,
		"track_active_queries":          true,
		"assistive_experiences_enabled": true,
	}

	context2 := map[string]interface{}{
		"random_suffix":                 random_suffix,
		"network_name":                  networkName,
		"enabled":                       false,
		"preserve_comments":             false,
		"track_wait_events":             false,
		"track_wait_event_types":        false,
		"max_query_string_length":       1023,
		"record_application_tags":       false,
		"query_plans_per_minute":        8,
		"track_client_address":          false,
		"track_active_queries":          false,
		"assistive_experiences_enabled": false,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_ObservabilityConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.preserve_comments", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_events", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_event_types", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.max_query_string_length", "1024"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.record_application_tags", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.query_plans_per_minute", "10"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_client_address", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_active_queries", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.assistive_experiences_enabled", "true"),
				),
			},
			{
				Config: testAccAlloydbInstance_ObservabilityConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.enabled", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.preserve_comments", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_events", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_event_types", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.max_query_string_length", "1023"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.record_application_tags", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.query_plans_per_minute", "8"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_client_address", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_active_queries", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.assistive_experiences_enabled", "false"),
				),
			},
			{
				Config: testAccAlloydbInstance_ObservabilityConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.preserve_comments", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_events", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_wait_event_types", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.max_query_string_length", "1024"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.record_application_tags", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.query_plans_per_minute", "10"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_client_address", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.track_active_queries", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "observability_config.0.assistive_experiences_enabled", "true"),
				),
			},
		},
	})
}

func testAccAlloydbInstance_ObservabilityConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
  observability_config {
    enabled                        = %{enabled}
    preserve_comments              = %{preserve_comments}
    track_wait_events              = %{track_wait_events}
    track_wait_event_types         = %{track_wait_event_types}
    max_query_string_length        = %{max_query_string_length}
    record_application_tags        = %{record_application_tags}
    query_plans_per_minute         = %{query_plans_per_minute}
    track_client_address           = %{track_client_address}
    track_active_queries           = %{track_active_queries}
    assistive_experiences_enabled  = %{assistive_experiences_enabled}
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
  deletion_protection = false
}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
data "google_project" "project" {}
`, context)
}

func TestAccAlloydbInstance_Update_QueryInsightsConfig(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1")

	context1 := map[string]interface{}{
		"random_suffix":           random_suffix,
		"network_name":            networkName,
		"query_string_length":     256,
		"record_application_tags": true,
		"record_client_address":   true,
		"query_plans_per_minute":  5,
	}

	context2 := map[string]interface{}{
		"random_suffix":           random_suffix,
		"network_name":            networkName,
		"query_string_length":     257,
		"record_application_tags": false,
		"record_client_address":   false,
		"query_plans_per_minute":  10,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_QueryInsightsConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_string_length", "256"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_application_tags", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_client_address", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_plans_per_minute", "5"),
				),
			},
			{
				Config: testAccAlloydbInstance_QueryInsightsConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_string_length", "257"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_application_tags", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_client_address", "false"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_plans_per_minute", "10"),
				),
			},
			{
				Config: testAccAlloydbInstance_QueryInsightsConfig(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_string_length", "256"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_application_tags", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.record_client_address", "true"),
					resource.TestCheckResourceAttr("google_alloydb_instance.default", "query_insights_config.0.query_plans_per_minute", "5"),
				),
			},
		},
	})
}

func testAccAlloydbInstance_QueryInsightsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  machine_config {
    cpu_count = 2
  }
  query_insights_config {
    query_string_length   = %{query_string_length}
    record_application_tags = %{record_application_tags}
    record_client_address = %{record_client_address}
    query_plans_per_minute  = %{query_plans_per_minute}
  }
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
  deletion_protection = false
}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
data "google_project" "project" {}
`, context)
}
