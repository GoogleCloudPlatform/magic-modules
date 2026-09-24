package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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

func TestAccComputeSnapshot_encryptionCMEKUpdate(t *testing.T) {
	t.Parallel()

	key1 := kms.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name
	key2 := kms.BootstrapKMSKey(t).CryptoKey.Name
	suffix := acctest.RandString(t, 10)
	snapshotName := fmt.Sprintf("tf-test-%s", suffix)
	diskName := fmt.Sprintf("tf-test-%s", suffix)
	var creationTimestamp string

	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@compute-system.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, "", ""),
			},
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, key1, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_snapshot.foobar", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotKmsKey(t, "google_compute_snapshot.foobar", key1, ""),
					resource.TestCheckResourceAttrWith("google_compute_snapshot.foobar", "creation_timestamp", func(v string) error {
						creationTimestamp = v
						return nil
					}),
				),
			},
			{
				ResourceName:      "google_compute_snapshot.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// Not returned by the API.
				ImportStateVerifyIgnore: []string{"zone", "source_disk"},
			},
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, key2, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_snapshot.foobar", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotKmsKey(t, "google_compute_snapshot.foobar", key2, ""),
					resource.TestCheckResourceAttr("google_compute_snapshot.foobar", "snapshot_encryption_key.0.kms_key_self_link", key2),
					resource.TestCheckResourceAttrPtr("google_compute_snapshot.foobar", "creation_timestamp", &creationTimestamp),
				),
			},
			{
				ResourceName:      "google_compute_snapshot.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// Not returned by the API.
				ImportStateVerifyIgnore: []string{"zone", "source_disk"},
			},
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, "", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_snapshot.foobar", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotKmsKey(t, "google_compute_snapshot.foobar", "", ""),
				),
			},
		},
	})
}

// updateKmsKey clears kmsKeyServiceAccount, so this recreates the snapshot.
func TestAccComputeSnapshot_encryptionCMEKUpdateWithServiceAccount(t *testing.T) {
	t.Parallel()

	key1 := kms.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name
	key2 := kms.BootstrapKMSKey(t).CryptoKey.Name
	suffix := acctest.RandString(t, 10)
	snapshotName := fmt.Sprintf("tf-test-%s", suffix)
	diskName := fmt.Sprintf("tf-test-%s", suffix)

	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@compute-system.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:{project_number}-compute@developer.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})
	serviceAccount := fmt.Sprintf("%s-compute@developer.gserviceaccount.com", envvar.GetTestProjectNumberFromEnv())

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, key1, serviceAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotKmsKey(t, "google_compute_snapshot.foobar", key1, serviceAccount),
				),
			},
			{
				ResourceName:      "google_compute_snapshot.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// Not returned by the API.
				ImportStateVerifyIgnore: []string{"zone", "source_disk"},
			},
			{
				Config: testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, key2, serviceAccount),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_snapshot.foobar", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotKmsKey(t, "google_compute_snapshot.foobar", key2, serviceAccount),
				),
			},
		},
	})
}

func testAccCheckComputeSnapshotKmsKey(t *testing.T, resourceName, wantKey, wantServiceAccount string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		config := acctest.GoogleProviderConfig(t)
		project := envvar.GetTestProjectFromEnv()
		url := fmt.Sprintf("%sprojects/%s/global/snapshots/%s", transport_tpg.BaseUrl(tpgcompute.Product, config), project, rs.Primary.Attributes["name"])
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}
		var gotKey, gotServiceAccount string
		if encryptionKey, ok := res["snapshotEncryptionKey"].(map[string]interface{}); ok {
			gotKey, _ = encryptionKey["kmsKeyName"].(string)
			gotServiceAccount, _ = encryptionKey["kmsKeyServiceAccount"].(string)
		}
		gotKey = strings.Split(gotKey, "/cryptoKeyVersions/")[0]
		if gotKey != wantKey {
			return fmt.Errorf("snapshot has Cloud KMS key %q, want %q", gotKey, wantKey)
		}
		if gotServiceAccount != wantServiceAccount {
			return fmt.Errorf("snapshot has kmsKeyServiceAccount %q, want %q", gotServiceAccount, wantServiceAccount)
		}
		return nil
	}
}

func testAccComputeSnapshot_encryptionCMEKUpdate(snapshotName, diskName, kmsKey, serviceAccount string) string {
	encryptionKey := ""
	if kmsKey != "" {
		serviceAccountLine := ""
		if serviceAccount != "" {
			serviceAccountLine = fmt.Sprintf("\n    kms_key_service_account = %q", serviceAccount)
		}
		encryptionKey = fmt.Sprintf(`
  snapshot_encryption_key {
    kms_key_self_link = %q%s
  }`, kmsKey, serviceAccountLine)
	}
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-balanced"
  zone = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
%s
}
`, diskName, snapshotName, encryptionKey)
}

func testAccComputeSnapshot_encryption(snapshotName string, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-13"
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
  family  = "debian-13"
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
  family  = "debian-13"
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
