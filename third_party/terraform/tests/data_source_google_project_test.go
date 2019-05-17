package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleProject_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	project := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleProjectCheck("data.google_project.project", "google_project.project"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleProjectCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		errMsg := ""
		for k, attr := range rsAttr {
			if dsAttr[k] != attr {
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], attr)
			}
		}

		if errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		return nil
	}
}

func testAccCheckGoogleProjectConfig(project, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "%s"
	name = "%s"
	org_id = "%s"
}
	
data "google_project" "project" {
	project_id = "${google_project.project.project_id}"
}`, project, project, org)
}
