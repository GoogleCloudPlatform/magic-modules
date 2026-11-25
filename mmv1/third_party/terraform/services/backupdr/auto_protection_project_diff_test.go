package backupdr

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func resetBackupDRProjectCacheForTest(t *testing.T) {
	t.Helper()
	backupDRProjectIDByNumberMu.Lock()
	saved := backupDRProjectIDByNumber
	backupDRProjectIDByNumber = map[string]string{}
	backupDRProjectIDByNumberMu.Unlock()

	t.Cleanup(func() {
		backupDRProjectIDByNumberMu.Lock()
		backupDRProjectIDByNumber = saved
		backupDRProjectIDByNumberMu.Unlock()
	})
}

func TestBackupDRAutoProtectionProjectDiffSuppress(t *testing.T) {
	resetBackupDRProjectCacheForTest(t)

	recordBackupDRProjectIDAndNumber("my-project-a", "111111111111")
	recordBackupDRProjectIDAndNumber("my-project-b", "222222222222")

	cases := map[string]struct {
		old      string
		new      string
		expected bool
	}{
		"exact match scope project ID": {
			old:      "projects/my-project-a",
			new:      "projects/my-project-a",
			expected: true,
		},
		"exact match scope project number": {
			old:      "projects/111111111111",
			new:      "projects/111111111111",
			expected: true,
		},
		"same project scope number to ID": {
			old:      "projects/111111111111",
			new:      "projects/my-project-a",
			expected: true,
		},
		"same project scope ID to number": {
			old:      "projects/my-project-a",
			new:      "projects/111111111111",
			expected: true,
		},
		"different project scope number to other ID": {
			old:      "projects/111111111111",
			new:      "projects/my-project-b",
			expected: false,
		},
		"different project scope number to uncached ID": {
			old:      "projects/111111111111",
			new:      "projects/uncached-project",
			expected: false,
		},
		"different project scope uncached number to ID": {
			old:      "projects/333333333333",
			new:      "projects/my-project-a",
			expected: false,
		},
		"different project scope number to different number": {
			old:      "projects/111111111111",
			new:      "projects/222222222222",
			expected: false,
		},
		"different project scope ID to different ID": {
			old:      "projects/my-project-a",
			new:      "projects/my-project-b",
			expected: false,
		},
		"same project backup_plan number to ID": {
			old:      "projects/111111111111/locations/us-central1/backupPlans/bp-1",
			new:      "projects/my-project-a/locations/us-central1/backupPlans/bp-1",
			expected: true,
		},
		"same project backup_plan ID to number": {
			old:      "projects/my-project-a/locations/us-central1/backupPlans/bp-1",
			new:      "projects/111111111111/locations/us-central1/backupPlans/bp-1",
			expected: true,
		},
		"same project backup_plan different plan name": {
			old:      "projects/111111111111/locations/us-central1/backupPlans/bp-1",
			new:      "projects/my-project-a/locations/us-central1/backupPlans/bp-2",
			expected: false,
		},
		"different project backup_plan same plan name": {
			old:      "projects/111111111111/locations/us-central1/backupPlans/bp-1",
			new:      "projects/my-project-b/locations/us-central1/backupPlans/bp-1",
			expected: false,
		},
		"different numeric project backup_plan same plan name": {
			old:      "projects/111111111111/locations/us-central1/backupPlans/bp-1",
			new:      "projects/222222222222/locations/us-central1/backupPlans/bp-1",
			expected: false,
		},
		"scope vs resource path mismatch": {
			old:      "projects/111111111111",
			new:      "projects/my-project-a/locations/us-central1/backupPlans/bp-1",
			expected: false,
		},
		"non-project scope": {
			old:      "folders/111111111111",
			new:      "folders/222222222222",
			expected: false,
		},
		"malformed empty project": {
			old:      "projects/",
			new:      "projects/my-project-a",
			expected: false,
		},
		"malformed double slash": {
			old:      "projects//locations/us-central1/backupPlans/bp-1",
			new:      "projects/my-project-a/locations/us-central1/backupPlans/bp-1",
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := backupDRAutoProtectionProjectDiffSuppress("scope", tc.old, tc.new, nil)
			if got != tc.expected {
				t.Fatalf("backupDRAutoProtectionProjectDiffSuppress(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.expected)
			}
		})
	}
}

func TestRecordBackupDRAutoProtectionPolicyBindingScopeProjectNumber(t *testing.T) {
	t.Run("resolves cross-project scope via resolver and memoizes", func(t *testing.T) {
		resetBackupDRProjectCacheForTest(t)

		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"project": "mgmt-project",
			"scope":   "projects/wrong-unverified-project",
		})
		res := map[string]interface{}{
			"name":  "projects/999999999999/locations/us-central1/autoProtectionPolicies/p1/bindings/b1",
			"scope": "projects/123456789012",
		}

		calls := 0
		resolver := func(projectNumber string) (string, error) {
			calls++
			if projectNumber != "123456789012" {
				t.Fatalf("unexpected projectNumber = %q", projectNumber)
			}
			return "target-project-a", nil
		}

		recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(d, &transport_tpg.Config{}, res, resolver)
		recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(d, &transport_tpg.Config{}, res, resolver)

		if calls != 1 {
			t.Fatalf("expected resolver to be called once due to caching, got %d", calls)
		}
		if !backupDRAutoProtectionProjectDiffSuppress("scope", "projects/123456789012", "projects/target-project-a", d) {
			t.Fatalf("expected verified scope project ID to be suppressed")
		}
		if backupDRAutoProtectionProjectDiffSuppress("scope", "projects/123456789012", "projects/wrong-unverified-project", d) {
			t.Fatalf("expected unverified config scope not to be suppressed")
		}
	})

	t.Run("does not trust config scope when cross-project resolver fails", func(t *testing.T) {
		resetBackupDRProjectCacheForTest(t)

		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"project": "mgmt-project",
			"scope":   "projects/some-other-project",
		})
		res := map[string]interface{}{
			"name":  "projects/999999999999/locations/us-central1/autoProtectionPolicies/p1/bindings/b1",
			"scope": "projects/123456789012",
		}

		resolver := func(string) (string, error) {
			return "", errors.New("permission denied")
		}

		recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(d, &transport_tpg.Config{}, res, resolver)

		if backupDRAutoProtectionProjectDiffSuppress("scope", "projects/123456789012", "projects/some-other-project", d) {
			t.Fatalf("expected failed cross-project lookup not to suppress diff")
		}
	})

	t.Run("uses API-verified same-project fallback when resolver fails", func(t *testing.T) {
		resetBackupDRProjectCacheForTest(t)

		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"project": "same-project",
			"scope":   "projects/same-project",
		})
		res := map[string]interface{}{
			"name":  "projects/123456789012/locations/us-central1/autoProtectionPolicies/p1/bindings/b1",
			"scope": "projects/123456789012",
		}

		resolver := func(string) (string, error) {
			return "", errors.New("crm unavailable")
		}

		recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(d, &transport_tpg.Config{}, res, resolver)

		if !backupDRAutoProtectionProjectDiffSuppress("scope", "projects/123456789012", "projects/same-project", d) {
			t.Fatalf("expected API-verified same-project fallback to suppress diff")
		}
		if backupDRAutoProtectionProjectDiffSuppress("scope", "projects/123456789012", "projects/different-project", d) {
			t.Fatalf("expected different project not to be suppressed")
		}
	})
}

func TestRecordBackupDRAutoProtectionPolicyProjectNumbers(t *testing.T) {
	t.Run("resolves backup_plan project number via resolver", func(t *testing.T) {
		resetBackupDRProjectCacheForTest(t)

		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"project": "123456789012",
		})
		res := map[string]interface{}{
			"name": "projects/123456789012/locations/us-central1/autoProtectionPolicies/p1",
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		resolver := func(projectNumber string) (string, error) {
			if projectNumber != "123456789012" {
				t.Fatalf("unexpected projectNumber = %q", projectNumber)
			}
			return "policy-project", nil
		}

		recordBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, &transport_tpg.Config{}, res, resolver)

		if !backupDRAutoProtectionProjectDiffSuppress(
			"backup_plan_details.0.backup_plan",
			"projects/123456789012/locations/us-central1/backupPlans/bp-1",
			"projects/policy-project/locations/us-central1/backupPlans/bp-1",
			d,
		) {
			t.Fatalf("expected same-project backup_plan diff to be suppressed")
		}
		if backupDRAutoProtectionProjectDiffSuppress(
			"backup_plan_details.0.backup_plan",
			"projects/123456789012/locations/us-central1/backupPlans/bp-1",
			"projects/other-project/locations/us-central1/backupPlans/bp-1",
			d,
		) {
			t.Fatalf("expected different-project backup_plan diff not to be suppressed")
		}
	})

	t.Run("uses API-verified same-project fallback when resolver fails", func(t *testing.T) {
		resetBackupDRProjectCacheForTest(t)

		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"project": "policy-project",
		})
		res := map[string]interface{}{
			"name": "projects/123456789012/locations/us-central1/autoProtectionPolicies/p1",
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		resolver := func(string) (string, error) {
			return "", errors.New("crm unavailable")
		}

		recordBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, &transport_tpg.Config{}, res, resolver)

		if !backupDRAutoProtectionProjectDiffSuppress(
			"backup_plan_details.0.backup_plan",
			"projects/123456789012/locations/us-central1/backupPlans/bp-1",
			"projects/policy-project/locations/us-central1/backupPlans/bp-1",
			d,
		) {
			t.Fatalf("expected fallback mapping to suppress same-project diff")
		}
	})
}
