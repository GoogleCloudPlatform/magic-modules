package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatabaseMigrationServiceConversionWorkspace_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceConversionWorkspaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConversionWorkspace_basic(context),
			},
			{
				ResourceName:            "google_database_migration_service_conversion_workspace.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conversion_workspace_id", "location"},
			},
		},
	})
}

func TestAccDatabaseMigrationServiceConversionWorkspace_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceConversionWorkspaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConversionWorkspace_basic(context),
			},
			{
				ResourceName:            "google_database_migration_service_conversion_workspace.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conversion_workspace_id", "location"},
			},
			{
				Config: testAccDatabaseMigrationServiceConversionWorkspace_update(context),
			},
			{
				ResourceName:            "google_database_migration_service_conversion_workspace.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conversion_workspace_id", "location"},
			},
		},
	})
}

func TestAccDatabaseMigrationServiceConversionWorkspace_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceConversionWorkspaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConversionWorkspace_full(context),
			},
			{
				ResourceName:            "google_database_migration_service_conversion_workspace.example_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conversion_workspace_id", "location"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConversionWorkspace_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_conversion_workspace" "example" {
	location = "us-central1"
	conversion_workspace_id = "tf-test-conversion-workspace%{random_suffix}"
	display_name = "Test conversion workspace"
	
	source {
		engine = "ORACLE"
		version = "21c"
	}
	
	destination {
		engine = "POSTGRESQL"
		version = "15"
	}
}
`, context)
}

func testAccDatabaseMigrationServiceConversionWorkspace_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_conversion_workspace" "example" {
	location = "us-central1"
	conversion_workspace_id = "tf-test-conversion-workspace%{random_suffix}"
	display_name = "Updated conversion workspace"
	
	source {
		engine = "ORACLE"
		version = "21c"
	}
	
	destination {
		engine = "POSTGRESQL"
		version = "15"
	}
	
	global_settings = {
		"skip_triggers" = "true"
		"max_parallel_level" = "10"
	}
}
`, context)
}

func testAccDatabaseMigrationServiceConversionWorkspace_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_conversion_workspace" "example_full" {
	location = "us-central1"
	conversion_workspace_id = "tf-test-conversion-workspace-full%{random_suffix}"
	display_name = "Full conversion workspace example"
	
	source {
		engine = "ORACLE"
		version = "19c"
	}
	
	destination {
		engine = "POSTGRESQL"
		version = "15"
	}
	
	global_settings = {
		"convert_foreign_key_to_interleave" = "true"
		"skip_triggers" = "false"
		"ignore_non_table_synonyms" = "true"
		"max_parallel_level" = "5"
	}
}
`, context)
}
