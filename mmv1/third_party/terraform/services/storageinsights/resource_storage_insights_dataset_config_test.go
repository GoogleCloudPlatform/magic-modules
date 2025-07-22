package storageinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_project_to_org(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_project(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_org(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_folder_to_org(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_folder(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_org(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_project_to_folder(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_project(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_folder(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_org_to_folder(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_org(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_folder(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_org_to_project(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_org(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_project(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func TestAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_folder_to_project(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_folder(context),
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
			{
				Config: testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_project(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_insights_dataset_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_insights_dataset_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_config_id", "location"},
			},
		},
	})
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_project(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 1
    source_projects {
        project_numbers = ["123", "456", "789"]
    }
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
}
`, context)
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_folder(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 1
    source_folders {
        folder_numbers = ["123", "456", "789"]
    }
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
}
`, context)
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_full_org(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 1
    organization_scope = true
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
}
`, context)
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_project(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 2
    source_projects {
		project_numbers = ["123", "456"]
	}
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
	description = "A sample description for dataset"
	link_dataset = false
}
`, context)
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_folder(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 2
    source_folders {
		folder_numbers = ["123", "456"]
	}
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
	description = "A sample description for dataset"
	include_newly_created_buckets = true
	exclude_cloud_storage_locations {
		locations = ["us-east1", "europe-west2"]
	}
	include_cloud_storage_buckets {
		cloud_storage_buckets {
			bucket_name = "gs://sample-bucket1/"
		}
		cloud_storage_buckets {
			bucket_prefix_regex = "gs://sample*/"
		}
	}
}
`, context)
}

func testAccStorageInsightsDatasetConfig_storageInsightsDatasetConfigExample_update_org(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_insights_dataset_config" "config" {
    location = "us-central1"
    dataset_config_id = "tf_test_my_config%{random_suffix}"
    retention_period_days = 2
    organization_scope = true
    identity {
        type = "IDENTITY_TYPE_PER_CONFIG"
    }
	description = "A sample description for dataset"
	include_newly_created_buckets = true
	include_cloud_storage_locations {
		locations = ["us-east1", "europe-west2"]
	}
	exclude_cloud_storage_buckets {
		cloud_storage_buckets {
			bucket_name = "gs://sample-bucket1/"
		}
		cloud_storage_buckets {
			bucket_prefix_regex = "gs://sample*/"
		}
	}
}
`, context)
}
