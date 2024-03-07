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
			{
				Config: testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_upsert(context),
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
			org_id			= "%{org_id}"
			billing_account	= "%{billing_account}"
		}

		resource "google_project_iam_binding" "project_iam" {
			project 	= google_project.new_project.project_id
			role    	= "roles/cloudquotas.admin"
			members 	= [
				"user:testinguser2@google.com"
			]
			depends_on	= [google_project.new_project]
		}

		resource "google_project_service" "cloudquotas" {
			project  	= google_project.new_project.project_id
			service 	= "cloudquotas.googleapis.com"
			depends_on	= [google_project_iam_binding.project_iam]
		}

		resource "google_project_service" "compute" {
			project  	= google_project.new_project.project_id
			service 	= "compute.googleapis.com"
			depends_on	= [google_project_service.cloudquotas]
		}

		resource "google_project_service" "billing" {
			project  	= google_project.new_project.project_id
			service 	= "cloudbilling.googleapis.com"
			depends_on	= [google_project_service.compute]
		}
		
		resource "time_sleep" "wait_180_seconds" {
			create_duration	= "180s"
			depends_on		= [
				google_project_service.billing, 
			]
		}

		resource "google_cloud_quotas_quota_preference" "my_preference"{
			parent				= "projects/${google_project_service.billing.project}"
			name 				= "compute_googleapis_com-CPUS-per-project_us-central1"
			dimensions          = { region = "us-central1" }
			service             = "compute.googleapis.com"
			quota_id            = "CPUS-per-project-region"
			contact_email       = "testinguser2@google.com"
			quota_config  {
				preferred_value = 70
			}
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_increaseQuota(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_project" "new_project" {
			project_id 		= "tf-test%{random_suffix}"
			name       		= "tf-test%{random_suffix}"
			org_id			= "%{org_id}"
			billing_account	= "%{billing_account}"
		}

		resource "google_project_iam_binding" "project_iam" {
			project 	= google_project.new_project.project_id
			role    	= "roles/cloudquotas.admin"
			members 	= [
				"user:testinguser2@google.com"
			]
			depends_on	= [google_project.new_project]
		}

		resource "google_project_service" "cloudquotas" {
			project  	= google_project.new_project.project_id
			service 	= "cloudquotas.googleapis.com"
			depends_on	= [google_project_iam_binding.project_iam]
		}

		resource "google_project_service" "compute" {
			project  	= google_project.new_project.project_id
			service 	= "compute.googleapis.com"
			depends_on	= [google_project_service.cloudquotas]
		}

		resource "google_project_service" "billing" {
			project  	= google_project.new_project.project_id
			service 	= "cloudbilling.googleapis.com"
			depends_on	= [google_project_service.compute]
		}

		resource "time_sleep" "wait_180_seconds" {
			create_duration	= "180s"
			depends_on		= [
				google_project_service.billing, 
			]
		}

		resource "google_cloud_quotas_quota_preference" "my_preference"{
			dimensions			= { region = "us-central1" }
			contact_email		= "testinguser2@google.com"
			justification		= "Ignore. Increase quota for Terraform testing."
			quota_config  {
				preferred_value = 72
				annotations 	= { label = "terraform" }
			}
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_decreaseQuota(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_project" "new_project" {
			project_id 		= "tf-test%{random_suffix}"
			name       		= "tf-test%{random_suffix}"
			org_id			= "%{org_id}"
			billing_account	= "%{billing_account}"
		}

		resource "google_project_iam_binding" "project_iam" {
			project 	= google_project.new_project.project_id
			role    	= "roles/cloudquotas.admin"
			members 	= [
				"user:testinguser2@google.com"
			]
			depends_on	= [google_project.new_project]
		}

		resource "google_project_service" "cloudquotas" {
			project  	= google_project.new_project.project_id
			service 	= "cloudquotas.googleapis.com"
			depends_on	= [google_project_iam_binding.project_iam]
		}

		resource "google_project_service" "compute" {
			project  	= google_project.new_project.project_id
			service 	= "compute.googleapis.com"
			depends_on	= [google_project_service.cloudquotas]
		}

		resource "google_project_service" "billing" {
			project  	= google_project.new_project.project_id
			service 	= "cloudbilling.googleapis.com"
			depends_on	= [google_project_service.compute]
		}

		resource "time_sleep" "wait_180_seconds" {
			create_duration	= "180s"
			depends_on		= [
				google_project_service.billing, 
			]
		}

		resource "google_cloud_quotas_quota_preference" "my_preference"{
			dimensions				= { region = "us-central1" }
			contact_email			= "testinguser2@google.com"
			ignore_safety_checks	= "QUOTA_DECREASE_PERCENTAGE_TOO_HIGH"
			quota_config  {
				preferred_value		= 65
			}
		}
	`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_upsert(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "new_project" {
			project_id 		= "tf-test%{random_suffix}"
			name       		= "tf-test%{random_suffix}"
			org_id			= "%{org_id}"
			billing_account	= "%{billing_account}"
		}

		resource "google_project_iam_binding" "project_iam" {
			project = google_project.new_project.project_id
			role    = "roles/cloudquotas.admin"
			members = [
				"user:testinguser3@google.com"
			]
			depends_on	= [google_project.new_project]
		}

		resource "google_project_service" "cloudquotas" {
			project  	= google_project.new_project.project_id
			service 	= "cloudquotas.googleapis.com"
			depends_on	= [google_project_iam_binding.project_iam]
		}

		resource "google_project_service" "compute" {
			project  	= google_project.new_project.project_id
			service 	= "compute.googleapis.com"
			depends_on	= [google_project_service.cloudquotas]
		}

		resource "google_project_service" "billing" {
			project  	= google_project.new_project.project_id
			service 	= "cloudbilling.googleapis.com"
			depends_on	= [google_project_service.compute]
		}

		resource "time_sleep" "wait_180_seconds" {
			create_duration	= "180s"
			depends_on		= [
				google_project_service.billing, 
			]
		}

		resource "google_cloud_quotas_quota_preference" "my_preference"{
			parent				= "projects/${google_project_service.billing.project}"
			name 				= "libraryagent_googleapis_com-ReadsPerMinutePerProject"
			service				= "libraryagent.googleapis.com"
			quota_id			= "ReadsPerMinutePerProject"
			contact_email		= "testinguser3@google.com"
			quota_config  {
				preferred_value	= 15
			}
			allow_missing		= true
		}
	`, context)
}
