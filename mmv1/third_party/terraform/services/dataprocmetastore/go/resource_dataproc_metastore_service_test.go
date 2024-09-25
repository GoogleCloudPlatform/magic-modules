package dataprocmetastore_test

import (
	"fmt"
	"testing"
	"regexp"
	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataprocMetastoreService_updateAndImport(t *testing.T) {
	t.Parallel()

	name := "tf-test-metastore-" + acctest.RandString(t, 10)
	tier := [2]string{"DEVELOPER", "ENTERPRISE"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[0]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[1]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccDataprocMetastoreService_updateLocation_deletionProtection(t *testing.T) {
        t.Parallel()

        name := "tf-test-metastore-" + acctest.RandString(t, 10)
        tier := [2]string{"DEVELOPER", "ENTERPRISE"}

        acctest.VcrTest(t, resource.TestCase{
                PreCheck:                 func() { acctest.AccTestPreCheck(t) },
                ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
                Steps: []resource.TestStep{
                        {
                                Config: testAccDataprocMetastoreService_updateLocation_deletionProtection(name, "us-central1", tier[0]),
                        },
                        {
                                ResourceName:            "google_dataproc_metastore_service.my_metastore",
                                ImportState:             true,
                                ImportStateVerify:       true,
                                ImportStateVerifyIgnore: []string{"deletion_protection"},
                        },
                        {
                                Config: testAccDataprocMetastoreService_updateLocation_deletionProtection(name, "us-west2", tier[1]),
                                ExpectError: regexp.MustCompile("deletion_protection"),
                        },
                        {
                                Config: testAccDataprocMetastoreService_updateAndImport(name, tier[1]),
                        },
                },
        })
}

func testAccDataprocMetastoreService_updateAndImport(name, tier string) string {
	return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
	service_id = "%s"
	location   = "us-central1"
	tier       = "%s"
	deletion_protection = true

	hive_metastore_config {
		version = "2.3.6"
	}
}
`, name, tier)
}

func testAccDataprocMetastoreService_updateLocation_deletionProtection(name, location, tier string) string {
        return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
        service_id = "%s"
        location   = "%s"
        tier       = "%s"
        deletion_protection = true

        hive_metastore_config {
                version = "2.3.6"
        }
}
`, name, location, tier)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location", "labels", "terraform_labels", "deletion_protection"},
			},
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context),
			},
		},
	})
}

func TestAccDataprocMetastoreService_PrivateServiceConnect(t *testing.T) {
	t.Skip("Skipping due to https://github.com/hashicorp/terraform-provider-google/issues/13710")
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_PrivateServiceConnect(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location"},
			},
		},
	})
}

func testAccDataprocMetastoreService_PrivateServiceConnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
// Use data source instead of creating a subnetwork due to a bug on API side.
// With the bug, the new created subnetwork cannot be deleted when deleting the dataproc metastore service.
data "google_compute_subnetwork" "subnet" {
  name   = "default"
  region = "us-central1"
}

resource "google_dataproc_metastore_service" "default" {
  service_id = "tf-test-metastore-srv%{random_suffix}"
  location   = "us-central1"

  hive_metastore_config {
    version = "3.1.2"
  }

  network_config {
    consumers {
      subnetwork = data.google_compute_subnetwork.subnet.id
    }
  }
}
`, context)
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "backup" {
  service_id = "tf-test-backup%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  scheduled_backup {
    enabled         = true
    cron_schedule   = "0 0 * * 0"
    time_zone       = "America/Los_Angeles"
    backup_location = "gs://${google_storage_bucket.bucket.name}"
  }

  labels = {
    env = "test"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-backup%{random_suffix}"
  location = "us-central1"
}
`, context)
}
