package kms_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsAutokeyConfig_basic(t *testing.T) {
	kmsAutokey := acctest.BootstrapKMSAutokeyKeyHandle(t)
	folder := fmt.Sprintf("folders/%s", strings.Split(kmsAutokey.AutokeyConfig.Name, "/")[1])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsAutokeyConfig_basic(folder),
				Check:  resource.TestMatchResourceAttr("data.google_kms_autokey_config.kms_autokey_config", "id", regexp.MustCompile(kmsAutokey.AutokeyConfig.Name)),
			},
		},
	})
}

func testAccDataSourceGoogleKmsAutokeyConfig_basic(folder string) string {

	return fmt.Sprintf(`
data "google_kms_autokey_config" "kms_autokey_config" {
  folder = "%s"
}
`, folder)
}