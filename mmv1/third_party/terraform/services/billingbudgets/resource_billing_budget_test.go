package billingbudgets_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/billingbudgets"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccBillingBudget_billingBudgetCurrencycode(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestMasterBillingAccountFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetCurrencycode(context),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccBillingBudget_billingBudgetCurrencycode(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
    labels  = {
      label = "bar"
    }
  }

  amount {
    specified_amount {
      units         = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 1.0
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }
}
`, context)
}

func TestAccBillingBudget_billingBudgetUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestMasterBillingAccountFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetUpdateStart(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdateRemoveFilter(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetCalendarUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetCustomPeriodUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBillingBudget_billingFilterSubaccounts(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"master_billing_acct": envvar.GetTestMasterBillingAccountFromEnv(t),
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingFilterSubaccounts(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingFilterRemoveSubaccounts(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_billingFilterSubaccounts(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{master_billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_subaccount" "subaccount" {
  display_name = "My Billing Account"
  master_billing_account = data.google_billing_account.account.id
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
    labels  = {
      label = "bar"
    }

    subaccounts = ["billingAccounts/${google_billing_subaccount.subaccount.billing_account_id}"]
  }

  amount {
    specified_amount {
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 1.0
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }
}
`, context)
}

func testAccBillingBudget_billingFilterRemoveSubaccounts(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{master_billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
    labels  = {
      label = "bar"
    }
  }

  amount {
    specified_amount {
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 1.0
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdateStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}

data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
    labels  = {
      label = "bar"
    }
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic1.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdateRemoveFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = []
    labels = {}
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic1.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = []
    labels  = {
      label1 = "bar2"
    }
    credit_types_treatment = "INCLUDE_SPECIFIED_CREDITS"
    services               = ["services/24E6-581D-38E5"] # Bigquery
    credit_types           = ["PROMOTION", "FREE_TIER"]
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetCalendarUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = []
    labels  = {
      label1 = "bar2"
    }
  calendar_period = "YEAR"
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetCustomPeriodUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = []
    labels  = {
      label1 = "bar2"
    }
    custom_period {
      start_date {
      year = 2022
      month = 1
      day = 1
      }
      end_date {
      year = 2023
      month = 12
      day = 31
      }
    }
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
  }
}
`, context)
}

func TestBillingBudgetStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]string
		Meta       interface{}
	}{
		"shorten long name": {
			Attributes: map[string]interface{}{
				"name": "billingAccounts/000000-111111-222222/budgets/9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Expected: map[string]string{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Meta: &transport_tpg.Config{},
		},
		"short name stays": {
			Attributes: map[string]interface{}{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Expected: map[string]string{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Meta: &transport_tpg.Config{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := billingbudgets.ResourceBillingBudgetsBudgetUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

			if err != nil {
				t.Error(err)
			}

			for _, expectedName := range tc.Expected {
				if actual["name"] != expectedName {
					t.Errorf("expected: name -> %#v\n got: name -> %#v\n in: %#v",
						expectedName, actual["name"], actual)
				}
			}
		})
	}
}

func TestAccBillingBudget_budgetFilterProjectsOrdering(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org":                  envvar.GetTestOrgFromEnv(t),
		"billing_acct":         envvar.GetTestMasterBillingAccountFromEnv(t),
		"project_billing_acct": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix_1":      acctest.RandString(t, 10),
		"random_suffix_2":      acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_budgetFilterProjectsOrdering1(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:             testAccBillingBudget_budgetFilterProjectsOrdering2(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_budgetFilterProjectsOrdering1(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

resource "google_project" "project1" {
  project_id      = "tf-test-%{random_suffix_1}"
  name            = "tf-test-%{random_suffix_1}"
  org_id          = "%{org}"
  billing_account = "%{project_billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_project" "project2" {
  project_id      = "tf-test-%{random_suffix_2}"
  name            = "tf-test-%{random_suffix_2}"
  org_id          = "%{org}"
  billing_account = "%{project_billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget"

  budget_filter {
    projects = [
      "projects/${google_project.project1.number}",
      "projects/${google_project.project2.number}",
    ]
  }

  amount {
    last_period_amount = true
  }

  threshold_rules {
    threshold_percent =  10.0
  }
}

`, context)
}

func testAccBillingBudget_budgetFilterProjectsOrdering2(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

resource "google_project" "project1" {
  project_id      = "tf-test-%{random_suffix_1}"
  name            = "tf-test-%{random_suffix_1}"
  org_id          = "%{org}"
  billing_account = "%{project_billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_project" "project2" {
  project_id      = "tf-test-%{random_suffix_2}"
  name            = "tf-test-%{random_suffix_2}"
  org_id          = "%{org}"
  billing_account = "%{project_billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget"

  budget_filter {
    projects = [
      "projects/${google_project.project2.number}",
      "projects/${google_project.project1.number}",
    ]
  }

  amount {
    last_period_amount = true
  }

  threshold_rules {
    threshold_percent =  10.0
  }
}

`, context)
}

// A capped budget needs spendCap.inputState in every update mask, and the API
// never returns it. This changes only the amount, then removes the cap while
// changing the amount again, and checks import after each step.
//
// The API allows one spend cap per project and service. Each spend cap test
// creates its own project and also uses its own service, as a second
// safeguard against overlapping caps. The sample uses C7E2-9256-1C43 (Gemini
// Enterprise Agent Platform), this test AEFD-7695-64FA (Gemini API) and
// TestAccBillingBudget_spendCapAdd 152E-C115-5142 (Cloud Run).
func TestAccBillingBudget_spendCapUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		// The budget's account must bill the filtered project. For regular budgets,
		// a parent of the project's billing account also seems to work
		// (TestAccBillingBudget_budgetFilterProjectsOrdering), but that is not
		// documented. So the budget and the project share BILLING_ACCT, and these
		// tests do not rely on a parent account.
		"org":           envvar.GetTestOrgFromEnv(t),
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"service":       "AEFD-7695-64FA",
		"units":         "100",
	}
	contextUpdate := map[string]interface{}{
		"org":           context["org"],
		"billing_acct":  context["billing_acct"],
		"random_suffix": context["random_suffix"],
		"service":       context["service"],
		"units":         "200",
	}
	contextRemoved := map[string]interface{}{
		"org":           context["org"],
		"billing_acct":  context["billing_acct"],
		"random_suffix": context["random_suffix"],
		"service":       context["service"],
		"units":         "300",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_spendCap(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_spendCap(contextUpdate),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_billing_budget.budget", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_spendCapRemoved(contextRemoved),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_billing_budget.budget", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_billing_budget.budget", "spend_cap.#", "0"),
					resource.TestCheckResourceAttr("google_billing_budget.budget", "amount.0.specified_amount.0.units", "300"),
				),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// The API refuses to turn a regular budget into a spend cap budget, so adding
// a spend_cap block must replace the budget.
func TestAccBillingBudget_spendCapAdd(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		// See TestAccBillingBudget_spendCapUpdate.
		"org":           envvar.GetTestOrgFromEnv(t),
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"service":       "152E-C115-5142",
		"units":         "100",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBillingBudgetsBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_spendCapRemoved(context),
				Check:  resource.TestCheckResourceAttr("google_billing_budget.budget", "spend_cap.#", "0"),
			},
			{
				Config: testAccBillingBudget_spendCap(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_billing_budget.budget", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttr("google_billing_budget.budget", "spend_cap.0.input_state", "CONFIGURED"),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_spendCap(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org}"
  billing_account = "%{billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Spend Cap Budget%{random_suffix}"

  budget_filter {
    projects               = ["projects/${google_project.project.number}"]
    services               = ["services/%{service}"]
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
    calendar_period        = "MONTH"
  }

  amount {
    specified_amount {
      # No currency_code, so the budget uses the billing account's currency.
      units = "%{units}"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
    spend_basis       = "CURRENT_SPEND"
  }
  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "CURRENT_SPEND"
  }

  all_updates_rule {
    monitoring_notification_channels = []
    enable_project_level_recipients  = true
  }

  spend_cap {
    input_state = "CONFIGURED"
  }
}
`, context)
}

func testAccBillingBudget_spendCapRemoved(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org}"
  billing_account = "%{billing_acct}"
  deletion_policy = "DELETE"
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Spend Cap Budget%{random_suffix}"

  budget_filter {
    projects               = ["projects/${google_project.project.number}"]
    services               = ["services/%{service}"]
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
    calendar_period        = "MONTH"
  }

  amount {
    specified_amount {
      # No currency_code, so the budget uses the billing account's currency.
      units = "%{units}"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
    spend_basis       = "CURRENT_SPEND"
  }
  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "CURRENT_SPEND"
  }

  all_updates_rule {
    monitoring_notification_channels = []
    enable_project_level_recipients  = true
  }
}
`, context)
}
