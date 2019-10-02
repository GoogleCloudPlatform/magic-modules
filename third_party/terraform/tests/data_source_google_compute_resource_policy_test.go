package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceComputeResourcePolicy(t *testing.T) {
	t.Parallel()

	rsName := "foo"
	rsFullName := fmt.Sprintf("google_compute_resource_policy.%s", rsName)
	dsName := "my_policy"
	dsFullName := fmt.Sprintf("data.google_compute_resource_policy.%s", dsName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceComputeResourcePolicyDestroy(rsFullName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeResourcePolicyConfig(rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeResourcePolicyCheck(dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeResourcePolicyCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
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

		policyAttrsToTest := []string{
			"self_link",
			"name",
		}

		for _, attrToCheck := range policyAttrsToTest {
			if dsAttr[attrToCheck] != rsAttr[attrToCheck] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attrToCheck,
					dsAttr[attrToCheck],
					rsAttr[attrToCheck],
				)
			}
		}

		return nil
	}
}

func testAccCheckDataSourceComputeResourcePolicyDestroy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		policyAttrs := rs.Primary.Attributes

		_, err := config.clientCompute.ResourcePolicies.Get(
			config.Project, policyAttrs["region"], policyAttrs["name"]).Do()
		if err == nil {
			return fmt.Errorf("Resource Policy still exists")
		}

		return nil
	}
}

func testAccDataSourceComputeResourcePolicyConfig(rsName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "%s" {
  name = "policy"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time = "04:00"
      }
    }
  }
}

data "google_compute_resource_policy" "%s" {
  name     = "${google_compute_resource_policy.%s.name}"
  region   = "${google_compute_resource_policy.%s.region}"
}
`, rsName, dsName, rsName, rsName)
}
