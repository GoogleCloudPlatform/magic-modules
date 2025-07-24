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
	subPrefixResourceName := "google_compute_public_delegated_prefix.subprefix"
	parentDataSourceName := "data.google_compute_public_delegated_prefix.parent"

	context := map[string]interface{}{
		"parent_pdp_id": "projects/tf-static-byoip/regions/us-central1/publicDelegatedPrefixes/tf-test-delegation-mode-sub-pdp",
		"project":       "byoipv6-fr-prober",
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
					// // 1. Verify that the sub-prefixes list contains exactly one block.
					// resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.#", "1"),

					// // 2. Verify the description of the first sub-prefix in the list.
					// resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.description", "A nested address"),

					// // 3. Verify the 'is_address' attribute of the first sub-prefix.
					// resource.TestCheckResourceAttr("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.is_address", "false"),

					// // 4. Verify that the name of the first sub-prefix is set to some value.
					// resource.TestCheckResourceAttrSet("google_compute_public_delegated_prefix.prefixes", "public_delegated_sub_prefixs.0.name"),

					// Verify that the parent's sub-prefix list now contains at least one item.
					// We check for "1" assuming a clean test environment.
					resource.TestCheckResourceAttr(parentDataSourceName, "public_delegated_sub_prefixes.#", "1"),

					// Verify that the 'name' of the first sub-prefix in the parent's list
					// matches the name of the sub-prefix resource we just created.
					resource.TestCheckResourceAttrPair(
						parentDataSourceName, "public_delegated_sub_prefixes.0.name",
						subPrefixResourceName, "name",
					),

					// Verify its IP range as well
					resource.TestCheckResourceAttrPair(
						parentDataSourceName, "public_delegated_sub_prefixes.0.ip_cidr_range",
						subPrefixResourceName, "ip_cidr_range",
					),
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
// resource "google_compute_public_delegated_prefix" "prefixes" {
//   name = "tf-test-prefix-with-sub-prefixes%{random_suffix}"
//   region = "us-central1"
//   description = "Public delegated prefix with sub prefix for testing "
//   ip_cidr_range = "2600:1901:4500:2::/64"
//   parent_prefix = "%{parent_pdp_id}"
//   project = "%{project}"
//   mode = "DELEGATION"

  
// }
data "google_compute_public_delegated_prefix" "parent" {
  name = "tf-test-delegation-mode-sub-pdp"
  region = "us-central1"
  project = "%{project}"
}
resource "google_compute_public_delegated_prefix" "subprefix" {
  name = "tf-test-sub-prefix-1%{random_suffix}"
  description = "A nested address"
  region = "us-central1"
  ip_cidr_range = "2600:1901:4500:1:2::/96"
  parent_prefix = data.google_compute_public_delegated_prefix.parent.id
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
