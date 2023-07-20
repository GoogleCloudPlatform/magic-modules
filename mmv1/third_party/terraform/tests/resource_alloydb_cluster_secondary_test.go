package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbClusterSecondary_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_update(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
		},
	})
}

func testAccAlloydbClusterSecondary_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  labels = {
	foo = "bar" 
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// Test if adding automatedBackupPolicy AND initialUser re-creates the cluster.
// Ideally, cluster shouldn't be re-created. This test will only pass if the cluster
// isn't re-created but updated in-place.
func TestAccAlloydbClusterSecondary_addAutomatedBackupPolicyAndInitialUser(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_withoutInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_withInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
		},
	})
}

// Test if deleting automatedBackupPolicy AND initialUser re-creates the cluster.
// Ideally, cluster shouldn't be re-created. This test will only pass if the cluster
// isn't re-created but updated in-place.
func TestAccAlloydbClusterSecondary_deleteAutomatedBackupPolicyAndInitialUser(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_withInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_withoutInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
		},
	})
}

func testAccAlloydbClusterSecondary_withInitialUserAndAutomatedBackupPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-west1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
  
  initial_user {
    user     = "tf-test-alloydb-cluster%{random_suffix}"
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = 23
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }

    quantity_based_retention {
      count = 1
    }

    labels = {
      test = "tf-test-alloydb-cluster%{random_suffix}"
    }
  }
  lifecycle {
    prevent_destroy = true
  }  
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func testAccAlloydbClusterSecondary_withoutInitialUserAndAutomatedBackupPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }

  lifecycle {
    prevent_destroy = true
  }  
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// The cluster creation should work fine even without a weekly schedule.
func TestAccAlloydbClusterSecondary_missingWeeklySchedule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_missingWeeklySchedule(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlloydbClusterSecondary_missingWeeklySchedule(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
 
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
 
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true
    quantity_based_retention {
	  count = 1
	}
    labels = {
	  test = "tf-test-alloydb-cluster%{random_suffix}"
	}
  }
}
data "google_project" "project" {}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// The cluster creation should succeed with minimal number of arguments.
func TestAccAlloydbClusterSecondary_mandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
		},
	})
}

// Deletion of time-based retention policy should be an in-place operation
func TestAccAlloydbClusterSecondary_deleteTimeBasedRetentionPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_withTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbClusterSecondary_withoutTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbClusterSecondary_AlloydbClusterSecondaryBasicExample(context),
			},
		},
	})
}

func testAccAlloydbClusterSecondary_withTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
 
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
 
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = 23
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
    time_based_retention {
      retention_period = "4.5s"
    }
  }
  lifecycle {
    ignore_changes = [
      automated_backup_policy[0].time_based_retention
    ]
    prevent_destroy = true
  }
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func testAccAlloydbClusterSecondary_withoutTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
 
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = 23
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }
  lifecycle {
    ignore_changes = [
      automated_backup_policy[0].time_based_retention
    ]
    prevent_destroy = true
  }
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}
func TestAccAlloydbClusterSecondary_usingCMEK(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_usingCMEK(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
		},
	})
}

func testAccAlloydbClusterSecondary_usingCMEK(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
 
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
 
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
data "google_project" "project" {}
resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-west1"
}  
resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}
resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}
`, context)
}

func TestAccAlloydbClusterSecondary_CMEKInAutomatedBackupIsUpdatable(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterSecondaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterSecondary_usingCMEKInClusterAndAutomatedBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_updateCMEKInAutomatedBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterSecondary_usingCMEKallowDeletion(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
		},
	})
}

func testAccAlloydbClusterSecondary_usingCMEKInClusterAndAutomatedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-west1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}
`, context)
}

func testAccAlloydbClusterSecondary_updateCMEKInAutomatedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  } 
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key2.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-west1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key" "key2" {
	name     = "%{key_name}-2"
	key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	members = [
	  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
	]
}
`, context)
}

func testAccAlloydbClusterSecondary_usingCMEKallowDeletion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id  = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location    = "us-central1"
  network     = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-west1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  secondary_config {
    primary_cluster_name = "${google_alloydb_cluster.primary.cluster_id}"
  }
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-west1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key2.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-west1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key" "key2" {
	name     = "%{key_name}-2"
	key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	members = [
	  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
	]
}
`, context)
}
