package dataplex_test

import (
	"testing"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexGlossary_dataplexGlossaryBasicExample(t *testing.T) {
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
				Config: testAccDataplexGlossary_dataplexGlossaryBasicExample(context),
			},
			{
				ResourceName:            "google_dataplex_glossary.glossary_test_id",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"glossary_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexGlossary_dataplexGlossaryBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_glossary" "glossary_test_id" {
  glossary_id = "tf-test-glossary-basic%{random_suffix}"
}
`, context)
}

func TestAccDataplexGlossary_dataplexGlossaryFullExample(t *testing.T) {
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
				Config: testAccDataplexGlossary_dataplexGlossaryFullExample(context),
			},
			{
				ResourceName:            "google_dataplex_glossary.glossary_test_id_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"glossary_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexGlossary_dataplexGlossaryFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_glossary" "glossary_test_id_full" {
  glossary_id = "tf-test-glossary-full%{random_suffix}"
  labels = { "tag": "test-tf" }
  display_name = "terraform glossary"
  description = "glossary created by Terraform"
}
`, context)
}
