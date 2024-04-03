package datalossprevention_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateOrg(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgRunning(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgFolderPaused(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateActions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActions(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateConditions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigConditionsCadence(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigFilterRegexesAndConditions(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_discovery_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"
	display_name = "display name"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
    inspect_templates = ["FAKE"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
			conditions {
				or_conditions {
					min_row_count = 10
					minAge = "10800s"
				}
			}
        }
    }
    inspect_templates = ["FAKE_NEW"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	actions {
        export_data {
            profile_table {
                project_id = "project"
                dataset_id = "dataset"
                table_id = "table"
            }
        }
    }
    actions { 
        pub_sub_notification {
                topic = "fake-topic"
                event = "NEW_PROFILE"
                pub_sub_condition {
                    expressions {
                        logical_operator = "OR"
                        conditions {
                            minimum_risk_score = "HIGH"
                            minimum_sensitivity_score = "HIGH"
                        }
                    }
                }
                detail_of_message = "TABLE_PROFILE"
            }
    }
    inspect_templates = ["FAKE"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgRunning(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "organizations/%{organization}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	org_config {
		project_id = "%{project}"
		location {
			organization_id = "%{organization}"
		}
	}
    inspect_templates = ["FAKE"]
	status = "RUNNING"
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgFolderPaused(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "organizations/%{organization}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	org_config {
		project_id = "%{project}"
		location {
			folder_id = 123
		}
	}
    inspect_templates = ["FAKE"]
	status = "PAUSED"
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigConditionsCadence(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"

	targets {
        big_query_target {
            filter {
                other_tables {}
            }
            conditions {
                type_collection = "BIG_QUERY_COLLECTION_ALL_TYPES"
            }
            cadence {
                schema_modified_cadence {
                    types = ["SCHEMA_NEW_COLUMNS"]
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                table_modified_cadence {
                    types = ["TABLE_MODIFIED_TIMESTAMP"]
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
            }
        }
    }
    inspect_templates = ["FAKE_NEW"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigFilterRegexesAndConditions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"

	targets {
        big_query_target {
            filter {
                tables {
                    include_regexes {
                        patterns {
                            project_id_regex = ".*"
                            dataset_id_regex = ".*"
                            table_id_regex = ".*"
                        }
                    }
                }
            }
            conditions {
                created_after = "2023-10-02T15:01:23Z"
                types {
                    types = ["BIG_QUERY_TABLE_TYPE_TABLE", "BIG_QUERY_TABLE_TYPE_EXTERNAL_BIG_LAKE"]
                }
                or_conditions {
                    min_row_count = 10
                    min_age = "10d"
                }
            }
        }
    }
    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
    inspect_templates = ["FAKE_NEW"]
}
`, context)
}
