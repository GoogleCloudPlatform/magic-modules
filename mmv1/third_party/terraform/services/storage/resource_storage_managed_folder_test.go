package storage_test

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageManagedFolder_storageManagedFolderUpdate(t *testing.T) {
	t.Parallel()
	bucketName := fmt.Sprintf("tf-test-managed-folder-%s", acctest.RandString(t, 10))
	folderName := "managed/folder/name/"
	objectName := folderName + "file.txt"
	content := "This file will affect the folder being deleted if allowNonEmpty=false"
	h := md5.New()
	if _, err := h.Write([]byte(content)); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	contentMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(folderName, false),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(folderName, true),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(folderName, true) + testAccStorageManagedFolder_object(objectName, content),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_object(objectName, content),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, contentMd5),
			},
		},
	})
}

func testAccStorageManagedFolder_bucket(name string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
}
`, name)
}

func testAccStorageManagedFolder_managedFolder(folderName string, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "%s"
  force_destroy = %t
}
`, folderName, forceDestroy)
}

func testAccStorageManagedFolder_object(objectName string, content string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket_object" "object" {
  name       = "%s"
  content    = "%s"
  bucket     = google_storage_bucket.bucket.name
}`, objectName, content)
}
