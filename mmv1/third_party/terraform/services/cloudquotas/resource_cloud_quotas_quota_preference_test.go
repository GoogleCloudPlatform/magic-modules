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
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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
resource "google_cloud_quotas_quota_preference" "preference" {
  name          = "projects/%{project}/locations/global/quotaPreferences/tf-test-compute_googleapis_com-CPUS-per-project_us-central2%{random_suffix}"
  quota_config  {
    preferred_value = 210
  }
  dimensions = {
    region = "us-central1"
  }
  service       = "compute.googleapis.com"
  quota_id      = "CPUS-per-project-region"
  contact_email = "liulola@google.com"
  justification = "Increase quota for terraform testing."
}
`, context)
}

func testAccCloudQuotasQuotaPreference_cloudquotasQuotaPreferenceBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_quotas_quota_preference" "preference" {
  name          = "projects/%{project}/locations/global/quotaPreferences/tf-test-compute_CPUS-per-project_compute_googleapis_com-CPUS-per-project_us-east1%{random_suffix}"
  quota_config  {
    preferred_value = 200
  }
  dimensions = {
    region = "us-east1"
  }
  service       = "compute.googleapis.com"
  quota_id      = "CPUS-per-project-region"
  contact_email = "liulola@google.com"
  justification = "Increase quota for terraform testing."
  validate_only = true
  allow_missing = true
}
`, context)
}
