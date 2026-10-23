package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceFindings_empty(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageControlProjectIntelligenceFindings_empty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_storage_control_project_intelligence_findings.empty", "findings.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceFindings_empty() string {
	return `
data "google_storage_control_project_intelligence_findings" "empty" {
}
`
}
