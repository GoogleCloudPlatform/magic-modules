package securitycenterv2_test

import (
	// "math/rand"
	// "strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2OrganizationSourceFinding_basic(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	// source_id := strconv.FormatInt(rand.Int63n(1e10), 10)

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		// "source_id":     source_id,
		"finding_id":    random_suffix,
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationSourceFinding_basic(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_source_finding.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterV2OrganizationSourceFinding_update(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_source_finding.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterV2OrganizationSourceFinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_scc_v2_organization_source" "default" {
  organization = "%{org_id}"
  display_name = "TF Test Source-%{random_suffix}"
  description  = "Test source for findings"
}

// locals {
//   source_id = regex("sources/(.+)$", google_scc_v2_organization_source.default.name)[0]
// }

locals {
  source_name_parts = regex("sources/([0-9]+)$", google_scc_v2_organization_source.default.name)
  source_id         = length(local.source_name_parts) > 0 ? local.source_name_parts[0] : "INVALID"
}

resource "google_scc_v2_organization_source_finding" "default" {
  organization = "%{org_id}"
  source       = local.source_id
  location     = "global"
  finding_id   = "%{finding_id}"
  
  state        = "ACTIVE"
  category     = "MEDIUM_RISK_ONE"
  event_time   = "2024-01-01T00:00:00Z"
}
`, context)
}

func testAccSecurityCenterV2OrganizationSourceFinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_scc_v2_organization_source" "default" {
  organization = "%{org_id}"
  display_name = "TF Test Source-%{random_suffix}"
  description  = "Test source for findings"
}

// locals {
//   source_id = regex("sources/(.+)$", google_scc_v2_organization_source.default.name)[0]
// }

locals {
  source_name_parts = regex("sources/([0-9]+)$", google_scc_v2_organization_source.default.name)
  source_id         = length(local.source_name_parts) > 0 ? local.source_name_parts[0] : "INVALID"
}

resource "google_scc_v2_organization_source_finding" "default" {
  organization = "%{org_id}"
  source       = local.source_id
  location     = "global"
  finding_id   = "%{finding_id}"

  state        = "INACTIVE"
  category     = "XSS"
  event_time   = "2024-01-01T00:00:00Z"
}
`, context)
}
