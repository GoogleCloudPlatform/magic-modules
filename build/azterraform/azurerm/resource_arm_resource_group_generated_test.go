// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package azurerm

import (
  "fmt"
  "testing"

  "github.com/hashicorp/terraform/helper/acctest"
  "github.com/hashicorp/terraform/helper/resource"
)


func TestAccArmResourceGroup_resourceGroupExample(t *testing.T) {
  t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckArmResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArmResourceGroup_resourceGroupExample(acctest.RandString(10)),
			},
			{
				ResourceName:      "google_arm_resource_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccArmResourceGroup_resourceGroupExample(val string) string {
  return fmt.Sprintf(`
resource "azurerm_resource_group" "example" {
  name     = "ExampleRG-%s"
  location = "West US"

  tags {
    environment = "Production"
  }
}
`, val,
  )
}


func testAccCheckArmResourceGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_arm_resource_group" {
			continue
		}

	config := testAccProvider.Meta().(*Config)

	url, err := replaceVarsForTest(rs, "NotUsedInAzureNotUsedInAzure/{{name}}")
	if err != nil {
		return err
	}

	_, err = sendRequest(config, "GET", url, nil)
	if err == nil {
			return fmt.Errorf("ArmResourceGroup still exists at %s", url)
		}
	}

	return nil
}
