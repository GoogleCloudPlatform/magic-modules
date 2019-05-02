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

    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationNullVariable_basic(t *testing.T) {
    resourceName := "azurerm_automation_null_variable.test"
    ri := tf.AccRandTimeInt()
    location := testLocation()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMAutomationNullVariableExists(resourceName),
                ),
            },
            {
                ResourceName:      resourceName,
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func TestAccAzureRMAutomationNullVariable_complete(t *testing.T) {
    resourceName := "azurerm_automation_null_variable.test"
    ri := tf.AccRandTimeInt()
    location := testLocation()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMAutomationNullVariable_complete(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMAutomationNullVariableExists(resourceName),
                    resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
                ),
            },
            {
                ResourceName:      resourceName,
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func TestAccAzureRMAutomationNullVariable_basicCompleteUpdate(t *testing.T) {
    resourceName := "azurerm_automation_null_variable.test"
    ri := tf.AccRandTimeInt()
    location := testLocation()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMAutomationNullVariableExists(resourceName),
                ),
            },
            {
                Config: testAccAzureRMAutomationNullVariable_complete(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMAutomationNullVariableExists(resourceName),
                    resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
                ),
            },
            {
                Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMAutomationNullVariableExists(resourceName),
                ),
            },
        },
    })
}


func testCheckAzureRMAutomationNullVariableExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("Automation Null Variable not found: %s", resourceName)
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]
        accountName := rs.Primary.Attributes["automation_account_name"]

        client := testAccProvider.Meta().(*ArmClient).automationVariableClient
        ctx := testAccProvider.Meta().(*ArmClient).StopContext

        if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
            if utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Automation Null Variable %q (Automation Account Name %q / Resource Group %q) does not exist", name, accountName, resourceGroup)
            }
            return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
        }

        return nil
    }
}

func testCheckAzureRMAutomationNullVariableDestroy(s *terraform.State) error {
    client := testAccProvider.Meta().(*ArmClient).automationVariableClient
    ctx := testAccProvider.Meta().(*ArmClient).StopContext

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "azurerm_automation_null_variable" {
            continue
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]
        accountName := rs.Primary.Attributes["automation_account_name"]

        if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
            }
        }

        return nil
    }

    return nil
}

func testAccAzureRMAutomationNullVariable_basic(rInt int, location string) string {
    return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_null_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationNullVariable_complete(rInt int, location string) string {
    return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_null_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
}
`, rInt, location, rInt, rInt)
}
