package dataprocmetastore_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
	"regexp"
        "fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetastoreFederation_deletionprotection(t *testing.T) {
	t.Parallel()
	
	name := "tf-test-metastore-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreFederationDeletionProtection(name, "us-central1"),
			},
			{
				ResourceName:            "google_dataproc_metastore_federation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
                                Config: testAccMetastoreFederationDeletionProtection(name, "us-west2"),
                                ExpectError: regexp.MustCompile("deletion_protection"),
                        },
			{
				Config: testAccMetastoreFederationDeletionProtectionFalse(name, "us-central1"),
			},
			{
				ResourceName:            "google_dataproc_metastore_federation.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccMetastoreFederationDeletionProtection(name string, location string) string {

	return fmt.Sprintf(`
       resource "google_dataproc_metastore_service" "default" {
         service_id = "%s"
         location   = "us-central1"
         tier       = "DEVELOPER"
         hive_metastore_config {
           version           = "3.1.2"
           endpoint_protocol = "GRPC"
         }
         }
       resource "google_dataproc_metastore_federation" "default" {
          federation_id = "%s"
          location      = "%s"
          version       = "3.1.2"
          deletion_protection = true
          backend_metastores {
            rank           = "1"
            name           = google_dataproc_metastore_service.default.id
            metastore_type = "DATAPROC_METASTORE" 
         }
}
`,name, name, location)
}

func testAccMetastoreFederationDeletionProtectionFalse(name string, location string) string {

	return fmt.Sprintf(`
       resource "google_dataproc_metastore_service" "default" {
         service_id = "%s"
         location   = "us-central1"
         tier       = "DEVELOPER"
         hive_metastore_config {
           version           = "3.1.2"
           endpoint_protocol = "GRPC"
         }
         }
       resource "google_dataproc_metastore_federation" "default" {
          federation_id = "%s"
          location      = "%s"
          version       = "3.1.2"
          deletion_protection = false
          backend_metastores {
            rank           = "1"
            name           = google_dataproc_metastore_service.default.id
            metastore_type = "DATAPROC_METASTORE" 
         }
}
`,name, name, location)
}
	  
