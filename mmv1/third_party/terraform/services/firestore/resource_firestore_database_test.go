package firestore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirestoreDatabase_tags(t *testing.T) {
	t.Parallel()

	// Bootstrap shared tag key and value
	tagKey := acctest.BootstrapSharedTestProjectTagKey(t, "firestore-databases-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"pid":           envvar.GetTestProjectFromEnv(),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestProjectTagValue(t, "firestore-databases-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirestoreDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirestoreDatabaseTags(context),
			},
			{
				ResourceName:            "google_firestore_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "etag", "deletion_policy", "tags"},
			},
		},
	})
}

func testAccFirestoreDatabaseTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_firestore_database" "database" {
      name                              = "tf-test-database-%{random_suffix}"
      location_id                       = "nam5"
      type                              = "FIRESTORE_NATIVE"
      delete_protection_state           = "DELETE_PROTECTION_DISABLED"
      deletion_policy                   = "DELETE"
      tags = {
        "%{pid}/%{tagKey}" = "%{tagValue}"
      }
    }
  `, context)
}
