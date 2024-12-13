package gkehub2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataSourceGoogleGKEHub2MembershipBinding_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHub2MembershipBindingDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleGKEHub2MembershipBinding_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_gke_hub_membership_binding.example", "google_gke_hub_membership_binding.example"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleGKEHub2MembershipBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_membership" "example" {
  membership_id = "tf-test-membership%{random_suffix}"
}

resource "google_gke_hub_scope" "example" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_membership_binding" "example" {
  membership_binding_id = "tf-test-membership-binding%{random_suffix}"
  scope = google_gke_hub_scope.example.name
  membership_id = "tf-test-membership%{random_suffix}"
  location = "global"
  labels = {
      keyb = "valueb"
      keya = "valuea"
      keyc = "valuec" 
  }
  depends_on = [
    google_gke_hub_membership.example,
    google_gke_hub_scope.example
  ]
}

data "google_gke_hub_membership_binding" "example" {
  location = google_gke_hub_membership_binding.example.location
  membership_id = google_gke_hub_membership_binding.example.membership_id
  membership_binding_id = google_gke_hub_membership_binding.example.membership_binding_id
}
`, context)
}

func testAccCheckGKEHub2MembershipBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_hub_membership_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{GKEHub2BasePath}}projects/{{project}}/locations/{{location}}/memberships/{{membership_id}}/bindings/{{membership_binding_id}}")
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
				return fmt.Errorf("GKEHub2MembershipBinding still exists at %s", url)
			}
		}

		return nil
	}
}
