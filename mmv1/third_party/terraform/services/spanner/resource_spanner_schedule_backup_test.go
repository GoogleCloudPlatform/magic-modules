package spanner_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests

func TestAccSpannerBackupSchedule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerBackupScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerBackupSchedule_basic(context),
			},
			{
				ResourceName:      "google_spanner_backup_schedule.backup_schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerBackupSchedule_update(context),
			},
			{
				ResourceName:      "google_spanner_backup_schedule.backup_schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerBackupSchedule_CMEKIncrementalBackup(t *testing.T) {
	t.Parallel()
	suffix := acctest.RandString(t, 10)
	kms := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-spanner-key")

	context := map[string]interface{}{
		"random_suffix": suffix,
		"key_name":      kms.CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerBackupScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerBackupSchedule_CMEKIncremental(context),
			},
			{
				ResourceName:      "google_spanner_backup_schedule.backup_schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerBackupSchedule_CMEKFullBackup(t *testing.T) {
	t.Parallel()
	suffix := acctest.RandString(t, 10)
	kms := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-spanner-key")

	context := map[string]interface{}{
		"random_suffix": suffix,
		"key_name":      kms.CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerBackupScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerBackupSchedule_basic(context),
			},
			{
				ResourceName:      "google_spanner_backup_schedule.backup_schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerBackupSchedule_CMEKFull(context),
			},
			{
				ResourceName:      "google_spanner_backup_schedule.backup_schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerBackupSchedule_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "instance" {
  name         = "my-instance-%{random_suffix}"
  config       = "regional-us-central1"
  display_name = "My Instance"
  num_nodes    = 1
  edition      = "ENTERPRISE"
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "my-database-%{random_suffix}"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
  deletion_protection = false
}

resource "google_spanner_backup_schedule" "backup_schedule" {
  instance = google_spanner_instance.instance.name
  database = google_spanner_database.database.name
  name     = "my-backup-schedule-%{random_suffix}"

  retention_duration = "172800s"

  spec {
    cron_spec {
      text = "0 12 * * *"
    }
  }

  full_backup_spec {}
}
`, context)
}

func testAccSpannerBackupSchedule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "instance" {
  name         = "my-instance-%{random_suffix}"
  config       = "regional-us-central1"
  display_name = "My Instance"
  num_nodes    = 1
  edition      = "ENTERPRISE"
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "my-database-%{random_suffix}"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
  deletion_protection = false
}

resource "google_spanner_backup_schedule" "backup_schedule" {
  instance = google_spanner_instance.instance.name
  database = google_spanner_database.database.name
  name     = "my-backup-schedule-%{random_suffix}"

  retention_duration = "172900s"

  spec {
    cron_spec {
      text = "0 0 * * *"
    }
  }

  full_backup_spec {}
}
`, context)
}

func testAccSpannerBackupSchedule_CMEKIncremental(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "instance" {
  name         = "my-instance-%{random_suffix}"
  config       = "regional-us-central1"
  display_name = "My Instance"
  num_nodes    = 1
  edition      = "ENTERPRISE"
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "my-database-%{random_suffix}"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
  deletion_protection = false

  encryption_config {
    kms_key_name = "%{key_name}"
  }
}

resource "google_spanner_backup_schedule" "backup_schedule" {
  instance = google_spanner_instance.instance.name
  database = google_spanner_database.database.name
  name     = "my-backup-schedule-%{random_suffix}"

  retention_duration = "172800s"

  spec {
    cron_spec {
      text = "0 12 * * *"
    }
  }

  incremental_backup_spec {}

  encryption_config {
    encryption_type = "GOOGLE_DEFAULT_ENCRYPTION"
  }
}
`, context)
}

func testAccSpannerBackupSchedule_CMEKFull(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "instance" {
  name         = "my-instance-%{random_suffix}"
  config       = "regional-us-central1"
  display_name = "My Instance"
  num_nodes    = 1
  edition      = "ENTERPRISE"
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "my-database-%{random_suffix}"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
  deletion_protection = false
}

resource "google_spanner_backup_schedule" "backup_schedule" {
  instance = google_spanner_instance.instance.name
  database = google_spanner_database.database.name
  name     = "my-backup-schedule-%{random_suffix}"

  retention_duration = "172800s"

  spec {
    cron_spec {
      text = "0 12 * * *"
    }
  }

  full_backup_spec {}

  encryption_config {
    encryption_type = "CUSTOMER_MANAGED_ENCRYPTION"
    kms_key_name    = "%{key_name}"
  }
}
`, context)
}
