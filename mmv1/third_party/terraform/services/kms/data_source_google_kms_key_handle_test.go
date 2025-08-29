package kms_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsKeyHandle_basic(t *testing.T) {
	kmsAutokey := acctest.BootstrapKMSAutokeyKeyHandle(t)
	keyParts := strings.Split(kmsAutokey.KeyHandle.Name, "/")
	project := keyParts[1]
	location := keyParts[3]
	keyHandleName := keyParts[5]

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyHandle_basic(project, location, keyHandleName),
				Check:  resource.TestMatchResourceAttr("data.google_kms_key_handle.kms_key_handle", "id", regexp.MustCompile(kmsAutokey.KeyHandle.Name)),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyHandle_basic(project string, location string, keyHandleName string) string {

	return fmt.Sprintf(`
data "google_kms_key_handle" "kms_key_handle" {
  name = "%s"
  location = "%s"
  project = "%s"
}
`, keyHandleName, location, project)
}
