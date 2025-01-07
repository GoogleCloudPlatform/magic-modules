package redis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Validate that cluster endpoints are updated for the cluster
func TestAccRedisCluster_updateClusterEndpoints(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with no user created connections
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// create cluster with one user created connection
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 1}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 2 endpoints
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 2}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 0 endpoints
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
		},
	})
}

type ClusterParams struct {
	name                      string
	replicaCount              int
	shardCount                int
	deletionProtectionEnabled bool
	nodeType                  string
	redisConfigs              map[string]string
	zoneDistributionMode      string
	zone                      string
	maintenanceDay            string
	maintenanceHours          int
	maintenanceMinutes        int
	maintenanceSeconds        int
	maintenanceNanos          int
	persistenceBlock          string
	shouldCreateSecondary     bool
	secondaryClusterName      string
	ccrRole                   string
	userEndpointCount         int
}

func createRedisClusterEndpoints(params *ClusterParams) string {
	if params.userEndpointCount == 2 {
		return createRedisClusterEndpointsWithTwoUserCreatedConnections(params)
	} else if params.userEndpointCount == 1 {
		return createRedisClusterEndpointsWithOneUserCreatedConnections(params)
	}
	return ``
}

func createRedisClusterEndpointsWithOneUserCreatedConnections(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_redis_cluster_user_created_connections" "default" {
		name = "%s"
		region = "us-central1"
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					address = google_compute_address.ip1_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network = google_compute_network.network1.id
					project_id = data.google_project.project.project_id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					address = google_compute_address.ip2_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network = google_compute_network.network1.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		}
		%s
		`,
		params.name,
		createRedisClusterUserCreatedConnection1(params),
	)

}

func createRedisClusterEndpointsWithTwoUserCreatedConnections(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_redis_cluster_user_created_connections" "default" {
		name = "%s"
		region = "us-central1"
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					address = google_compute_address.ip1_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network = google_compute_network.network1.id
					project_id = data.google_project.project.project_id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					address = google_compute_address.ip2_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network = google_compute_network.network1.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network2.psc_connection_id
					address = google_compute_address.ip1_network2.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network2.id
					network = google_compute_network.network2.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network2.psc_connection_id
					address = google_compute_address.ip2_network2.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network2.id
					network = google_compute_network.network2.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		}
		%s
		%s
		`,
		params.name,
		createRedisClusterUserCreatedConnection1(params),
		createRedisClusterUserCreatedConnection2(params),
	)
}

func createRedisClusterUserCreatedConnection1(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network1" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip1_network1.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network1.id
		target                 = google_redis_cluster.test.psc_service_attachments[0].service_attachment
		}

		resource "google_compute_forwarding_rule" "forwarding_rule2_network1" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip2_network1.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network1.id
		target                 = google_redis_cluster.test.psc_service_attachments[1].service_attachment
		}

		resource "google_compute_address" "ip1_network1" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network1.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network1" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network1.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network1" {
		name          = "%s"
		ip_cidr_range = "10.0.0.248/29"
		region        = "us-central1"
		network       = google_compute_network.network1.id
		}

		resource "google_compute_network" "network1" {
		name                    = "%s"
		auto_create_subnetworks = false
		}
		
		data "google_project" "project" {
		}
		`,
		params.name+"-11", // fwd-rule1-net1
		params.name+"-12", // fwd-rule2-net1
		params.name+"-11", // ip1-net1
		params.name+"-12", // ip2-net1
		params.name+"-1",  // subnet-net1
		params.name+"-1",  // net1
	)
}

func createRedisClusterUserCreatedConnection2(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network2" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip1_network2.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network2.id
		target                 = google_redis_cluster.test.psc_service_attachments[0].service_attachment
		}

		resource "google_compute_forwarding_rule" "forwarding_rule2_network2" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip2_network2.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network2.id
		target                 = google_redis_cluster.test.psc_service_attachments[1].service_attachment
		}

		resource "google_compute_address" "ip1_network2" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network2.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network2" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network2.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network2" {
		name          = "%s"
		ip_cidr_range = "10.0.0.248/29"
		region        = "us-central1"
		network       = google_compute_network.network2.id
		}

		resource "google_compute_network" "network2" {
		name                    = "%s"
		auto_create_subnetworks = false
		}

		data "google_project" "project" {
		}
		`,
		params.name+"-21", // fwd-rule1-net2
		params.name+"-22", // fwd-rule2-net2
		params.name+"-21", // ip1-net2
		params.name+"-22", // ip2-net2
		params.name+"-2",  // subnet-net2
		params.name+"-2",  // net2
	)

}

func createOrUpdateRedisCluster(params *ClusterParams) string {
	clusterResourceBlock := createRedisClusterResourceConfig(params /*isSecondaryCluster*/, false)
	secClusterResourceBlock := ``
	if params.shouldCreateSecondary {
		secClusterResourceBlock = createRedisClusterResourceConfig(params /*isSecondaryCluster*/, true)
	}
	endpointBlock := ``
	if params.userEndpointCount > 0 {
		endpointBlock = createRedisClusterUserCreatedConnection(params)
	}

	return fmt.Sprintf(`
		%s
		%s
		resource "google_network_connectivity_service_connection_policy" "default" {
			name = "%s"
			location = "us-central1"
			service_class = "gcp-memorystore-redis"
			description   = "my basic service connection policy"
			network = google_compute_network.producer_net.id
			psc_config {
			subnetworks = [google_compute_subnetwork.producer_subnet.id]
			}
		}

		resource "google_compute_subnetwork" "producer_subnet" {
			name          = "%s"
			ip_cidr_range = "10.0.0.16/28"
			region        = "us-central1"
			network       = google_compute_network.producer_net.id
		}

		resource "google_compute_network" "producer_net" {
			name                    = "%s"
			auto_create_subnetworks = false
		}
		%s
		`,
		clusterResourceBlock,
		secClusterResourceBlock,
		params.name,
		params.name,
		params.name,
		endpointBlock)
}

func createRedisClusterResourceConfig(params *ClusterParams, isSecondaryCluster bool) string {
	tfClusterResourceName := "test"
	clusterName := params.name
	dependsOnBlock := "google_network_connectivity_service_connection_policy.default"

	var redsConfigsStrBuilder strings.Builder
	for key, value := range params.redisConfigs {
		redsConfigsStrBuilder.WriteString(fmt.Sprintf("%s =  \"%s\"\n", key, value))
	}

	zoneDistributionConfigBlock := ``
	if params.zoneDistributionMode != "" {
		zoneDistributionConfigBlock = fmt.Sprintf(`
		zone_distribution_config {
			mode = "%s"
			zone = "%s"
		}
		`, params.zoneDistributionMode, params.zone)
	}

	maintenancePolicyBlock := ``
	if params.maintenanceDay != "" {
		maintenancePolicyBlock = fmt.Sprintf(`
		maintenance_policy {
			weekly_maintenance_window {
				day = "%s"
				start_time {
					hours = %d
					minutes = %d
					seconds = %d
					nanos = %d
				}
			}
		}
		`, params.maintenanceDay, params.maintenanceHours, params.maintenanceMinutes, params.maintenanceSeconds, params.maintenanceNanos)
	}

	crossClusterReplicationConfigBlock := ``
	if isSecondaryCluster {
		tfClusterResourceName = "test_secondary"
		clusterName = params.secondaryClusterName
		dependsOnBlock = dependsOnBlock + ", google_redis_cluster.test"

		// Construct cross_cluster_replication_config block
		pcBlock := ``
		scsBlock := ``
		if params.ccrRole == "SECONDARY" {
			pcBlock = fmt.Sprintf(`
			primary_cluster {
				cluster = google_redis_cluster.test.id
			}
			`)
		} else if params.ccrRole == "PRIMARY" {
			scsBlock = fmt.Sprintf(`
			secondary_clusters {
				cluster = google_redis_cluster.test.id
			}
			`)
		}
		crossClusterReplicationConfigBlock = fmt.Sprintf(`
		cross_cluster_replication_config {
			cluster_role = "%s"
			%s
			%s
		}
		`, params.ccrRole, pcBlock, scsBlock)
	}

	return fmt.Sprintf(`
		resource "google_redis_cluster" "%s" {
		name           = "%s"
		replica_count = %d
		shard_count = %d
		node_type = "%s"
		deletion_protection_enabled = %v
		region         = "us-central1"
		psc_configs {
				network = google_compute_network.producer_net.id
		}
		redis_configs = {
			%s
		}
		%s
		%s
		%s
		%s
		depends_on = [
				%s
			]
		}
		`,
		tfClusterResourceName,
		clusterName,
		params.replicaCount,
		params.shardCount,
		params.nodeType,
		params.deletionProtectionEnabled,
		redsConfigsStrBuilder.String(),
		zoneDistributionConfigBlock,
		maintenancePolicyBlock,
		params.persistenceBlock,
		crossClusterReplicationConfigBlock,
		dependsOnBlock)
}
