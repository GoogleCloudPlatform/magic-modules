package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestSnapshotKmsKeyDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old                string
		New                string
		ExpectDiffSuppress bool
	}{
		"matching base keys": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			ExpectDiffSuppress: true,
		},
		"base key old, versioned new": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/1",
			ExpectDiffSuppress: true,
		},
		"versioned old, base key new": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/1",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			ExpectDiffSuppress: true,
		},
		"different versions same key": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/1",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/2",
			ExpectDiffSuppress: true,
		},
		"self link old, relative versioned new": {
			Old:                "https://www.googleapis.com/compute/v1/projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/1",
			ExpectDiffSuppress: true,
		},
		"different keys": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/other-key",
			ExpectDiffSuppress: false,
		},
		"different keys with versions": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key/cryptoKeyVersions/1",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/other-key/cryptoKeyVersions/1",
			ExpectDiffSuppress: false,
		},
		"empty old": {
			Old:                "",
			New:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			ExpectDiffSuppress: false,
		},
		"empty new": {
			Old:                "projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key",
			New:                "",
			ExpectDiffSuppress: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tpgcompute.SnapshotKmsKeyDiffSuppress("snapshot_encryption_key.0.kms_key_self_link", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("SnapshotKmsKeyDiffSuppress %q failed: old=%s, new=%s, expected=%v, got=%v",
					name, tc.Old, tc.New, tc.ExpectDiffSuppress, !tc.ExpectDiffSuppress)
			}
		})
	}
}

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
	kmsKeyName := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key1").CryptoKey.Name

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


func TestAccComputeSnapshot_kmsKeyVersionBehaviors(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t) // Based on the other test having this

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	kmsKey1 := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key3").CryptoKey.Name
	kmsKey2 := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key4").CryptoKey.Name

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey1),
			},
			{
				Config:             testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey1+"/cryptoKeyVersions/1"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:      testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey2+"/cryptoKeyVersions/1"),
				ExpectError: regexp.MustCompile("Changing the KMS key to a specific version is not supported\\."),
			},
		},
	})
}
func TestAccComputeSnapshot_updateKmsKey(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	kmsKey1 := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key1").CryptoKey.Name
	kmsKey2 := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key2").CryptoKey.Name

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			// Step 1: Create snapshot with CMEK Key 1
			{
				Config: testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_compute_snapshot.foobar", "snapshot_encryption_key.0.kms_key_self_link", kmsKey1),
				),
			},
			// Step 2: In-place update to change to CMEK Key 2
			{
				Config: testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_compute_snapshot.foobar", "snapshot_encryption_key.0.kms_key_self_link", kmsKey2),
				),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "source_disk_encryption_key"},
			},
		},
	})
}

func TestAccComputeSnapshot_disallowCmekToGmek(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	kmsKey := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-compute-snapshot-key1").CryptoKey.Name

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey),
			},
			{
				Config:      testAccComputeSnapshot_updateKmsKey_step1(snapshotName, diskName),
				ExpectError: regexp.MustCompile("Removing 'snapshot_encryption_key' is not supported in-place. To remove the KMS key, please destroy and recreate the snapshot without a KMS key."),
			},
		},
	})
}

func testAccComputeSnapshot_updateKmsKey_step1(snapshotName, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
}
`, diskName, snapshotName)
}

func testAccComputeSnapshot_updateKmsKey_step2(snapshotName, diskName, kmsKey string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    kms_key_self_link = "%s"
  }
}
`, diskName, snapshotName, kmsKey)
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

func TestAccComputeSnapshot_snapshotType(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	context1 := map[string]interface{}{
		"random_suffix": randomSuffix,
		"snapshot_type": "ARCHIVE",
	}

	context2 := map[string]interface{}{
		"random_suffix": randomSuffix,
		"snapshot_type": "STANDARD",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_snapshotType(context1),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "source_disk", "source_disk_encryption_key", "terraform_labels", "zone"},
			},
			{
				Config: testAccComputeSnapshot_snapshotType(context2),
			},
			{
				ResourceName:            "google_compute_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "snapshot_encryption_key.0.raw_key", "snapshot_encryption_key.0.rsa_encrypted_key", "source_disk", "source_disk_encryption_key", "terraform_labels", "zone"},
			},
		},
	})
}

func testAccComputeSnapshot_snapshotType(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_snapshot" "snapshot" {
  name        = "tf-test-my-snapshot%{random_suffix}"
  source_disk = google_compute_disk.persistent.id
  zone        = "us-central1-a"
  labels = {
    my_label = "value"
  }
  storage_locations = ["us-central1"]
  snapshot_type     = "%{snapshot_type}"
}

data "google_compute_image" "debian" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
  name  = "tf-test-debian-disk%{random_suffix}"
  image = data.google_compute_image.debian.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}
`, context)
}

func TestAccComputeSnapshot_resourceManagerTags(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    pid,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_resourceManagerTags(context),
			},
		},
	})
}

func testAccComputeSnapshot_resourceManagerTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "tag_key" {
  parent     = "projects/%{project_id}"
  short_name = "tf-test-key-%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value" {
  parent     = "tagKeys/${google_tags_tag_key.tag_key.name}"
  short_name = "tf-test-value-%{random_suffix}"
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-disk-%{random_suffix}"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instant_snapshot" "foobar" {
  name        = "tf-test-instant-snapshot-%{random_suffix}"
  zone        = "us-central1-a"
  source_disk = google_compute_disk.foobar.id
}

resource "google_compute_snapshot" "foobar" {
  name                    = "tf-test-snapshot-%{random_suffix}"
  zone                    = "us-central1-a"
  source_instant_snapshot = google_compute_instant_snapshot.foobar.id
  params {
    resource_manager_tags = {
      "${google_tags_tag_key.tag_key.id}" = "${google_tags_tag_value.tag_value.id}"
    }
  }
}
`, context)
}
