// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleQuotaInfos_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_cloud_quotas_quota_infos.my_quota_infos"
	project := envvar.GetTestProjectFromEnv()
	service := "libraryagent.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleQuotaInfos_basic(project, service),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_infos.#", "9"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.%"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.quota_id"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.metric"),
					resource.TestCheckResourceAttr(resourceName, "quota_infos.0.service", service),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.is_precise"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.container_type"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.quota_increase_eligibility.0.is_eligible"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.dimensions_infos.0.details.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.dimensions_infos.0.applicable_locations.0"),
				),
			},
			{
				Config: testAccDataSourceGoogleQuotaInfos_withPagination(project, service),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_infos.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.%"),
					resource.TestCheckResourceAttr(resourceName, "next_page_token", "3"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleQuotaInfos_basic(project, service string) string {
	return acctest.Nprintf(`
	data "google_cloud_quotas_quota_infos" "my_quota_infos" {
		parent      = "projects/%{project}"	
		service 	= "%{service}"
	}
`, map[string]interface{}{"project": project, "service": service})
}

func testAccDataSourceGoogleQuotaInfos_withPagination(project, service string) string {
	return acctest.Nprintf(`
	data "google_cloud_quotas_quota_infos" "my_quota_infos" {
		parent      = "projects/%{project}"	
		service 	= "%{service}"
		page_size	= 2
		page_token	= 2
	}
`, map[string]interface{}{"project": project, "service": service})
}
