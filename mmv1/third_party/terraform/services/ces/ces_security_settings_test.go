// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCESSecuritySettings_cesSecuritySettingsExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckCESSecuritySettingsDestroyNoOp,
		Steps: []resource.TestStep{
			{
				Config: testAccCESSecuritySettings_cesSecuritySettingsExample_full(context),
			},
			{
				ResourceName:            "google_ces_security_settings.security_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "project"},
			},
			{
				Config: testAccCESSecuritySettings_cesSecuritySettingsExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_security_settings.security_settings", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_security_settings.security_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "project"},
			},
		},
	})
}

func testAccCheckCESSecuritySettingsDestroyNoOp(s *terraform.State) error {
	return nil
}

func testAccCESSecuritySettings_cesSecuritySettingsExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  provider        = google-beta
  project_id      = "tf-test-project-%{random_suffix}"
  name            = "tf-test-project-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "ces" {
  provider = google-beta
  project  = google_project.project.project_id
  service  = "ces.googleapis.com"
}

resource "google_ces_security_settings" "security_settings" {
  provider = google-beta
  project  = google_project.project.project_id
  location = "us"

  endpoint_control_policy {
    enforcement_scope = "ALWAYS"
    allowed_origins   = ["https://example.com", "https://google.com"]
  }

  depends_on = [google_project_service.ces]
}
`, context)
}

func testAccCESSecuritySettings_cesSecuritySettingsExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  provider        = google-beta
  project_id      = "tf-test-project-%{random_suffix}"
  name            = "tf-test-project-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "ces" {
  provider = google-beta
  project  = google_project.project.project_id
  service  = "ces.googleapis.com"
}

resource "google_ces_security_settings" "security_settings" {
  provider = google-beta
  project  = google_project.project.project_id
  location = "us"

  endpoint_control_policy {
    enforcement_scope = "VPCSC_ONLY"
    allowed_origins   = ["https://google.com"]
  }

  depends_on = [google_project_service.ces]
}
`, context)
}
