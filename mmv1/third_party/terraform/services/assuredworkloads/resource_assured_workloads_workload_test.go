package assuredworkloads_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccAssuredWorkloadsWorkload_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy:             testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "workload_options", "provisioned_resources_parent", "partner_services_billing_account", "labels", "terraform_labels"},
			},
			{
				Config: testAccAssuredWorkloadsWorkload_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "workload_options", "provisioned_resources_parent", "partner_services_billing_account", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccAssuredWorkloadsWorkload_FullHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_FullHandWritten(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "workload_options", "provisioned_resources_parent", "partner_services_billing_account", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAssuredWorkloadsWorkload_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  labels = {
    a = "a"
  }
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
  workload_options {
    kaj_enrollment_type = "KEY_ACCESS_TRANSPARENCY_OFF"
  }
  resource_settings {
    resource_type = "CONSUMER_FOLDER"
    display_name = "folder-display-name"
  }
  violation_notifications_enabled = true
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_folder.folder1]
}
`, context)
}

func testAccAssuredWorkloadsWorkload_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  labels = {
    a = "b"
  }
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
  resource_settings {
    resource_type = "CONSUMER_FOLDER"
    display_name = "folder-display-name"
  }
  violation_notifications_enabled = true
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_folder.folder1]
}
`, context)
}

func testAccAssuredWorkloadsWorkload_FullHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  organization = "%{org_id}"
  location = "us-central1"
  kms_settings {
    next_rotation_time = "2022-10-02T15:01:23Z"
    rotation_period = "864000s"
  }
  provisioned_resources_parent = google_folder.folder1.name
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_folder.folder1]
}
`, context)
}

