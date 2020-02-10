package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccStorageHmacKey_update(t *testing.T) {
	t.Parallel()

	saName := fmt.Sprintf("%v%v", "service-account", acctest.RandString(10))
	bucketName := testBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageHmacKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleStorageHmacKeyBasic(saName, bucketName, "ACTIVE"),
			},
			{
				ResourceName:      "google_storage_hmac_key.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGoogleStorageHmacKeyBasic(saName, bucketName, "INACTIVE"),
			},
			{
				ResourceName:      "google_storage_hmac_key.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStorageHmacKeyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_hmac_key" {
			continue
		}
		accessId := rs.Primary.Attributes["accessId"]

		_, err := config.clientStorage.HmacKeys.Get(accessId).Do()
		if err == nil {
			return fmt.Errorf("Hmac key still exists.")
		}
	}

	return nil
}

func testAccGoogleStorageHmacKeyBasic(saName, bucketName, state string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id = "%s"
}

resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_hmac_key" "key" {
	service_account_email = google_service_account.service_account.email
	state = "%s"
}
`, saName, bucketName, state)
}
