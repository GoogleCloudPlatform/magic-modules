package backupdr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/backupdr"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccBackupDRAutoProtectionPolicyBinding_scopeChangeForcesReplacement(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBackupDRAutoProtectionPolicyBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDRAutoProtectionPolicyBinding_scopeProjectID(context),
			},
			{
				ResourceName:            "google_backup_dr_auto_protection_policy_binding.my_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_protection_policy_binding_id", "auto_protection_policy_id", "location", "wait_for_delete"},
			},
			{
				Config: testAccBackupDRAutoProtectionPolicyBinding_scopeProjectIDWaitForDelete(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy_binding.my_binding", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicyBinding_scopeDifferentProjectID(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy_binding.my_binding", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicyBinding_scopeDifferentProjectNumber(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy_binding.my_binding", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func TestAccBackupDRAutoProtectionPolicyBinding_numericScope(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBackupDRAutoProtectionPolicyBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDRAutoProtectionPolicyBinding_scopeProjectNumber(context),
			},
			{
				Config:             testAccBackupDRAutoProtectionPolicyBinding_scopeDifferentProjectNumber(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_backup_dr_auto_protection_policy_binding.my_binding", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func testAccBackupDRAutoProtectionPolicyBinding_scopeProjectID(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-scope-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-scope-%{random_suffix}"
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
  auto_protection_policy_id = "tf-test-app-scope-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "backup-policy"
        values = ["tf-test-scope-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = google_backup_dr_backup_plan.my_backup_plan.name
  }
}

resource "google_backup_dr_auto_protection_policy_binding" "my_binding" {
  location                          = "us-central1"
  auto_protection_policy_id         = google_backup_dr_auto_protection_policy.my_policy.auto_protection_policy_id
  auto_protection_policy_binding_id = "tf-test-binding-scope-%{random_suffix}"
  scope                             = "projects/${data.google_project.project.project_id}"
  wait_for_delete                   = false
}
`, context)
}

func testAccBackupDRAutoProtectionPolicyBinding_scopeProjectIDWaitForDelete(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-scope-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-scope-%{random_suffix}"
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
  auto_protection_policy_id = "tf-test-app-scope-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "backup-policy"
        values = ["tf-test-scope-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = google_backup_dr_backup_plan.my_backup_plan.name
  }
}

resource "google_backup_dr_auto_protection_policy_binding" "my_binding" {
  location                          = "us-central1"
  auto_protection_policy_id         = google_backup_dr_auto_protection_policy.my_policy.auto_protection_policy_id
  auto_protection_policy_binding_id = "tf-test-binding-scope-%{random_suffix}"
  scope                             = "projects/${data.google_project.project.project_id}"
  wait_for_delete                   = true
}
`, context)
}

func testAccBackupDRAutoProtectionPolicyBinding_scopeProjectNumber(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-scope-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-scope-%{random_suffix}"
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
  auto_protection_policy_id = "tf-test-app-scope-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "backup-policy"
        values = ["tf-test-scope-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = google_backup_dr_backup_plan.my_backup_plan.name
  }
}

resource "google_backup_dr_auto_protection_policy_binding" "my_binding" {
  location                          = "us-central1"
  auto_protection_policy_id         = google_backup_dr_auto_protection_policy.my_policy.auto_protection_policy_id
  auto_protection_policy_binding_id = "tf-test-binding-scope-%{random_suffix}"
  scope                             = "projects/${data.google_project.project.number}"
  wait_for_delete                   = true
}
`, context)
}

func testAccBackupDRAutoProtectionPolicyBinding_scopeDifferentProjectID(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-scope-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-scope-%{random_suffix}"
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
  auto_protection_policy_id = "tf-test-app-scope-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "backup-policy"
        values = ["tf-test-scope-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = google_backup_dr_backup_plan.my_backup_plan.name
  }
}

resource "google_backup_dr_auto_protection_policy_binding" "my_binding" {
  location                          = "us-central1"
  auto_protection_policy_id         = google_backup_dr_auto_protection_policy.my_policy.auto_protection_policy_id
  auto_protection_policy_binding_id = "tf-test-binding-scope-%{random_suffix}"
  scope                             = "projects/different-project-id"
  wait_for_delete                   = true
}
`, context)
}

func testAccBackupDRAutoProtectionPolicyBinding_scopeDifferentProjectNumber(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                   = "us-central1"
  backup_vault_id                            = "tf-test-bv-app-scope-%{random_suffix}"
  backup_minimum_enforced_retention_duration = "100000s"
}

resource "google_backup_dr_backup_plan" "my_backup_plan" {
  location       = "us-central1"
  backup_plan_id = "tf-test-bp-app-scope-%{random_suffix}"
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
  auto_protection_policy_id = "tf-test-app-scope-%{random_suffix}"

  criteria {
    matching_conditions {
      label_condition {
        key    = "backup-policy"
        values = ["tf-test-scope-%{random_suffix}"]
      }
    }
  }

  backup_plan_details {
    resource_type = "compute.googleapis.com/Instance"
    backup_plan   = google_backup_dr_backup_plan.my_backup_plan.name
  }
}

resource "google_backup_dr_auto_protection_policy_binding" "my_binding" {
  location                          = "us-central1"
  auto_protection_policy_id         = google_backup_dr_auto_protection_policy.my_policy.auto_protection_policy_id
  auto_protection_policy_binding_id = "tf-test-binding-scope-%{random_suffix}"
  scope                             = "projects/999999999999"
  wait_for_delete                   = true
}
`, context)
}
