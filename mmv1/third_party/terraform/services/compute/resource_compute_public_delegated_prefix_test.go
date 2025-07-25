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
	parentProject := "byoipv6-fr-prober"
	parentRegion := "us-central1"
	parentName := "tf-test-delegation-mode-sub-pdp"

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
					// First, a basic check that the sub-prefix was created
					resource.TestCheckResourceAttrSet(subPrefixResourceName, "id"),

					// Now, the custom check function
					testAccCheckParentHasSubPrefix(t, parentProject, parentRegion, parentName, subPrefixResourceName),
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

resource "google_compute_public_delegated_prefix" "subprefix" {
  name = "tf-test-sub-prefix-1%{random_suffix}"
  description = "A nested address"
  region = "us-central1"
  ip_cidr_range = "2600:1901:4500:1:2::/96"
  parent_prefix = "%{parent_pdp_id}"
  mode = "DELEGATION"
  allocatable_prefix_length = 64
}
`, context)
}

func testAccCheckParentHasSubPrefix(t *testing.T, project, region, parentName, subPrefixResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the newly created sub-prefix resource in the Terraform state
		rs, ok := s.RootModule().Resources[subPrefixResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", subPrefixResourceName)
		}
		newSubPrefixName := rs.Primary.Attributes["name"]

		// Get the authenticated GCP client from the test helper
		config := acctest.GoogleProviderConfig(t)
		computeService, err := config.NewComputeClient(config.UserAgent)
		if err != nil {
			return err
		}

		// Make a direct API call to GCP to get the parent's details
		parent, err := computeService.PublicDelegatedPrefixes.Get(project, region, parentName).Do()
		if err != nil {
			return err
		}

		// Loop through the list of sub-prefixes returned by the API
		for _, sub := range parent.PublicDelegatedSubPrefixs {
			if sub.Name == newSubPrefixName {
				// Success! We found the new sub-prefix in the parent's list.
				return nil
			}
		}

		// If the loop finishes without finding the sub-prefix, return an error.
		return fmt.Errorf("sub-prefix %q not found in parent %q's sub-prefix list", newSubPrefixName, parentName)
	}
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
