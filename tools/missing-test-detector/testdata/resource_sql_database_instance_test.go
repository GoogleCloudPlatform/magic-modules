package google

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// Fields that should be ignored in import tests because they aren't returned
// from GCP (and thus can't be imported)
var ignoredReplicaConfigurationFields = []string{
	"replica_configuration.0.ca_certificate",
	"replica_configuration.0.client_certificate",
	"replica_configuration.0.client_key",
	"replica_configuration.0.connect_retry_interval",
	"replica_configuration.0.dump_file_path",
	"replica_configuration.0.master_heartbeat_period",
	"replica_configuration.0.password",
	"replica_configuration.0.ssl_cipher",
	"replica_configuration.0.username",
	"replica_configuration.0.verify_server_certificate",
	"deletion_protection",
}

func init() {
	resource.AddTestSweepers("gcp_sql_db_instance", &resource.Sweeper{
		Name: "gcp_sql_db_instance",
		F:    testSweepDatabases,
	})
}

func TestMaintenanceVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New       string
		ShouldSuppress bool
	}{
		"older configuration maintenance version than current version should suppress diff": {
			Old:            "MYSQL_8_0_26.R20220508.01_09",
			New:            "MYSQL_5_7_37.R20210508.01_03",
			ShouldSuppress: true,
		},
		"older configuration maintenance version than current version should suppress diff with lexicographically smaller database version": {
			Old:            "MYSQL_5_8_10.R20220508.01_09",
			New:            "MYSQL_5_8_7.R20210508.01_03",
			ShouldSuppress: true,
		},
		"newer configuration maintenance version than current version should not suppress diff": {
			Old:            "MYSQL_5_7_37.R20210508.01_03",
			New:            "MYSQL_8_0_26.R20220508.01_09",
			ShouldSuppress: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if maintenanceVersionDiffSuppress("version", tc.Old, tc.New, nil) != tc.ShouldSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ShouldSuppress)
			}
		})
	}
}

func testSweepDatabases(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	found, err := config.NewSqlAdminClient(config.userAgent).Instances.List(config.Project).Do()
	if err != nil {
		log.Printf("error listing databases: %s", err)
		return nil
	}

	if len(found.Items) == 0 {
		log.Printf("No databases found")
		return nil
	}

	running := map[string]struct{}{}

	for _, d := range found.Items {
		var testDbInstance bool
		for _, testName := range []string{"tf-lw-", "sqldatabasetest"} {
			// only destroy instances we know to fit our test naming pattern
			if strings.HasPrefix(d.Name, testName) {
				testDbInstance = true
			}
		}

		if !testDbInstance {
			continue
		}
		if d.State != "RUNNABLE" {
			continue
		}
		running[d.Name] = struct{}{}
	}

	for _, d := range found.Items {
		// don't delete replicas, we'll take care of that
		// when deleting the database they replicate
		if d.ReplicaConfiguration != nil {
			continue
		}
		log.Printf("Destroying SQL Instance (%s)", d.Name)

		// replicas need to be stopped and destroyed before destroying a master
		// instance. The ordering slice tracks replica databases for a given master
		// and we call destroy on them before destroying the master
		var ordering []string
		for _, replicaName := range d.ReplicaNames {
			// don't try to stop replicas that aren't running
			if _, ok := running[replicaName]; !ok {
				ordering = append(ordering, replicaName)
				continue
			}

			// need to stop replication before being able to destroy a database
			op, err := config.NewSqlAdminClient(config.userAgent).Instances.StopReplica(config.Project, replicaName).Do()

			if err != nil {
				log.Printf("error, failed to stop replica instance (%s) for instance (%s): %s", replicaName, d.Name, err)
				return nil
			}

			err = sqlAdminOperationWaitTime(config, op, config.Project, "Stop Replica", config.userAgent, 10*time.Minute)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("Replication operation not found")
				} else {
					log.Printf("Error waiting for sqlAdmin operation: %s", err)
					return nil
				}
			}

			ordering = append(ordering, replicaName)
		}

		// ordering has a list of replicas (or none), now add the primary to the end
		ordering = append(ordering, d.Name)

		for _, db := range ordering {
			// destroy instances, replicas first
			op, err := config.NewSqlAdminClient(config.userAgent).Instances.Delete(config.Project, db).Do()

			if err != nil {
				if strings.Contains(err.Error(), "409") {
					// the GCP api can return a 409 error after the delete operation
					// reaches a successful end
					log.Printf("Operation not found, got 409 response")
					continue
				}

				log.Printf("Error, failed to delete instance %s: %s", db, err)
				return nil
			}

			err = sqlAdminOperationWaitTime(config, op, config.Project, "Delete Instance", config.userAgent, 10*time.Minute)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("SQL instance not found")
					continue
				}
				log.Printf("Error, failed to delete instance %s: %s", db, err)
				return nil
			}
		}
	}

	return nil
}

func TestAccSqlDatabaseInstance_basicInferredName(t *testing.T) {
	// Randomness
	skipIfVcr(t)
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_basic2,
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicSecondGen(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicMSSQL(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	rootPassword := randString(t, 15)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_mssql, databaseName, rootPassword),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_update_mssql, databaseName, rootPassword),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_dontDeleteDefaultUserOnReplica(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + randString(t, 10)
	failoverName := "sql-instance-test-failover-" + randString(t, 10)
	// 1. Create an instance.
	// 2. Add a root@'%' user.
	// 3. Create a replica and assert it succeeds (it'll fail if we try to delete the root user thinking it's a
	//    default user)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withoutReplica(databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				PreConfig: func() {
					// Add a root user
					config := googleProviderConfig(t)
					user := sqladmin.User{
						Name:     "root",
						Host:     "%",
						Password: randString(t, 26),
					}
					op, err := config.NewSqlAdminClient(config.userAgent).Users.Insert(config.Project, databaseName, &user).Do()
					if err != nil {
						t.Errorf("Error while inserting root@%% user: %s", err)
						return
					}
					err = sqlAdminOperationWaitTime(config, op, config.Project, "Waiting for user to insert", config.userAgent, 10*time.Minute)
					if err != nil {
						t.Errorf("Error while waiting for user insert operation to complete: %s", err.Error())
					}
					// User was created, now create replica
				},
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_basic(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_secondary(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_secondary, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_deletionProtection(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "true"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "true"),
				Destroy:     true,
				ExpectError: regexp.MustCompile("Error, failed to delete instance because deletion_protection is set to true. Set it to false to proceed with instance deletion"),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "false"),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_maintenanceVersion(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion, databaseName),
				ExpectError: regexp.MustCompile(
					`.*Maintenance version \(MYSQL_5_7_37.R20210508.01_03\) must not be set.*`),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_checkServiceNetworking(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_checkServiceNetworking, databaseName, databaseName),
				ExpectError: regexp.MustCompile("Error, failed to create instance because the network doesn't have at least 1 private services connection. Please see https://cloud.google.com/sql/docs/mysql/private-ip#network_requirements for how to create this connection."),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_replica(t *testing.T) {
	t.Parallel()

	databaseID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_replica, databaseID, databaseID, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.replica1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
			{
				ResourceName:            "google_sql_database_instance.replica2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_slave(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)
	slaveID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_slave, masterID, slaveID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance_slave",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_highAvailability(t *testing.T) {
	t.Parallel()

	instanceID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_highAvailability, instanceID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_diskspecs(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_diskspecs, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_maintenance(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_upgrade(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settingsDowngrade(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

// GH-4222
func TestAccSqlDatabaseInstance_authNets(t *testing.T) {
	t.Parallel()

	databaseID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step2, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

// Tests that a SQL instance can be referenced from more than one other resource without
// throwing an error during provisioning, see #9018.
func TestAccSqlDatabaseInstance_multipleOperations(t *testing.T) {
	t.Parallel()

	databaseID, instanceID, userID := randString(t, 8), randString(t, 8), randString(t, 8)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_multipleOperations, databaseID, instanceID, userID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basic_with_user_labels(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels_update, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-private")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(t *testing.T) {
	// Service Networking
	skipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-private-allocated-ip-range")
	addressName_update := "tf-test-" + randString(t, 10) + "update"
	networkName_update := BootstrapSharedTestNetwork(t, "sql-instance-private-allocated-ip-range-update")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName_update, addressName_update),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(t *testing.T) {
	// Service Networking
	skipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-private-replica")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.replica1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(t *testing.T) {
	// Service Networking
	skipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-private-clone")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.clone1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_createFromBackup(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    randString(t, 10),
		"original_db_name": BootstrapSharedSQLInstanceBackupRun(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_restoreFromBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "restore_backup_context"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_backupUpdate(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    randString(t, 10),
		"original_db_name": BootstrapSharedSQLInstanceBackupRun(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_beforeBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_restoreFromBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "restore_backup_context"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicClone(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    randString(t, 10),
		"original_db_name": BootstrapSharedSQLInstanceBackupRun(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_basicClone(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_cloneWithSettings(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    randString(t, 10),
		"original_db_name": BootstrapSharedSQLInstanceBackupRun(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_cloneWithSettings(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func testAccSqlDatabaseInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := googleProviderConfig(t)
			if rs.Type != "google_sql_database_instance" {
				continue
			}

			_, err := config.NewSqlAdminClient(config.userAgent).Instances.Get(config.Project,
				rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Database Instance still exists")
			}
		}

		return nil
	}
}

func testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t *testing.T, instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		users, err := config.NewSqlAdminClient(config.userAgent).Users.List(config.Project, instance).Do()

		if err != nil {
			return fmt.Errorf("Could not list database users for %q: %s", instance, err)
		}

		for _, u := range users.Items {
			if u.Name == "root" && u.Host == "%" {
				return fmt.Errorf("%v@%v user still exists", u.Name, u.Host)
			}
		}

		return nil
	}
}

func TestAccSqlDatabaseInstance_BackupRetention(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_BackupRetention(masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_PointInTimeRecoveryEnabled(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_insights(t *testing.T) {
	t.Parallel()

	masterID := randInt(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_insights, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_encryptionKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    getTestProjectFromEnv(),
		"key_name":      "tf-test-key-" + randString(t, 10),
		"instance_name": "tf-test-sql-" + randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: Nprintf(
					testGoogleSqlDatabaseInstance_encryptionKey, context),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    getTestProjectFromEnv(),
		"key_name":      "tf-test-key-" + randString(t, 10),
		"instance_name": "tf-test-sql-" + randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: Nprintf(
					testGoogleSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion, context),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ActiveDirectory(t *testing.T) {
	t.Parallel()
	databaseName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-private-test-ad")
	addressName := "tf-test-" + randString(t, 10)
	rootPassword := randString(t, 15)
	adDomainName := BootstrapSharedTestADDomain(t, "test-domain", networkName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, addressName, rootPassword, adDomainName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance-with-ad",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_SqlServerAuditConfig(t *testing.T) {
	// Service Networking
	skipIfVcr(t)
	t.Parallel()
	databaseName := "tf-test-" + randString(t, 10)
	rootPassword := randString(t, 15)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-sqlserver-audit")
	bucketName := fmt.Sprintf("%s-%d", "tf-test-bucket", randInt(t))
	uploadInterval := "900s"
	retentionInterval := "86400s"
	bucketNameUpdate := fmt.Sprintf("%s-%d", "tf-test-bucket", randInt(t)) + "update"
	uploadIntervalUpdate := "1200s"
	retentionIntervalUpdate := "172800s"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerAuditConfig(networkName, addressName, databaseName, rootPassword, bucketName, uploadInterval, retentionInterval),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerAuditConfig(networkName, addressName, databaseName, rootPassword, bucketNameUpdate, uploadIntervalUpdate, retentionIntervalUpdate),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_Timezone(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)
	rootPassword := randString(t, 15)
	addressName := "tf-test-" + randString(t, 10)
	networkName := BootstrapSharedTestNetwork(t, "sql-instance-sqlserver-audit")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_Timezone(networkName, addressName, databaseName, rootPassword, "Pacific Standard Time"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_mysqlMajorVersionUpgrade(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3_update, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(context),
			},
			{
				ResourceName:            "google_sql_database_instance.mysql_pvp_instance_name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
		},
	})
}

func testAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "mysql_pvp_instance_name" {
  name             = "tf-test-mysql-pvp-instance-name%{random_suffix}"
  region           = "asia-northeast1"
  database_version = "MYSQL_8_0"
  root_password = "abcABC123!"
  settings {
    tier              = "db-f1-micro"
    password_validation_policy {
      min_length  = 6
      complexity  =  "COMPLEXITY_DEFAULT"
      reuse_interval = 2
      disallow_username_substring = true
      enable_password_policy = true
    }
  }
  deletion_protection =  "%{deletion_protection}"
}
`, context)
}

var testGoogleSqlDatabaseInstance_basic2 = `
resource "google_sql_database_instance" "instance" {
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic3 = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic3_update = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic_mssql = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  database_version    = "SQLSERVER_2019_STANDARD"
  root_password       = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    collation = "Polish_CI_AS"
  }
}
`

var testGoogleSqlDatabaseInstance_update_mssql = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  database_version    = "SQLSERVER_2019_STANDARD"
  root_password       = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    collation = "Polish_CI_AS"
    ip_configuration {
      ipv4_enabled = true
      require_ssl = true
    }
  }
}
`

func testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, addressRangeName, rootPassword, adDomainName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance-with-ad" {
  depends_on = [google_service_networking_connection.foobar]
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-2-7680"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }

    active_directory_config {
      domain = "%s"
    }
  }
}`, networkName, addressRangeName, databaseName, rootPassword, adDomainName)
}

func testGoogleSqlDatabaseInstance_SqlServerAuditConfig(networkName, addressName, databaseName, rootPassword, bucketName, uploadInterval, retentionInterval string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "gs-bucket" {
  name                      	= "%s"
  location                  	= "US"
  uniform_bucket_level_access = true
}

data "google_compute_network" "servicenet" {
  name = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
	depends_on = [google_service_networking_connection.foobar]
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    sql_server_audit_config {
      bucket = "gs://%s"
      retention_interval = "%s"
      upload_interval = "%s"
    }
  }
}
`, bucketName, networkName, addressName, databaseName, rootPassword, bucketName, retentionInterval, uploadInterval)
}

func testGoogleSqlDatabaseInstance_Timezone(networkName, addressName, databaseName, rootPassword, timezone string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}
resource "google_sql_database_instance" "instance" {
	depends_on = [google_service_networking_connection.foobar]
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    time_zone = "%s"
  }
}
`, networkName, addressName, databaseName, rootPassword, timezone)
}

func testGoogleSqlDatabaseInstanceConfig_withoutReplica(instanceName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}
`, instanceName)
}

func testGoogleSqlDatabaseInstanceConfig_withReplica(instanceName, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  database_version     = "MYSQL_5_7"
  master_instance_name = google_sql_database_instance.instance.name
  deletion_protection  = false

  replica_configuration {
    failover_target = "true"
  }

  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
  }
}
`, networkName, addressRangeName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = google_compute_global_address.foobar.name
    }
  }
}
`, networkName, addressRangeName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}
resource "google_sql_database_instance" "replica1" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s-replica1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = google_compute_global_address.foobar.name
    }
  }

  master_instance_name = google_sql_database_instance.instance.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}
`, networkName, addressRangeName, databaseName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "clone1" {
  name                = "%s-clone1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  clone {
    source_instance_name = google_sql_database_instance.instance.name
    allocated_ip_range   = google_compute_global_address.foobar.name
  }

}
`, networkName, addressRangeName, databaseName, databaseName)
}

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
    location_preference {
      zone = "us-central1-f"
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ALWAYS"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_secondary = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
    location_preference {
      zone           = "us-central1-f"
	  secondary_zone = "us-central1-a"	  
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ALWAYS"
    connector_enforcement = "REQUIRED"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_deletionProtection = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = %s
  settings {
    tier                   = "db-f1-micro"
    location_preference {
      zone = "us-central1-f"
	}

    ip_configuration {
	  ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ALWAYS"
  }
}
`
var testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  maintenance_version = "MYSQL_5_7_37.R20210508.01_03"
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_checkServiceNetworking = `
resource "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = "false"
      private_network = google_compute_network.servicenet.self_link
    }
  }
}
`

var testGoogleSqlDatabaseInstance_replica = `
resource "google_sql_database_instance" "instance_master" {
  name                = "tf-lw-%d"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "replica1" {
  name                = "tf-lw-%d-1"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"
		backup_configuration {
      binary_log_enabled = true
		}
  }

  master_instance_name = google_sql_database_instance.instance_master.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}

resource "google_sql_database_instance" "replica2" {
  name                = "tf-lw-%d-2"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"
  }

  master_instance_name = google_sql_database_instance.instance_master.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}
`

var testGoogleSqlDatabaseInstance_slave = `
resource "google_sql_database_instance" "instance_master" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "instance_slave" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  master_instance_name = google_sql_database_instance.instance_master.name

  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_highAvailability = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    availability_type = "REGIONAL"

    backup_configuration {
      enabled  = true
      location = "us"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_diskspecs = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier                  = "db-f1-micro"
    disk_autoresize       = true
    disk_autoresize_limit = 50
    disk_size             = 15
    disk_type             = "PD_HDD"
  }
}
`

var testGoogleSqlDatabaseInstance_maintenance = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "canary"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_authNets_step1 = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"

    ip_configuration {
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }
  }
}
`

var testGoogleSqlDatabaseInstance_authNets_step2 = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-lw-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = "true"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_multipleOperations = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
  }
}

resource "google_sql_database" "database" {
  name     = "tf-test-%s"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_user" "user" {
  name     = "tf-test-%s"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "hunter2"
}
`

var testGoogleSqlDatabaseInstance_basic_with_user_labels = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    user_labels = {
      track    = "production"
      location = "western-division"
    }
  }
}
`
var testGoogleSqlDatabaseInstance_basic_with_user_labels_update = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    user_labels = {
      track = "production"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_insights = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 256
      record_application_tags = true
      record_client_address   = true
      query_plans_per_minute  = 10
    }
  }
}
`
var testGoogleSqlDatabaseInstance_encryptionKey = `
data "google_project" "project" {
  project_id = "%{project_id}"
}
resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "master" {
  name                = "%{instance_name}-master"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false
  encryption_key_name = google_kms_crypto_key.key.id

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "replica" {
  name                 = "%{instance_name}-replica"
  database_version     = "MYSQL_5_7"
  region               = "us-central1"
  master_instance_name = google_sql_database_instance.master.name
  deletion_protection  = false

  settings {
    tier = "db-n1-standard-1"
  }

  depends_on = [google_sql_database_instance.master]
}
`

var testGoogleSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion = `

data "google_project" "project" {
  project_id = "%{project_id}"
}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {

  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "master" {
  name                = "%{instance_name}-master"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false
  encryption_key_name = google_kms_crypto_key.key.id

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_kms_key_ring" "keyring-rep" {

  name     = "%{key_name}-rep"
  location = "us-east1"
}

resource "google_kms_crypto_key" "key-rep" {

  name     = "%{key_name}-rep"
  key_ring = google_kms_key_ring.keyring-rep.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_rep" {
  crypto_key_id = google_kms_crypto_key.key-rep.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "replica" {
  name                 = "%{instance_name}-replica"
  database_version     = "MYSQL_5_7"
  region               = "us-east1"
  master_instance_name = google_sql_database_instance.master.name
  encryption_key_name = google_kms_crypto_key.key-rep.id
  deletion_protection  = false

  settings {
    tier = "db-n1-standard-1"
  }

  depends_on = [google_sql_database_instance.master]
}
`

func testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID int, pointInTimeRecoveryEnabled bool) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    backup_configuration {
      enabled                        = true
      start_time                     = "00:00"
      point_in_time_recovery_enabled = %t
    }
  }
}
`, masterID, pointInTimeRecoveryEnabled)
}

func testGoogleSqlDatabaseInstance_BackupRetention(masterID int) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    backup_configuration {
      enabled                        = true
      start_time                     = "00:00"
      binary_log_enabled             = true
	  transaction_log_retention_days = 2
	  backup_retention_settings {
	    retained_backups = 4
	  }
    }
  }
}
`, masterID)
}

func testAccSqlDatabaseInstance_beforeBackup(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  deletion_protection = false
}
`, context)
}

func testAccSqlDatabaseInstance_restoreFromBackup(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  restore_backup_context {
    backup_run_id = data.google_sql_backup_run.backup.backup_id
    instance_id = data.google_sql_backup_run.backup.instance
  }

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [restore_backup_context[0].backup_run_id]
  }

  deletion_protection = false
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func testAccSqlDatabaseInstance_basicClone(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
  }

  deletion_protection = false

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [clone[0].point_in_time]
  }
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func testAccSqlDatabaseInstance_cloneWithSettings(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = false
	}
  }

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
  }

  deletion_protection = false

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [clone[0].point_in_time]
  }
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}
