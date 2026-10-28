package billingbudgets

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitBillingBudget_spendCapAddedCustomizeDiff(t *testing.T) {
	t.Parallel()

	configured := []interface{}{
		map[string]interface{}{"input_state": "CONFIGURED"},
	}
	lifted := []interface{}{
		map[string]interface{}{"input_state": "AWAITING_NEXT_PERIOD"},
	}

	cases := map[string]struct {
		id       string
		before   interface{}
		after    interface{}
		forceNew bool
	}{
		"new resource with spend_cap": {
			id:       "",
			before:   []interface{}{},
			after:    configured,
			forceNew: false,
		},
		"spend_cap added": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   []interface{}{},
			after:    configured,
			forceNew: true,
		},
		"spend_cap added, no prior value": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   nil,
			after:    configured,
			forceNew: true,
		},
		"spend_cap removed": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   configured,
			after:    []interface{}{},
			forceNew: false,
		},
		"spend_cap kept": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   configured,
			after:    configured,
			forceNew: false,
		},
		"spend_cap kept, input_state changed": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   configured,
			after:    lifted,
			forceNew: false,
		},
		"no spend_cap": {
			id:       "billingAccounts/000000-000000-000000/budgets/budget",
			before:   []interface{}{},
			after:    []interface{}{},
			forceNew: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			d := &tpgresource.ResourceDiffMock{
				Before: map[string]interface{}{"spend_cap": tc.before},
				After:  map[string]interface{}{"spend_cap": tc.after},
			}
			if err := billingBudgetSpendCapAddedCustomizeDiffFunc(d, tc.id); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if d.IsForceNew != tc.forceNew {
				t.Errorf("ForceNew = %t, want %t", d.IsForceNew, tc.forceNew)
			}
		})
	}
}
