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


func TestAccAzureRmResourceGroup_containerAnalysisNoteBasicExample(t *testing.T) {
  t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureRmResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRmResourceGroup_containerAnalysisNoteBasicExample(acctest.RandString(10)),
			},
			{
				ResourceName:      "google_azure_rm_resource_group.note",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRmResourceGroup_containerAnalysisNoteBasicExample(val string) string {
  return fmt.Sprintf(`
resource "google_container_analysis_note" "note" {
  name = "test-attestor-note-%s"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}
`, val,
  )
}


func testAccCheckAzureRmResourceGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_azure_rm_resource_group" {
			continue
		}

	config := testAccProvider.Meta().(*Config)

	url, err := replaceVarsForTest(rs, "https://pubsub.googleapis.com/v1/projects/{{project}}/topics/{{name}}")
	if err != nil {
		return err
	}

	_, err = sendRequest(config, "GET", url, nil)
	if err == nil {
			return fmt.Errorf("AzureRmResourceGroup still exists at %s", url)
		}
	}

	return nil
}
