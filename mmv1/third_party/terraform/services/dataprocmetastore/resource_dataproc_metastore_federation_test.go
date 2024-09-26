package dataprocmetastore_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccMetastoreFederation_tags(t *testing.T) {
	t.Parallel()

        org := envvar.GetTestOrgFromEnv(t)
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	tagKey := acctest.BootstrapSharedTestTagKey(t, "metastore-federations-tagkey")
	tagValue := acctest.BootstrapSharedTestTagValue(t, "metastore-federations-tagvalue", tagKey)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreFederationTags(context, map[string]string{org + "/" + tagKey: tagValue}),
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

func testAccMetastoreFederationTags(context map[string]interface{}, tags map[string]string) string {

	r := acctest.Nprintf(`

       resource "google_dataproc_metastore_service" "default" {
         service_id = "metastore-srv-sep16"
         location   = "us-central1"
         tier       = "DEVELOPER"


         hive_metastore_config {
           version           = "3.1.2"
           endpoint_protocol = "GRPC"
         }
         }
       resource "google_dataproc_metastore_federation" "default" {
          location      = "us-central1"
          federation_id = "metastore-fed-sep16"
          version       = "3.1.2"

          backend_metastores {
            rank           = "1"
            name           = google_dataproc_metastore_service.default.id
            metastore_type = "DATAPROC_METASTORE" 
         }
	  tags = {`, context)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}
