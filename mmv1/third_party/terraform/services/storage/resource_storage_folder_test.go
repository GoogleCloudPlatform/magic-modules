// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"google.golang.org/api/storage/v1"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageFolder_storageFolderBasic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageFolder(true),
			},
			{
				ResourceName:            "google_storage_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "recursive", "force_destroy"},
			},
		},
	})
}

func TestAccStorageFolder_hnsDisabled(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageFolder_storageBucket(bucketName, false, true) + testAccStorageFolder_storageFolder(true),
				ExpectError: regexp.MustCompile("Error creating Folder: googleapi: Error 409: The bucket does not support hierarchical namespace., conflict"),
			},
		},
	})
}

func TestAccStorageFolder_FolderForceDestroy(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	data := []byte("data data data")

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucketObject(bucketName, true, true, testFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketUploadItem(t, bucketName),
				),
			},
		},
	})
}

func TestAccStorageFolder_DeleteEmptyFolderWithForceDestroyDefault(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageOneFolder(false),
			},
		},
	})
}

func TestAccStorageFolder_FailDeleteNonEmptyFolder(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	folderName := "folder/"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageOneFolder(false),
				Check: resource.ComposeTestCheckFunc(
					testAccStorageCreatSubFolder(t, bucketName, folderName),
					testAccStorageDeleteFolder(t, bucketName, folderName),
				),
				ExpectError: regexp.MustCompile("googleapi: Error 409: The folder you tried to delete is not empty"),
			},
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageOneFolder(true),
			},
		},
	})
}

func testAccStorageFolder_storageBucket(bucketName string, hnsFlag bool, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = %t
  }
  force_destroy = %t
}
`, bucketName, hnsFlag, forceDestroy)
}

func testAccStorageFolder_storageBucketObject(bucketName string, hnsFlag bool, forceDestroy bool, fileName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = %t
  }
  force_destroy = true
}
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "${google_storage_folder.folder.name}subfolder/"
  force_destroy = %t
}  
resource "google_storage_bucket_object" "object" {
  name   = "${google_storage_folder.subfolder.name}tffile"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}  
`, bucketName, hnsFlag, forceDestroy, forceDestroy, fileName)
}

func testAccStorageFolder_storageFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "${google_storage_folder.folder.name}name/"
  force_destroy = %t
}  
`, forceDestroy, forceDestroy)
}

func testAccStorageFolder_storageOneFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
} 
`, forceDestroy)
}

func testAccStorageCreatSubFolder(t *testing.T, bucketName, parentFolder string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		subFolder := &storage.Folder{
			Name: parentFolder + "subfolder/",
		}
		if res, err := config.NewStorageClient(config.UserAgent).Folders.Insert(bucketName, subFolder).Do(); err == nil {
			log.Printf("sub folder created: %s", res.Name)
		} else {
			log.Printf("failed to create sub folder: %s", subFolder.Name)
		}
		return nil
	}
}

func testAccStorageDeleteFolder(t *testing.T, bucketName, parentFolder string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		var deleteError error
		if err := config.NewStorageClient(config.UserAgent).Folders.Delete(bucketName, parentFolder).Do(); err == nil {
			log.Printf("successfully deleted folder: %s", err)
		} else {
			deleteError = fmt.Errorf("failed to deleted folder: %s", err)
		}
		return deleteError
	}
}

func testAccCheckStorageBucketUploadItem(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		if res, err := config.NewStorageClient(config.UserAgent).Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		return nil
	}
}
