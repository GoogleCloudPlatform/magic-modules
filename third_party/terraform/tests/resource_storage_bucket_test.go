package google

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func TestAccStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "false"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportStateId:     fmt.Sprintf("%s/%s", getTestProjectFromEnv(), bucketName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_requesterPays(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "requester_pays", "true"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_lowercaseLocation(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lowercaseLocation(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_customAttributes(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_lifecycleRules(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRules(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_storageClass(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_storageClass(bucketName, "MULTI_REGIONAL", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "NEARLINE", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					// storage_class-only change should not recreate
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "REGIONAL", "US-CENTRAL1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					// Location change causes recreate
					testAccCheckStorageBucketWasRecreated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_update_requesterPays(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_requesterPays(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_update(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var recreated storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "false"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &recreated),
					testAccCheckStorageBucketWasRecreated(&recreated, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes_withLifecycle1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes_withLifecycle2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_forceDestroy(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(acctest.RandomWithPrefix("tf-test-acl-bucket")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketMissing(bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_forceDestroyWithVersioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_versioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_versioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_logging(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_logging(bucketName, "log-bucket"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_bucket", "log-bucket"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_object_prefix", bucketName),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_loggingWithPrefix(bucketName, "another-log-bucket", "object-prefix"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_bucket", "another-log-bucket"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_object_prefix", "object-prefix"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "0"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_cors(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsCors(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_encryption(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":    getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(10),
		"random_int":      acctest.RandInt(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_encryption(context),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_labels(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			// Going from two labels
			{
				Config: testAccStorageBucket_updateLabels(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Down to only one label (test single label deletion)
			{
				Config: testAccStorageBucket_labels(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// And make sure deleting all labels work
			{
				Config: testAccStorageBucket_basic(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckStorageBucketExists(n string, bucketName string, bucket *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project_ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientStorage.Buckets.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Bucket not found")
		}

		if found.Name != bucketName {
			return fmt.Errorf("expected name %s, got %s", bucketName, found.Name)
		}

		*bucket = *found
		return nil
	}
}

func testAccCheckStorageBucketWasUpdated(newBucket *storage.Bucket, b *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if newBucket.TimeCreated != b.TimeCreated {
			return fmt.Errorf("expected storage bucket to have been updated (had same creation time), instead was recreated - old creation time %s, new creation time %s", newBucket.TimeCreated, b.TimeCreated)
		}
		return nil
	}
}

func testAccCheckStorageBucketWasRecreated(newBucket *storage.Bucket, b *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if newBucket.TimeCreated == b.TimeCreated {
			return fmt.Errorf("expected storage bucket to have been recreated, instead had same creation time (%s)", b.TimeCreated)
		}
		return nil
	}
}

func testAccCheckStorageBucketPutItem(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		// This needs to use Media(io.Reader) call, otherwise it does not go to /upload API and fails
		if res, err := config.clientStorage.Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		return nil
	}
}

func testAccCheckStorageBucketMissing(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientStorage.Buckets.Get(bucketName).Do()
		if err == nil {
			return fmt.Errorf("Found %s", bucketName)
		}

		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return nil
		}

		return err
	}
}

func testAccStorageBucketDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket" {
			continue
		}

		_, err := config.clientStorage.Buckets.Get(rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Bucket still exists")
		}
	}

	return nil
}

func testAccStorageBucket_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}
`, bucketName)
}

func testAccStorageBucket_requesterPays(bucketName string, pays bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	requester_pays = %t
}
`, bucketName, pays)
}

func testAccStorageBucket_lowercaseLocation(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	location = "eu"
}
`, bucketName)
}

func testAccStorageBucket_customAttributes(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	location = "EU"
	force_destroy = "true"
}
`, bucketName)
}

func testAccStorageBucket_customAttributes_withLifecycle1(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	location = "EU"
	force_destroy = "true"
	lifecycle_rule {
		action {
			type = "Delete"
		}
		condition {
			age = 10
		}
	}
}
`, bucketName)
}

func testAccStorageBucket_customAttributes_withLifecycle2(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	location = "EU"
	force_destroy = "true"
	lifecycle_rule {
		action {
			type = "SetStorageClass"
			storage_class = "NEARLINE"
		}
		condition {
			age = 2
		}
	}
	lifecycle_rule {
		action {
			type = "Delete"
		}
		condition {
			age = 10
			num_newer_versions = 2
		}
	}
}
`, bucketName)
}

func testAccStorageBucket_storageClass(bucketName, storageClass, location string) string {
	var locationBlock string
	if location != "" {
		locationBlock = fmt.Sprintf(`
	location = "%s"`, location)
	}
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	storage_class = "%s"%s
}
`, bucketName, storageClass, locationBlock)
}

func testGoogleStorageBucketsCors(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	cors {
	  origin = ["abc", "def"]
	  method = ["a1a"]
	  response_header = ["123", "456", "789"]
	  max_age_seconds = 10
	}

	cors {
	  origin = ["ghi", "jkl"]
	  method = ["z9z"]
	  response_header = ["000"]
	  max_age_seconds = 5
	}
}
`, bucketName)
}

func testAccStorageBucket_forceDestroyWithVersioning(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	force_destroy = "true"
	versioning {
	  enabled = "true"
	}
}
`, bucketName)
}

func testAccStorageBucket_versioning(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	versioning {
	  enabled = "true"
	}
}
`, bucketName)
}

func testAccStorageBucket_logging(bucketName string, logBucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	logging {
		log_bucket = "%s"
	}
}
`, bucketName, logBucketName)
}

func testAccStorageBucket_loggingWithPrefix(bucketName string, logBucketName string, prefix string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	logging {
		log_bucket = "%s"
		log_object_prefix = "%s"
	}
}
`, bucketName, logBucketName, prefix)
}

func testAccStorageBucket_lifecycleRules(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	lifecycle_rule {
		action {
			type = "SetStorageClass"
			storage_class = "NEARLINE"
		}
		condition {
			age = 2
		}
  	}
	lifecycle_rule {
		action {
			type = "Delete"
		}
		condition {
			age = 10
		}
	}
}
`, bucketName)
}

func testAccStorageBucket_labels(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	labels = {
		my-label = "my-label-value"
	}
}
`, bucketName)
}

func testAccStorageBucket_encryption(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "acceptance" {
	name            = "terraform-%{random_suffix}"
	project_id      = "terraform-%{random_suffix}"
	org_id          = "%{organization}"
	billing_account = "%{billing_account}"
}

resource "google_project_services" "acceptance" {
	project = "%{google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	name     = "tf-test-%{random_suffix}"
	project  = "${google_project_services.acceptance.project}"
	location = "us"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "tf-test-%{random_suffix}"
	key_ring        = "${google_kms_key_ring.key_ring.id}"
	rotation_period = "1000000s"
}

resource "google_storage_bucket" "bucket" {
	name = "tf-test-crypto-bucket-%{random_int}"
	encryption {
		default_kms_key_name = "${google_kms_crypto_key.crypto_key.self_link}"
	}
}
	`, context)
}

func testAccStorageBucket_updateLabels(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
	labels = {
		my-label    = "my-updated-label-value"
		a-new-label = "a-new-label-value"
	}
}
`, bucketName)
}
