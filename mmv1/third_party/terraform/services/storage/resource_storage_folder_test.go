// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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
			},
		},
	})
}

func TestAccStorageFolder_DeleteSingleFolderDisableForceDestroy(t *testing.T) {
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
