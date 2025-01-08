// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccDataSourceGoogleKmsKeyHandles_basic(t *testing.T) {
	kmsAutokey := acctest.BootstrapKMSAutokeyKeyHandle(t)
	keyParts := strings.Split(kmsAutokey.KeyHandle.Name, "/")
	project := keyParts[1]
	location := keyParts[3]
	diskFilter := fmt.Sprintf("resourceTypeSelector=\\\"compute.googleapis.com/Disk\\\"")

	context := map[string]interface{}{
		"location": location,
		"project":  project,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyHandles_basic(context, diskFilter),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve the bootstrapped KMS key handles used by the test
					resource.TestCheckResourceAttr("data.google_kms_key_handles.all_key_handles", "id", ""),
					//resource.TestMatchResourceAttr("data.google_kms_key_rings.all_key_rings", "key_rings.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyHandles_basic(context map[string]interface{}, filter string) string {
	context["filter"] = filter
	str := acctest.Nprintf(`
data "google_kms_key_handles" "all_key_handles" {
  location = "%{location}"
  project = "%{project}"
  filter = "%{filter}"
}
`, context)
	return str
}
