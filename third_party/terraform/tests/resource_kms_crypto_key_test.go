package google

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/iterator"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func TestCryptoKeyIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedCryptoKeyId string
		Config              *Config
	}{
		"id is in project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id is in domain:project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/example.com:test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/test-key-ring/can-you-believe-that-this-cryptokey-name-is-this-extravagantly-long",
			ExpectedError: true,
		},
		"id is in location/keyRingName/cryptoKeyName format": {
			ImportId:            "us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/keyRingName/cryptoKeyName format without project in config": {
			ImportId:      "us-central1/test-key-ring/test-key-name",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		cryptoKeyId, err := parseKmsCryptoKeyId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if cryptoKeyId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, cryptoKeyId.terraformId())
		}

		if cryptoKeyId.cryptoKeyId() != tc.ExpectedCryptoKeyId {
			t.Fatalf("bad: %s, expected CryptoKey ID to be `%s` but is `%s`", tn, tc.ExpectedCryptoKeyId, cryptoKeyId.cryptoKeyId())
		}
	}
}

func TestCryptoKeyNextRotationCalculation(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	rotationPeriod, timestamp, err := kmsCryptoKeyNextRotation("1000000s", now)
	if err != nil {
		t.Fatal(err)
	}

	if act, exp := rotationPeriod.RotationPeriod.Seconds, int64(1000000); act != exp {
		t.Errorf("expected %d to be %d", act, exp)
	}

	if act, exp := timestamp.Seconds, now.Add(1000000*time.Second).Unix(); act != exp {
		t.Errorf("expected %d to be %d", act, exp)
	}
}

func TestAccKmsCryptoKey_basic(t *testing.T) {
	t.Parallel()

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_rotation(t *testing.T) {
	t.Parallel()

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	rotationPeriod := "100000s"
	updatedRotationPeriod := "7776000s"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, updatedRotationPeriod),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKey_rotationRemoved(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

/*
	KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource was removed from state,
	even though the server-side resource was not removed.
*/
func testAccCheckGoogleKmsCryptoKeyWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

/*
	KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource's CryptoKeyVersion
	sub-resources were scheduled to be destroyed, rendering the key itself inoperable.
*/

func testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		gcpResourceUri := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyRingName, cryptoKeyName)

		ctx := context.Background()
		it := config.clientKms.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
			Parent: gcpResourceUri,
		})

		for {
			ckv, err := it.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
				return fmt.Errorf("Failed to list crypto key versions: %s", err)
			}

			if ckv.State != kmspb.CryptoKeyVersion_DESTROYED &&
				ckv.State != kmspb.CryptoKeyVersion_DESTROY_SCHEDULED {
				return fmt.Errorf("CryptoKey %s should have no versions, but version %s has state %s",
					cryptoKeyName, ckv.Name, ckv.State)
			}
		}

		return nil
	}
}

/*
	This test runs in its own project, otherwise the test project would start to get filled
	with undeletable resources
*/
func testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.self_link}"
	rotation_period = "1000000s"
	version_template {
		algorithm =        "symmetric_encryption"
		protection_level = "software"
	}
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.self_link}"
	rotation_period = "%s"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod)
}

func testGoogleKmsCryptoKey_rotationRemoved(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.self_link}"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}
