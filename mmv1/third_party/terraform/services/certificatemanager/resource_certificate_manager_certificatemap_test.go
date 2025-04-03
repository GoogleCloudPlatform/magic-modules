// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerCertificateMap_tags(t *testing.T) {
	t.Parallel()
	org := envvar.GetTestOrgFromEnv(t)
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	tagKey := acctest.BootstrapSharedTestTagKey(t, "ccm-certificatemaps-tagkey")
	tagValue := acctest.BootstrapSharedTestTagValue(t, "ccm-certificatemaps-tagvalue", tagKey)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateMapTags(name, map[string]string{org + "/" + tagKey: tagValue}),
			},
			{
				ResourceName:            "google_certificate_manager_certificate_map.certificatemap",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateMapTags(name string, tags map[string]string) string {
	r := fmt.Sprintf(`
resource "google_certificate_manager_certificate_map" "certificatemap" {
  name = "tf-certificate-map-%s"
  description = "Global cert"
tags = {`, name)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}
