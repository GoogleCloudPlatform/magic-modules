package spanner_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/spanner"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccSpannerInstancePartition_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_basic(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_update(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstancePartition_processingUnits(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_processingUnits(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_processingUnitsUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "test-spanner-partition"
  node_count   = 1
}
`, context)
}

func testAccSpannerInstancePartition_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "updated-spanner-partition"
  node_count   = 2
}
`, context)
}

func testAccSpannerInstancePartition_processingUnits(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
  edition          = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "nam8"
  display_name     = "test-spanner-partition"
  processing_units = 1000
}
`, context)
}

func testAccSpannerInstancePartition_processingUnitsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
  edition          = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "nam8"
  display_name     = "updated-spanner-partition"
  processing_units = 2000
}
`, context)
}

func TestAccSpannerInstancePartition_autoscaling(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_autoscaling(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_autoscalingUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_autoscaling(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "test-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_processing_units = 1000
      max_processing_units = 2000
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      storage_utilization_percent           = 95
    }
  }
}
`, context)
}

func testAccSpannerInstancePartition_autoscalingUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "updated-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_processing_units = 1000
      max_processing_units = 3000
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 70
      storage_utilization_percent           = 90
    }
  }
}
`, context)
}

func TestAccSpannerInstancePartition_autoscalingWithNodes(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_autoscalingWithNodes(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_autoscalingWithNodesUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_autoscalingWithNodes(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "test-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_nodes = 1
      max_nodes = 2
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      storage_utilization_percent           = 95
    }
  }
}
`, context)
}

func testAccSpannerInstancePartition_autoscalingWithNodesUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "updated-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_nodes = 1
      max_nodes = 3
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 70
      storage_utilization_percent           = 90
    }
  }
}
`, context)
}

func TestAccSpannerInstancePartition_asymmetricAutoscaling(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_asymmetricAutoscaling(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_asymmetricAutoscalingUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_asymmetricAutoscaling(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "test-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_processing_units = 1000
      max_processing_units = 3000
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      total_cpu_utilization_percent         = 80
      storage_utilization_percent           = 90
    }
  }
}
`, context)
}

func testAccSpannerInstancePartition_asymmetricAutoscalingUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "nam8"
  display_name = "updated-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
      min_processing_units = 2000
      max_processing_units = 4000
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 70
      total_cpu_utilization_percent         = 85
      storage_utilization_percent           = 95
    }
  }
}
`, context)
}

func TestAccSpannerInstancePartition_virtualUpdate(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_virtualUpdate(suffix, false),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_virtualUpdate(suffix, true),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
				// force_destroy is client-side only, so import resets it to its default.
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testAccSpannerInstancePartition_virtualUpdate(suffix string, forceDestroy bool) string {
	context := map[string]interface{}{
		"random_suffix": suffix,
		"force_destroy": forceDestroy,
	}
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}

resource "google_spanner_instance_partition" "partition" {
  name          = "tf-test-partition-%{random_suffix}"
  instance      = google_spanner_instance.main.name
  config        = "nam8"
  display_name  = "test-spanner-partition"
  node_count    = 1
  force_destroy = %{force_destroy}
}
`, context)
}

func TestAccSpannerInstancePartition_forceDestroy(t *testing.T) {
	// Uses time.Now
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_forceDestroy(context),
				// The backup blocks the partition delete unless force_destroy removes it.
				Check: testAccSpannerInstancePartitionCreateBackup(t, "google_spanner_database.database"),
			},
		},
	})
}

// testAccSpannerInstancePartitionCreateBackup backs up the database and waits for the backup to be ready.
// A database cannot be deleted while one of its backups is still being created.
func testAccSpannerInstancePartitionCreateBackup(t *testing.T, databaseResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[databaseResource]
		if !ok {
			return fmt.Errorf("Not found: %s", databaseResource)
		}

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(spanner.Product, config)+"projects/{{project}}/instances/{{instance}}/backups?backupId={{name}}")
		if err != nil {
			return err
		}
		database, err := tpgresource.ReplaceVarsForTest(config, rs, "projects/{{project}}/instances/{{instance}}/databases/{{name}}")
		if err != nil {
			return err
		}

		project := rs.Primary.Attributes["project"]
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body: map[string]interface{}{
				"database": database,
				// Just over the 6 hour minimum, so that a leaked backup expires soon.
				"expireTime": time.Now().Add(7 * time.Hour).UTC().Format(time.RFC3339),
			},
		})
		if err != nil {
			return fmt.Errorf("Error creating backup of %s: %s", database, err)
		}
		return spanner.SpannerOperationWaitTime(config, res, project, "Creating Backup", config.UserAgent, 60*time.Minute)
	}
}

func testAccSpannerInstancePartition_forceDestroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"

  # A scheduled backup that is still being created would block the database delete.
  default_backup_schedule_type = "NONE"
}

resource "google_spanner_instance_partition" "partition" {
  name          = "tf-test-partition-%{random_suffix}"
  instance      = google_spanner_instance.main.name
  config        = "nam8"
  display_name  = "test-spanner-partition"
  node_count    = 1
  force_destroy = true
}

resource "google_spanner_database" "database" {
  instance            = google_spanner_instance.main.name
  name                = "tf-test-db-%{random_suffix}"
  deletion_protection = false
  ddl = [
    "CREATE PLACEMENT testplacement OPTIONS (instance_partition = '${google_spanner_instance_partition.partition.name}')",
    "CREATE TABLE Singers (SingerId INT64 NOT NULL, Location STRING(MAX) NOT NULL PLACEMENT KEY) PRIMARY KEY (SingerId)",
  ]
}
`, context)
}
