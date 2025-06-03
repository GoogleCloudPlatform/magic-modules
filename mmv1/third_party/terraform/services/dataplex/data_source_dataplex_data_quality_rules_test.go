package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataplexDataQualityRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":      "dataplex-back-end-dev-project",
		"location":     "us-central1",
		"data_scan_id": "a111",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDataQualityRules_config(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dataplex_data_quality_rules.generated_dq_rules", "rules.#", "17"),
				),
			},
		},
	})
}

func testAccDataplexDataQualityRules_config(context map[string]interface{}) string {
	return acctest.Nprintf(`
		data "google_dataplex_data_quality_rules" "generated_dq_rules" {
			project		 = "%{project}"
			location	 = "%{location}"
			data_scan_id = "%{data_scan_id}"
		}`, context)
}
