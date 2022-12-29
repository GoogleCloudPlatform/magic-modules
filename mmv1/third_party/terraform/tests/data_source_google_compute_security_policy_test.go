package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeSecurityPolicy_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeSecurityPolicy_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_compute_security_policy.foo", "google_compute_security_policy.policy"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeSecurityPolicy_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_security_policy" "policy" {
   name = "my-policy%{random_suffix}"

   rule {
     action   = "deny(403)"
     priority = "1000"
     match {
       versioned_expr = "SRC_IPS_V1"
       config {
	    src_ip_ranges = ["9.9.9.0/24"]
       }
      }
     description = "Deny access to IPs in 9.9.9.0/24"
    }

   rule {
     action   = "allow"
     priority = "2147483647"
     match {
     versioned_expr = "SRC_IPS_V1"
     config {
	   src_ip_ranges = ["*"]
     }
     }
     description = "default rule"
    }
}

data "google_compute_security_policy" "foo" {
  name = google_compute_security_policy.policy.name
  project = google_compute_security_policy.policy.project
}`, context)

}
