package certificatemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerIssuanceConfig_tags(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestTagKey(t, "certificate-manager-issuance-config-tagkey")
	tagValue := acctest.BootstrapSharedTestTagValue(t, "certificate-manager-issuance-config-tagvalue", tagKey)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerIssuanceConfigTags(context, map[string]string{org + "/" + tagKey: tagValue}),
			},
			{
				ResourceName:            "google_certificate_manager_certificate_issuance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCertificateManagerIssuanceConfigTags(context map[string]interface{}, tags map[string]string) string {
	r := acctest.Nprintf(`
resource "google_certificate_manager_certificate_issuance_config" "default" {
        name        = "tf-test-issuance-config%{random_suffix}"
        description = "sample description for the issaunce config"
        location    = "us-central1"
tags = {`, context)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}
