package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFirestoreDocument_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFirestoreDocument_update(name),
			},
			resource.TestStep{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccFirestoreDocument_update2(name),
			},
			resource.TestStep{
				ResourceName:      "google_firestore_document.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirestoreDocument_update(name string) string {
	return fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	database   = "(default)"
	collection = "somenewcollection"
	document_id = "%s"
	fields     = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"val1\"}}}}}"
}
`, name)
}

func testAccFirestoreDocument_update2(name string) string {
	return fmt.Sprintf(`
resource "google_firestore_document" "instance" {
	database   = "(default)"
	collection = "somenewcollection"
	document_id = "%s"
	fields     = "{\"something\":{\"mapValue\":{\"fields\":{\"yo\":{\"stringValue\":\"val2\"}}}}}"
}
`, name)
}
