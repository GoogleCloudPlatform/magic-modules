package dataprocmetastore_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccMetastoreFederation_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestTagKey(t, "dataproc_metastore_federation-tagkey")
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestTagValue(t, "dataproc_metastore_federation-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreFederationTags(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_federation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccMetastoreFederationTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_dataproc_metastore_service" "default" {
			service_id = "tf-test-service-%{random_suffix}"
			location   = "us-central1"
			tier       = "DEVELOPER"

			hive_metastore_config {
				version           = "3.1.2"
				endpoint_protocol = "GRPC"
			}
		}

		resource "google_dataproc_metastore_federation" "default" {
			location      = "us-central1"
			federation_id = "tf-test-federation-%{random_suffix}"
			version       = "3.1.2"

			backend_metastores {
				rank           = "1"
				name           = google_dataproc_metastore_service.default.id
				metastore_type = "DATAPROC_METASTORE"
			}

			tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
		}
	`, context)

}
