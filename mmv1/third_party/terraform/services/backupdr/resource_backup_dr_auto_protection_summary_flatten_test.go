package backupdr

import (
	"reflect"
	"testing"
)

func TestBackupDRAutoProtectionSummaryZeroCountFlatteners(t *testing.T) {
	t.Run("policy_evaluation_summary converts omitted and empty counts to 0", func(t *testing.T) {
		got := flattenBackupDRAutoProtectionPolicyPolicyEvaluationSummary(map[string]interface{}{
			"matchingResourceCount":  "5",
			"protectedResourceCount": "",
		}, nil, nil)

		want := []interface{}{
			map[string]interface{}{
				"matching_resource_count":  "5",
				"protected_resource_count": "0",
				"failed_resource_count":    "0",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("policy_evaluation_summary preserves nil and empty object behavior", func(t *testing.T) {
		if got := flattenBackupDRAutoProtectionPolicyPolicyEvaluationSummary(nil, nil, nil); got != nil {
			t.Fatalf("nil input: got %#v, want nil", got)
		}
		if got := flattenBackupDRAutoProtectionPolicyPolicyEvaluationSummary(map[string]interface{}{}, nil, nil); got != nil {
			t.Fatalf("empty object input: got %#v, want nil", got)
		}
	})

	t.Run("binding_summary converts empty count to 0 and preserves nil/empty object behavior", func(t *testing.T) {
		got := flattenBackupDRAutoProtectionPolicyBindingSummary(map[string]interface{}{
			"projectCount": "",
		}, nil, nil)
		want := []interface{}{
			map[string]interface{}{
				"project_count": "0",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}

		gotNonZero := flattenBackupDRAutoProtectionPolicyBindingSummary(map[string]interface{}{
			"projectCount": "3",
		}, nil, nil)
		wantNonZero := []interface{}{
			map[string]interface{}{
				"project_count": "3",
			},
		}
		if !reflect.DeepEqual(gotNonZero, wantNonZero) {
			t.Fatalf("got %#v, want %#v", gotNonZero, wantNonZero)
		}

		if got := flattenBackupDRAutoProtectionPolicyBindingSummary(nil, nil, nil); got != nil {
			t.Fatalf("nil input: got %#v, want nil", got)
		}
		if got := flattenBackupDRAutoProtectionPolicyBindingSummary(map[string]interface{}{}, nil, nil); got != nil {
			t.Fatalf("empty object input: got %#v, want nil", got)
		}
		if got := flattenBackupDRAutoProtectionPolicyBindingSummaryProjectCount(nil, nil, nil); got != "0" {
			t.Fatalf("nil projectCount: got %#v, want \"0\"", got)
		}
	})

	t.Run("binding_evaluation_summary converts omitted and empty counts to 0", func(t *testing.T) {
		got := flattenBackupDRAutoProtectionPolicyBindingBindingEvaluationSummary(map[string]interface{}{
			"matchingResourceCount": "2",
			"failedResourceCount":   "",
		}, nil, nil)

		want := []interface{}{
			map[string]interface{}{
				"matching_resource_count":  "2",
				"protected_resource_count": "0",
				"failed_resource_count":    "0",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}

		if got := flattenBackupDRAutoProtectionPolicyBindingBindingEvaluationSummary(nil, nil, nil); got != nil {
			t.Fatalf("nil input: got %#v, want nil", got)
		}
		if got := flattenBackupDRAutoProtectionPolicyBindingBindingEvaluationSummary(map[string]interface{}{}, nil, nil); got != nil {
			t.Fatalf("empty object input: got %#v, want nil", got)
		}
	})
}
