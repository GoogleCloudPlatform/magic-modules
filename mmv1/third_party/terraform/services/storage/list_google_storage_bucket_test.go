// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageBucketListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketListResource_queryIdentity_prereq(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_storage_bucket.prereq", "name", bucketName),
					resource.TestCheckResourceAttr("google_storage_bucket.prereq", "project", project),
				),
			},
			{
				Query:  true,
				Config: testAccStorageBucketListResource_queryIdentity_list(project, bucketName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_storage_bucket.prefixed_list", map[string]knownvalue.Check{
						"name":    knownvalue.StringExact(bucketName),
						"project": knownvalue.StringExact(project),
					}),
					querycheck.ExpectLengthAtLeast("google_storage_bucket.prefixed_list", 1),
				},
			},
		},
	})
}

func testAccStorageBucketListResource_queryIdentity_prereq(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "prereq" {
  name          = "%[1]s"
  location      = "US"
  force_destroy = true
}
`, bucketName)
}

func testAccStorageBucketListResource_queryIdentity_list(project, bucketName string) string {
	return fmt.Sprintf(`
list "google_storage_bucket" "prefixed_list" {
  provider = google

  config {
    project                  = "%[1]s"
    prefix                   = "%[2]s"
    max_results              = 500
    projection               = "noAcl"
    return_partial_success   = true
    soft_deleted             = false
  }
}
`, project, bucketName)
}
