package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeBackendBucketSignedUrlKey_basic(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	storageName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	keyName := "key_one"
	keyValue := fmt.Sprintf("%s==", acctest.RandString(22))
	keyValueMod := fmt.Sprintf("%s==", acctest.RandString(22))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendBucketSignedUrlKeyDestroy(bucketName),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucketSignedUrlKeys_basic(keyName, keyValue, bucketName, storageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendBucketKeysExists(bucketName, keyName)),
			},
			{
				Config: testAccComputeBackendBucketSignedUrlKeys_basic(keyName, keyValueMod, bucketName, storageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendBucketKeysExists(bucketName, keyName)),
			},
		},
	})
}

func testAccCheckComputeBackendBucketSignedUrlKeyDestroy(backendName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_backend_bucket_signed_url_key" {
				continue
			}

			bucket, err := config.clientCompute.BackendBuckets.Get(config.Project, backendName).Do()
			if err != nil {
				if !isGoogleApiErrorWithCode(err, 404) {
					return fmt.Errorf("error getting backend bucket %s: %v", rs.Primary.ID, err)
				}
				return nil
			}

			if bucket.CdnPolicy != nil {
				for _, keyName := range bucket.CdnPolicy.SignedUrlKeyNames {
					if keyName == rs.Primary.ID {
						return fmt.Errorf("backend bucket signed url key %s still exists", rs.Primary.ID)
					}
				}
			}
		}
		return nil
	}
}

func testAccCheckComputeBackendBucketKeysExists(backendBucketName string, keyNames ...string) resource.TestCheckFunc {
	return func(s *terraform.State) (err error) {
		config := testAccProvider.Meta().(*Config)
		bucket, err := config.clientCompute.BackendBuckets.Get(config.Project, backendBucketName).Do()
		if err != nil {
			return err
		}
		if bucket.CdnPolicy == nil {
			return fmt.Errorf("backend bucket %s not found", backendBucketName)
		}
		if len(bucket.CdnPolicy.SignedUrlKeyNames) != len(keyNames) {
			return fmt.Errorf("unexpected number of keys in backend bucket %s, expected keys: %#v actual keys: (%#v)",
				backendBucketName, keyNames, bucket.CdnPolicy.SignedUrlKeyNames)
		}
		for _, name := range keyNames {
			n := fmt.Sprintf("google_compute_backend_bucket_signed_url_key.%s", name)
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("No ID is set")
			}
			found := false
			for _, keyName := range bucket.CdnPolicy.SignedUrlKeyNames {
				if keyName == rs.Primary.ID {
					found = true
				}
			}
			if !found {
				fmt.Errorf("backend bucket %s signed url key %s not found", backendBucketName, rs.Primary.ID)
			}
		}
		return nil
	}
}

func testAccComputeBackendBucketSignedUrlKeys_basic(keyName, keyValue, backendName, storageName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket_signed_url_key" "%s" {
  name         = "mykey"
  key_value     = "%s"
  backend_name = "${google_compute_backend_bucket.my_backend_bucket.name}"
}

resource "google_compute_backend_bucket" "my_backend_bucket" {
  name        = "%s"
  bucket_name = "${google_storage_bucket.gcs_bucket_one.name}"
  enable_cdn  = true
}

resource "google_storage_bucket" "gcs_bucket_one" {
  name     = "%s"
  location = "EU"
}
`, keyName, keyValue, backendName, storageName)
}
