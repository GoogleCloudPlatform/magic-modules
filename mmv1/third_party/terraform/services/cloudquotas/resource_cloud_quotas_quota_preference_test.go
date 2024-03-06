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
				ResourceName:            "google_cloud_quotas_quota_preference.my_preference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "quota_preference_id", "allow_missing", "validate_only", "ignore_safety_checks", "contact_email"},
			},
			{
				Config: testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_increaseQuota(context),
			},
			{
				ResourceName:            "google_cloud_quotas_quota_preference.my_preference",
				ImportState:             true,
				ImportStateVerify:       true,
				ExpectNonEmptyPlan:      true,
				ImportStateVerifyIgnore: []string{"parent", "quota_preference_id", "allow_missing", "validate_only", "ignore_safety_checks", "contact_email", "justification", "quota_config.0.annotations"},
			},
			{
				Config: testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_decreaseQuota(context),
			},
			{
				ResourceName:            "google_cloud_quotas_quota_preference.my_preference",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "quota_preference_id", "allow_missing", "validate_only", "ignore_safety_checks", "contact_email", "justification"},
			},
		},
	})
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_project" "new_project" {
			project_id 		= "tf-test%{random_suffix}"
			name       		= "tf-test%{random_suffix}"
			org_id          = "%{org_id}"
			billing_account = "%{billing_account}"
		}

		resource "google_project_service" "cloudquotas" {
			project  = google_project.new_project.project_id
			service = "cloudquotas.googleapis.com"
			depends_on = [google_project.new_project]
		}

		resource "google_project_service" "billing" {
			project  = google_project.new_project.project_id
			service = "cloudbilling.googleapis.com"
			depends_on = [google_project.new_project]
		}

		resource "google_project_iam_binding" "project_iam" {
			project = google_project_service.cloudquotas.project
			role    = "roles/cloudquotas.admin"

			members = [
				"user:liulola@google.com"
			]
			depends_on = [google_project.new_project]
		}

		resource "time_sleep" "wait_120_seconds" {
			depends_on = [google_project_iam_binding.project_iam]
			create_duration = "120s"
		}

		resource "google_cloud_quotas_quota_preference" "my_preference"{
			parent				= "projects/${google_project_iam_binding.project_iam.project}"
			name 				= "compute_googleapis_com-A2-CPUS-per-project_asia-northeast1"
			dimensions          = { region = "asia-northeast1" }
			service             = "compute.googleapis.com"
			quota_id            = "A2-CPUS-per-project-region"
			contact_email       = "liulola@google.com"
			quota_config  {
				preferred_value = 12
			}
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_increaseQuota(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_cloud_quotas_quota_preference" "my_preference"{
			contact_email       = "liulola@google.com"
			justification		= "Increase quota for Terraform testing."
			quota_config  {
				preferred_value = 12
				annotations 	= { label = "terraform" }
			}
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_decreaseQuota(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_cloud_quotas_quota_preference" "my_preference"{
			contact_email			= "liulola@google.com"
			ignore_safety_checks	= "QUOTA_DECREASE_PERCENTAGE_TOO_HIGH"
			quota_config  {
				preferred_value 	= 10
			}
		}
	`, context)
}
