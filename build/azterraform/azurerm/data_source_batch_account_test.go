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
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMBatchAccount_basic(t *testing.T) {
    dataSourceName := "data.azurerm_batch_account.test"
    ri := tf.AccRandTimeInt()
    location := testLocation()

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccDataSourceBatchAccount_basic(ri, location),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(dataSourceName, "poolAllocationMode", "BatchService"),
                    resource.TestCheckResourceAttrSet(dataSourceName, "storageAccountId"),
                ),
            },
        },
    })
}

func testAccDataSourceBatchAccount_basic(rInt int, location string) string {
    config := testAccBatchAccount_basic(rInt, location)
    return fmt.Sprintf(`
%s

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_batch_account.test.resource_group_name}"
}
`, config)
}
