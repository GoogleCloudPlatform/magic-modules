// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
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
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_basic(context),
			},
			{
				ResourceName:            "google_cloud_quotas_quota_preference.preference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "quota_preference_id", "allow_missing", "validate_only", "ignore_safety_checks"},
			},
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

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_project" "my_project" {
			project_id 		= "tf-test%{random_suffix}"
			name       		= "tf-test%{random_suffix}"
			org_id          = "%{org_id}"
			billing_account = "%{billing_account}"
		}

		resource "google_project_iam_binding" "project" {
			project = google_project.my_project.project_id
			role    = "roles/cloudquotas.admin"

			members = [
				"user:liulola@google.com",
			]
		}

		# Wait for project being created.
		resource "time_sleep" "wait_120_seconds" {
			depends_on = [google_project_iam_binding.project]
			create_duration = "120s"
		}

		resource "google_cloud_quotas_quota_preference" "my-preference" {
			name			= "projects/tf-test%{random_suffix}/locations/global/quotaPreferences/compute_googleapis_com-CPUS-per-project-us-central1"
			quota_config {
				preferred_value = 50
			}
			dimensions		= { region = "us-central1" }
			service       	= "compute.googleapis.com"
			quota_id      	= "CPUS-per-project-region"
			contact_email 	= "liulola@google.com"
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_cloud_quotas_quota_preference" "my-preference" {
			name          = "projects/tf-test%{random_suffix}/locations/global/quotaPreferences/compute_googleapis_com-CPUS-per-project-us-central1"
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
