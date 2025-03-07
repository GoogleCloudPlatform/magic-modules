// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagebatchoperations_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageBatchOperationsJobs_storageBatchOperationsError(t *testing.T) {
	t.Parallel()
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageBatchOperationsJobs_storageBatchOperationsError(jobID),
				ExpectError: regexp.MustCompile("but `delete_object,put_metadata` were specified"),
			},
			{
				Config:      testAccStorageBatchOperationsJobs_storageBatchOperationsJobIDError(),
				ExpectError: regexp.MustCompile("doesn't match regexp"),
			},
			{
				Config:      testAccStorageBatchOperationsJobs_storageBatchOperationsManifestError(jobID),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_createJobWithPrefix(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsCreateJobWithPrefix(bucketName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_jobWithPrefixDeleteObjectAllVersions(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsJobWithPrefixDeleteObjectAllVersions(bucketName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_jobWithPrefixDeleteLiveObject(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsJobWithPrefixDeleteObject(bucketName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_jobWithPrefixObjectHold(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsJobWithPrefixObjectHold(bucketName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_createJobWithManifest(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsCreateJobWithManifest(bucketName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func TestAccStorageBatchOperationsJobs_batchOperationJobKmsKey(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	jobID := fmt.Sprintf("tf-test-job-%d", acctest.RandInt(t))
	keyRing := fmt.Sprintf("tf-test-keyring-%d", acctest.RandInt(t))
	cryptoKey := fmt.Sprintf("tf-test-cryptokey-%d", acctest.RandInt(t))
	objectName := fmt.Sprintf("tf-test-object-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBatchOperationsJobs_storageBatchOerationsJobKmsKey(keyRing, cryptoKey, bucketName, objectName, jobID),
			},
			{
				ResourceName:            "google_storage_batch_operations_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_id", "location", "delete_protection"},
			},
		},
	})
}

func testAccStorageBatchOperationsCreateJobWithPrefix(bucketName, jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global"
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			prefix_list {
				included_object_prefixes = [
					"bkt"
				]
			}
		}
	}

	put_metadata {
		custom_metadata = {
			"key"="value"
			"key1"="value1"
		}
	}
	delete_protection = false
}
`, bucketName, jobID)
}

func testAccStorageBatchOperationsJobWithPrefixDeleteObjectAllVersions(bucketName, jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global"
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			prefix_list {
				included_object_prefixes = [
					"bkt"
				]
			}
		}
	}
	delete_object {
		permanent_object_deletion_enabled = true
	}
	delete_protection = false
}
`, bucketName, jobID)
}

func testAccStorageBatchOperationsJobWithPrefixDeleteObject(bucketName, jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global"
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			prefix_list {
				included_object_prefixes = [
					"objprefix"
				]
			}
		}
	}
	delete_object {
		permanent_object_deletion_enabled = false
	}
	delete_protection = false
}
`, bucketName, jobID)
}

func testAccStorageBatchOperationsJobWithPrefixObjectHold(bucketName, jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global"
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			prefix_list {
				included_object_prefixes = [
					"objprefix", "prefix2"
				]
			}
		}
	}
	put_object_hold {
		event_based_hold= "SET"
	}

	delete_protection = false
}
`, bucketName, jobID)
}

func testAccStorageBatchOperationsCreateJobWithManifest(bucketName, jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global" 
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			manifest {
				manifest_location = "gs://%s/manifest.csv"
			}
		}
	}
	put_metadata {
		custom_metadata = {
			"key"="value"
			"key1"="value1"
		}
		cache_control = "public, max-age=3600"
		content_disposition = "sample"
		content_encoding = "text"
		content_language = "en-us"
		content_type = "application/json"
	}
	delete_protection = false
}
`, bucketName, jobID, bucketName)
}

func testAccStorageBatchOperationsJobs_storageBatchOperationsManifestError(jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "us-central1"
	bucket_list {
		buckets  {
			bucket = "test-bkt"
			prefix_list {
				included_object_prefixes = [
					"bkt"
				]
			}
			manifest {
				manifest_location = "gs://bucket/file.csv"
			}
		}
	}
	delete_object  {
		permanent_object_deletion_enabled = false
	}
}
`, jobID)
}

func testAccStorageBatchOperationsJobs_storageBatchOperationsError(jobID string) string {
	return fmt.Sprintf(`
resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "us-central1"
	bucket_list {
		buckets  {
			bucket = "test-bkt"
			manifest {
				manifest_location = "gs://bucket/file.csv"
			}
		}
	}
	delete_object  {
		permanent_object_deletion_enabled = false
	}
	put_metadata {
		content_type = "application/json"
	}
}
`, jobID)
}

func testAccStorageBatchOperationsJobs_storageBatchOperationsJobIDError() string {
	return fmt.Sprintf(`
resource "google_storage_batch_operations_job" "job" {
	job_id     = "tf-job@d-"
	location = "global"
	bucket_list {
		buckets  {
			bucket = "test-bkt"
			manifest {
				manifest_location = "gs://bucket/file.csv"
			}
		}
	}
	put_metadata {
		custom_metadata = {
			"key"="value"
			"key1"="value1"
		}
	}
}
`)
}

func testAccStorageBatchOperationsJobs_storageBatchOerationsJobKmsKey(kmsKeyRing, kmsKeyName, bucketName, objectName, jobID string) string {
	return fmt.Sprintf(`
resource "google_kms_key_ring" "tf_keyring" {
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "tf_crypto_key" {
  name            = "%s"
  key_ring        = google_kms_key_ring.tf_keyring.id
  rotation_period = "7776000s"

  lifecycle {
    prevent_destroy = false
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "us-central1"
  uniform_bucket_level_access = true
  force_destroy = true
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  content       = "test-content"
}

data "google_storage_project_service_account" "gcs_account" {
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.tf_crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
}

resource "google_storage_batch_operations_job" "job" {
	job_id     = "%s"
	location = "global"
	bucket_list {
		buckets  {
			bucket = google_storage_bucket.bucket.name
			prefix_list {
				included_object_prefixes = [
					"objprefix"
				]
			}
		}
	}
	rewrite_object {
		kms_key = google_kms_crypto_key.tf_crypto_key.id
	}

	delete_protection = false
}
`, kmsKeyRing, kmsKeyName, bucketName, objectName, jobID)
}
