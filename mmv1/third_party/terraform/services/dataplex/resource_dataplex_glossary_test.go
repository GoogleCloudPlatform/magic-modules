package dataplex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
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

func testAccCheckDataplexGlossaryDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataplex_glossary" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DataplexBasePath}}projects/{{project}}/locations/{{location}}/glossaries/{{glossary_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DataplexGlossary still exists at %s", url)
			}
		}

		return nil
	}
}
