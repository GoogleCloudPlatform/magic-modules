// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_update(context),
			},
			{
				ResourceName:            "google_cloud_quotas_quota_preference.preference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "quota_preference_id", "allow_missing", "validate_only", "ignore_safety_checks"},
			},
		},
	})
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_quotas_quota_preference" "my-preference" {
  name          = "projects/%{project}/locations/global/quotaPreferences/compute_googleapis_com-CPUS-per-project-us-central1"
  quota_config  {
    preferred_value = 200
  }
  dimensions 	= { region = "us-central1" }
  service       = "compute.googleapis.com"
  quota_id      = "CPUS-per-project-region"
  contact_email = "liulola@google.com"
  justification = "Increase quota for terraform testing."
  validate_only = true
}
`, context)
}
