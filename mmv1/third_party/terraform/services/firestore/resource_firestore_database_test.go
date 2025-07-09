package firestore_test

import (
  "fmt"
  "strings"
  "testing"

  "github.com/hashicorp/terraform-plugin-testing/helper/resource"
  "github.com/hashicorp/terraform-plugin-testing/terraform"

  "github.com/hashicorp/terraform-provider-google/google/acctest"
  "github.com/hashicorp/terraform-provider-google/google/envvar"
  "github.com/hashicorp/terraform-provider-google/google/services/firestore"
)

func TestAccFirestoreDatabase_withSharedTags(t *testing.T) {
  t.Parallel()
  name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
  org := envvar.GetTestOrgFromEnv(t)

  // Bootstrap shared tag key and value
  tagKey := acctest.BootstrapSharedTestTagKey(t, "firestore-db-tagkey")
  tagValue := acctest.BootstrapSharedTestTagValue(t, "firestore-db-tagvalue", tagKey)

  acctest.VcrTest(t, resource.TestCase{
    PreCheck:                 func() { acctest.AccTestPreCheck(t) },
    ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
    CheckDestroy:             testAccCheckFirestoreDatabaseDestroyProducer(t), // Assuming you have this function
    Steps: []resource.TestStep{
      {
        Config: testAccFirestoreDatabaseWithTags(name, tagKey, tagValue),
      },
      {
        ResourceName:            "google_firestore_database.database",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"project", "etag", "deletion_policy"},
      },
    },
  })
}

func testAccFirestoreDatabaseWithTags(name, tagKeyId, tagValueId string) string {
  return fmt.Sprintf(`
    resource "google_firestore_database" "database" {
      name                              = "%s"
      location_id                       = "nam5"
      type                              = "FIRESTORE_NATIVE"
      delete_protection_state           = "DELETE_PROTECTION_ENABLED"
      deletion_policy                   = "DELETE"
      tags = {
        "%s": "%s"
      }
    }
  `, name, tagKeyId, tagValueId)
}

func testAccCheckFirestoreDatabaseDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_firestore_database" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{FirestoreBasePath}}projects/{{project}}/databases/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				Project:              billingProject,
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorAbortPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429QuotaError},
			})
			if err == nil {
				return fmt.Errorf("Firestore Database still exists at %s", url)
			}
		}

		return nil
	}
}