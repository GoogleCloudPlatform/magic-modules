package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBillingBudget_billingBudgetCurrencycode(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckBillingBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetCurrencycode(context),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccBillingBudget_billingBudgetUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckBillingBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetUpdate(context),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdate2(context),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccBillingBudget_billingBudgetCurrencycode(context map[string]interface{}) string {
	return Nprintf(`
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

func testAccBillingBudget_billingBudgetUpdate(context map[string]interface{}) string {
	return Nprintf(`
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
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
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

func testAccBillingBudget_billingBudgetUpdate2(context map[string]interface{}) string {
	return Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget%{random_suffix}"

  budget_filter {
    credit_types_treatment = "INCLUDE_SPECIFIED_CREDITS"
  }

  amount {
    specified_amount {
      units         = "1000000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.3
    spend_basis       = "FORECASTED_SPEND"
  }
}
`, context)
}
