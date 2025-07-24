package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputePublicDelegatedPrefix_computePublicDelegatedPrefixWithSubPrefixExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"parent_pdp_id": "projects/tf-static-byoip/regions/us-central1/publicDelegatedPrefixes/tf-test-delegation-mode-sub-pdp",
		"project":       "tf-static-byoip",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputePublicDelegatedPrefixWithSubPrefixDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputePublicDelegatedPrefix_computePublicDelegatedPrefixWithSubPrefixExample(context),
				Check: resource.ComposeTestCheckFunc(
					// 1. Verify that the sub-prefixes list contains exactly one block.
					resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.#", "1"),

					// 2. Verify the description of the first sub-prefix in the list.
					resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.description", "A nested address"),

					// 3. Verify the 'is_address' attribute of the first sub-prefix.
					resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.is_address", "true"),

					// 4. Verify that the name of the first sub-prefix is set to some value.
					resource.TestCheckResourceAttrSet("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.name"),
				),
			},
			{
				ResourceName:            "google_compute_public_delegated_prefix.prefixes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccComputePublicDelegatedPrefix_computePublicDelegatedPrefixWithSubPrefixExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_public_delegated_prefix" "prefixes" {
  name = "tf-test-prefix-with-sub-prefixes%{random_suffix}"
  region = "us-central1"
  description = "Public delegated prefix with sub prefix for testing "
  ip_cidr_range = "2600:1901:4500:2::/64"
  parent_prefix = "%{parent_pdp_id}"
  project = "%{project}"
  mode = "DELEGATION"

  
}

resource "google_compute_public_delegated_prefix" "subprefix" {
  name = "tf-test-sub-prefix-1%{random_suffix}"
  description = "A nested address"
  region = "us-central1"
  ip_cidr_range = "2600:1901:4500:1:2::/96"
  parent_prefix = google_compute_public_delegated_prefix.prefixes.id
  mode = "DELEGATION"
  allocatable_prefix_length = 64
}
`, context)
}

func testAccCheckComputePublicDelegatedPrefixWithSubPrefixDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_public_delegated_prefix" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/publicDelegatedPrefixes/{{name}}")
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
				return fmt.Errorf("ComputePublicDelegatedPrefix still exists at %s", url)
			}
		}

		return nil
	}
}
