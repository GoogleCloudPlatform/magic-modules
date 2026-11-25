package backupdr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/backupdr"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccBackupDRAutoProtectionPolicy_backupPlanProjectChange(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBackupDRAutoProtectionPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDRAutoProtectionPolicy_backupPlanProjectID(context),
			},
			{
				ResourceName:            "google_backup_dr_auto_protection_policy.my_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_protection_policy_id", "location"},
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicy_backupPlanProjectNumber(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy.my_policy", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicy_backupPlanDifferentProjectID(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy.my_policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicy_backupPlanDifferentProjectNumber(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy.my_policy", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func testAccBackupDRAutoProtectionPolicy_backupPlanProjectID(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-bpproj-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-bpproj-%{random_suffix}"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id               = "rule-1"
    backup_retention_days = 5

    standard_schedule {
      recurrence_type  = "HOURLY"
      hourly_frequency = 6
      time_zone        = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 24
      }
    }
  }
}

resource "google_backup_dr_auto_protection_policy" "my_policy" {
  location                  = "us-central1"
  auto_protection_policy_id = "tf-test-app-bpproj-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "env"
        values = ["tf-test-bpproj-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = "projects/${data.google_project.project.project_id}/locations/us-central1/backupPlans/${google_backup_dr_backup_plan.my_backup_plan.backup_plan_id}"
  }
}
`, context)
}

func testAccBackupDRAutoProtectionPolicy_backupPlanProjectNumber(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-bpproj-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-bpproj-%{random_suffix}"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id               = "rule-1"
    backup_retention_days = 5

    standard_schedule {
      recurrence_type  = "HOURLY"
      hourly_frequency = 6
      time_zone        = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 24
      }
    }
  }
}

resource "google_backup_dr_auto_protection_policy" "my_policy" {
  location                  = "us-central1"
  auto_protection_policy_id = "tf-test-app-bpproj-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "env"
        values = ["tf-test-bpproj-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = "projects/${data.google_project.project.number}/locations/us-central1/backupPlans/${google_backup_dr_backup_plan.my_backup_plan.backup_plan_id}"
  }
}
`, context)
}

func testAccBackupDRAutoProtectionPolicy_backupPlanDifferentProjectID(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-bpproj-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-bpproj-%{random_suffix}"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id               = "rule-1"
    backup_retention_days = 5

    standard_schedule {
      recurrence_type  = "HOURLY"
      hourly_frequency = 6
      time_zone        = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 24
      }
    }
  }
}

resource "google_backup_dr_auto_protection_policy" "my_policy" {
  location                  = "us-central1"
  auto_protection_policy_id = "tf-test-app-bpproj-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "env"
        values = ["tf-test-bpproj-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = "projects/different-project-id/locations/us-central1/backupPlans/${google_backup_dr_backup_plan.my_backup_plan.backup_plan_id}"
  }
}
`, context)
}

func testAccBackupDRAutoProtectionPolicy_backupPlanDifferentProjectNumber(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-bpproj-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-bpproj-%{random_suffix}"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id               = "rule-1"
    backup_retention_days = 5

    standard_schedule {
      recurrence_type  = "HOURLY"
      hourly_frequency = 6
      time_zone        = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 24
      }
    }
  }
}

resource "google_backup_dr_auto_protection_policy" "my_policy" {
  location                  = "us-central1"
  auto_protection_policy_id = "tf-test-app-bpproj-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "env"
        values = ["tf-test-bpproj-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = "projects/999999999999/locations/us-central1/backupPlans/${google_backup_dr_backup_plan.my_backup_plan.backup_plan_id}"
  }
}
`, context)
}
