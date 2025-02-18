package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexGlossary_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexGlossaryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexGlossary_full(context),
			},
			{
				ResourceName:            "google_dataplex_glossary.test_glossary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "glossary_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataplexGlossary_update(context),
			},
			{
				ResourceName:            "google_dataplex_glossary.test_glossary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "glossary_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexGlossary_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dataplex_glossary" "test_glossary" {
	glossary_id = "tf-test-glossary%{random_suffix}"
	project = "%{project_name}"
	location = "us-central1"
}
`, context)
}

func testAccDataplexGlossary_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dataplex_glossary" "test_glossary" {
	glossary_id = "tf-test-glossary%{random_suffix}"
	project = "%{project_name}"
	location = "us-central1"

	labels = {"tag": "test-tf"}
	display_name = "terraform glossary"
	description = "glossary created by Terraform"
}
`, context)
}
