package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeSnapshot_encryption(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryption(snapshotName, diskName),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_encryption_key", "source_disk", "source_disk_encryption_key", "zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_encryptionCMEK(t *testing.T) {
	t.Parallel()
	// KMS causes errors due to rotation
	acctest.SkipIfVcr(t)

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	kmsKeyName := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key1").CryptoKey.Name

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryptionCMEK(snapshotName, diskName, kmsKeyName),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "snapshot_encryption_key", "source_disk_encryption_key"},
			},
		},
	})
}

func TestAccComputeSnapshot_encryptionRSA(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"resource_id":       acctest.RandString(t, 10),
		"rsa_encrypted_key": "fB6BS8tJGhGVDZDjGt1pwUo2wyNbkzNxgH1avfOtiwB9X6oPG94gWgenygitnsYJyKjdOJ7DyXLmxwQOSmnCYCUBWdKCSssyLV5907HL2mb5TfqmgHk5JcArI/t6QADZWiuGtR+XVXqiLa5B9usxFT2BTmbHvSKfkpJ7McCNc/3U0PQR8euFRZ9i75o/w+pLHFMJ05IX3JB0zHbXMV173PjObiV3ItSJm2j3mp5XKabRGSA5rmfMnHIAMz6stGhcuom6+bMri2u/axmPsdxmC6MeWkCkCmPjaKsVz1+uQUNCJkAnzesluhoD+R6VjFDm4WI7yYabu4MOOAOTaQXdEg==",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryptionRSA(context),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_instantSnapshot(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"resource_id":       acctest.RandString(t, 10),
		"kms_key_self_link": acctest.BootstrapKMSKey(t).CryptoKey.Name,
		"raw_key":           "",
		"rsa_encrypted_key": "fB6BS8tJGhGVDZDjGt1pwUo2wyNbkzNxgH1avfOtiwB9X6oPG94gWgenygitnsYJyKjdOJ7DyXLmxwQOSmnCYCUBWdKCSssyLV5907HL2mb5TfqmgHk5JcArI/t6QADZWiuGtR+XVXqiLa5B9usxFT2BTmbHvSKfkpJ7McCNc/3U0PQR8euFRZ9i75o/w+pLHFMJ05IX3JB0zHbXMV173PjObiV3ItSJm2j3mp5XKabRGSA5rmfMnHIAMz6stGhcuom6+bMri2u/axmPsdxmC6MeWkCkCmPjaKsVz1+uQUNCJkAnzesluhoD+R6VjFDm4WI7yYabu4MOOAOTaQXdEg==",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_instantSnapshot(context),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot-raw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot-rsa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot-kms",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_guestOSFeatures(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"resource_id":       acctest.RandString(t, 10),
		"guest_os_features": "",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_guestOsFeatures(context),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_snapshotType(t *testing.T) {
	t.Parallel()

	context_1 := map[string]any{
		"resource_id":   acctest.RandString(t, 10),
		"snapshot_type": "",
	}
	context_2 := map[string]any{
		"resource_id":   context_1["resource_id"],
		"snapshot_type": "",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_snapshotType(context_1),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
			{
				Config: testAccComputeSnapshot_snapshotType(context_2),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"architecture", "labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "snapshot_type", "source_disk", "source_disk_encryption_key", "source_disk_for_recovery_checkpoint", "source_instant_snapshot", "source_instant_snapshot_encryption_key.0.raw_key", "source_instant_snapshot_encryption_key.0.rsa_encrypted_key", "terraform_labels", "zone"},
			},
		},
	})
}

func testAccComputeSnapshot_encryption(snapshotName string, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }

  source_disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
`, diskName, snapshotName)
}

func testAccComputeSnapshot_encryptionCMEK(snapshotName, diskName, kmsKeyName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "KMS Ops Account"
}

resource "google_kms_crypto_key_iam_member" "example-key" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${google_service_account.test.email}"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
    kms_key_service_account = google_service_account.test.email
  }
  depends_on = [google_kms_crypto_key_iam_member.example-key]
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    kms_key_self_link = "%s"
    kms_key_service_account = google_service_account.test.email
  }
}
`, diskName, kmsKeyName, diskName, kmsKeyName, snapshotName, kmsKeyName)
}

func testAccComputeSnapshot_encryptionRSA(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-disk-%{resource_id}"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    rsa_encrypted_key = "%{rsa_encrypted_key}"
  }
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    rsa_encrypted_key = "%{rsa_encrypted_key}"
  }

  source_disk_encryption_key {
    rsa_encrypted_key = "%{rsa_encrypted_key}"
  }
}
`, context)
}

func testAccComputeSnapshot_instantSnapshot(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-disk-%{resource_id}"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    raw_key = "%{raw_key}"
  }
}

resource "google_compute_instant_snapshot" "foobar-raw" {
	name = "tf-test-isnapshot1-%{resource_id}"
  	source_disk = "google_compute_disk.foobar.id"

	instant_snapshot_encryption_key {
		raw_key = %{raw_key}
	}
}

resource "google_compute_instant_snapshot" "foobar-rsa" {
	name = "tf-test-isnapshot2-%{resource_id}"
	source_disk = "google_compute_disk.foobar.id"

	instant_snapshot_encryption_key {
		rsa_encrypted_key = %{rsa_encrypted_key}
	}
}

resource "google_compute_instant_snapshot" "foobar-kms" {
	name = "tf-test-isnapshot3-%{resource_id}"
	source_disk = "google_compute_disk.foobar.id"

	instant_snapshot_encryption_key {
		kms_key_self_link = %{kms_key_self_link}
	}
}

resource "google_compute_snapshot" "foobar-raw" {
	name = "tf-test-snapshot1-%{resource_id}"
	source_instant_snapshot = google_compute_instant_snapshot.foobar-raw.id

	source_instant_snapshot_encryption_key {
		raw_key = %{raw_key}
	}

	snapshot_encryption_key {
		raw_key = %{raw_key}
	}
}

resource "google_compute_snapshot" "foobar-rsa" {
	name = "tf-test-snapshot2-%{resource_id}"
	source_instant_snapshot = google_compute_instant_snapshot.foobar-rsa.id

	source_instant_snapshot_encryption_key {
		rsa_encrypted_key = %{rsa_encrypted_key}
	}

	snapshot_encryption_key {
		rsa_encrypted_key = %{rsa_encrypted_key}
	}
}

resource "google_compute_snapshot" "foobar-kms" {
	name = "tf-test-snapshot3-%{resource_id}"
	source_instant_snapshot = google_compute_instant_snapshot.foobar-kms.id

	source_instant_snapshot_encryption_key {
		kms_key_self_link = %{kms_key_self_link}
	}

	snapshot_encryption_key {
		kms_key_self_link = %{kms_key_self_link}
	}
}
`, context)
}

func testAccComputeSnapshot_snapshotType(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-disk-%{resource_id}"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
  name        = "tf-test-snapshot-%{resource_id}"
  source_disk = google_compute_disk.foobar.name
  snapshot_type = "%{snapshot_type}"
}
`, context)
}

func testAccComputeSnapshot_guestOsFeatures(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-disk-%{resource_id}"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
  name        = "tf-test-snapshot-%{resource_id}"
  source_disk = google_compute_disk.foobar.name

  guest_os_features {
	type = "%{guest_os_features}"
  }
`, context)
}
