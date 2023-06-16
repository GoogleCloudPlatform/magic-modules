package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageBucketAccessControl_update(t *testing.T) {
	t.Parallel()

<<<<<<< HEAD:mmv1/third_party/terraform/tests/resource_storage_bucket_access_control_test.go
	bucketName := testBucketName(t)
<<<<<<< HEAD:mmv1/third_party/terraform/tests/resource_storage_bucket_access_control_test.go
	acctest.VcrTest(t, resource.TestCase{
=======
	VcrTest(t, resource.TestCase{
=======
	bucketName := acctest.TestBucketName(t)
	acctest.VcrTest(t, resource.TestCase{
>>>>>>> 12945f953 (Generate Mmv1 test files to the service packages):mmv1/third_party/terraform/services/storage/resource_storage_bucket_access_control_test.go
>>>>>>> c13a90bef (Generate Mmv1 test files to the service packages):mmv1/third_party/terraform/services/storage/resource_storage_bucket_access_control_test.go
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageObjectAccessControlDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketAccessControlBasic(bucketName, "READER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_bucket_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleStorageBucketAccessControlBasic(bucketName, "OWNER", "allUsers"),
			},
			{
				ResourceName:      "google_storage_bucket_access_control.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleStorageBucketAccessControlBasic(bucketName, role, entity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket_access_control" "default" {
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  entity = "%s"
}

resource "google_storage_bucket" "bucket" {
	name     = "%s"
	location = "US"
}
`, role, entity, bucketName)
}
