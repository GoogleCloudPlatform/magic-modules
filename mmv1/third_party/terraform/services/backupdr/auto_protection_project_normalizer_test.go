package backupdr

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestNormalizeBackupDRAutoProtectionPolicyBindingScope(t *testing.T) {
	t.Run("rewrites verified project number back to configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"project": "mgmt-project",
			"scope":   "projects/target-project-a",
		})
		res := map[string]interface{}{
			"name":  "projects/mgmt-project/locations/us-central1/autoProtectionPolicies/p1/bindings/b1",
			"scope": "projects/123456789012",
		}

		calls := 0
		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(projectNumber string) (string, error) {
			calls++
			if projectNumber != "123456789012" {
				t.Fatalf("unexpected projectNumber = %q", projectNumber)
			}
			return "target-project-a", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 resolver call, got %d", calls)
		}
		if got := res["scope"]; got != "projects/target-project-a" {
			t.Fatalf("scope = %q, want %q", got, "projects/target-project-a")
		}
	})

	t.Run("does not rewrite when resolved project ID differs from configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"scope": "projects/target-project-b",
		})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(string) (string, error) {
			return "target-project-a", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := res["scope"]; got != "projects/123456789012" {
			t.Fatalf("scope = %q, want %q", got, "projects/123456789012")
		}
	})

	t.Run("skips resolver when configured scope already matches numeric project number", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"scope": "projects/123456789012",
		})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		called := false
		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(string) (string, error) {
			called = true
			return "", errors.New("should not be called")
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Fatal("expected resolver not to be called for numeric scope match")
		}
		if got := res["scope"]; got != "projects/123456789012" {
			t.Fatalf("scope = %q, want %q", got, "projects/123456789012")
		}
	})

	t.Run("skips resolver and keeps API value when numeric scopes differ", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"scope": "projects/999999999999",
		})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		called := false
		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(string) (string, error) {
			called = true
			return "", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Fatal("expected resolver not to be called when comparing two numeric scopes")
		}
		if got := res["scope"]; got != "projects/123456789012" {
			t.Fatalf("scope = %q, want %q", got, "projects/123456789012")
		}
	})

	t.Run("returns error when resolver fails for configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
			"scope": "projects/target-project-a",
		})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(string) (string, error) {
			return "", errors.New("403 permission denied")
		})
		if err == nil {
			t.Fatal("expected error when CRM lookup fails for configured project ID")
		}
		if !strings.Contains(err.Error(), "projects/123456789012") {
			t.Fatalf("expected actionable hint with numeric scope, got: %v", err)
		}
	})

	t.Run("resolves to project ID on import", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(projectNumber string) (string, error) {
			if projectNumber != "123456789012" {
				t.Fatalf("unexpected projectNumber = %q", projectNumber)
			}
			return "target-project-a", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := res["scope"]; got != "projects/target-project-a" {
			t.Fatalf("scope = %q, want %q", got, "projects/target-project-a")
		}
	})

	t.Run("falls back to numeric scope on import when resolver fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{})
		res := map[string]interface{}{
			"scope": "projects/123456789012",
		}

		err := normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(string) (string, error) {
			return "", errors.New("403 permission denied")
		})
		if err != nil {
			t.Fatalf("expected best-effort import not to fail on CRM error, got: %v", err)
		}
		if got := res["scope"]; got != "projects/123456789012" {
			t.Fatalf("scope = %q, want %q", got, "projects/123456789012")
		}
	})
}

func TestNormalizeBackupDRAutoProtectionPolicyProjectNumbers(t *testing.T) {
	t.Run("rewrites verified backup_plan project number back to configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"backup_plan_details": []interface{}{
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Instance",
					"backup_plan":   "projects/policy-project/locations/us-central1/backupPlans/bp-1",
				},
			},
		})
		res := map[string]interface{}{
			"name": "projects/policy-project/locations/us-central1/autoProtectionPolicies/p1",
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(projectNumber string) (string, error) {
			if projectNumber != "123456789012" {
				t.Fatalf("unexpected projectNumber = %q", projectNumber)
			}
			return "policy-project", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res["backupPlanDetails"].([]interface{})[0].(map[string]interface{})["backupPlan"]
		want := "projects/policy-project/locations/us-central1/backupPlans/bp-1"
		if got != want {
			t.Fatalf("backupPlan = %q, want %q", got, want)
		}
	})

	t.Run("does not rewrite when resolved project ID differs from configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"backup_plan_details": []interface{}{
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Instance",
					"backup_plan":   "projects/other-project/locations/us-central1/backupPlans/bp-1",
				},
			},
		})
		res := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(string) (string, error) {
			return "policy-project", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res["backupPlanDetails"].([]interface{})[0].(map[string]interface{})["backupPlan"]
		want := "projects/123456789012/locations/us-central1/backupPlans/bp-1"
		if got != want {
			t.Fatalf("backupPlan = %q, want %q", got, want)
		}
	})

	t.Run("skips resolver when configured backup_plan already matches numeric project number", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"backup_plan_details": []interface{}{
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Instance",
					"backup_plan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		})
		res := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		called := false
		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(string) (string, error) {
			called = true
			return "", errors.New("should not be called")
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Fatal("expected resolver not to be called for numeric backup_plan match")
		}
	})

	t.Run("returns error when resolver fails for configured project ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"backup_plan_details": []interface{}{
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Instance",
					"backup_plan":   "projects/policy-project/locations/us-central1/backupPlans/bp-1",
				},
			},
		})
		res := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(string) (string, error) {
			return "", errors.New("403 permission denied")
		})
		if err == nil {
			t.Fatal("expected error when CRM lookup fails for configured project ID")
		}
		if !strings.Contains(err.Error(), "projects/123456789012/locations/us-central1/backupPlans/bp-1") {
			t.Fatalf("expected actionable hint with numeric backup_plan, got: %v", err)
		}
	})

	t.Run("resolves to project ID on import and falls back on CRM error", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{})
		res := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}

		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(string) (string, error) {
			return "policy-project", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res["backupPlanDetails"].([]interface{})[0].(map[string]interface{})["backupPlan"]
		want := "projects/policy-project/locations/us-central1/backupPlans/bp-1"
		if got != want {
			t.Fatalf("backupPlan = %q, want %q", got, want)
		}

		resErr := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
			},
		}
		err = normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, resErr, func(string) (string, error) {
			return "", errors.New("403 permission denied")
		})
		if err != nil {
			t.Fatalf("expected best-effort import not to fail on CRM error, got: %v", err)
		}
		gotErr := resErr["backupPlanDetails"].([]interface{})[0].(map[string]interface{})["backupPlan"]
		wantErr := "projects/123456789012/locations/us-central1/backupPlans/bp-1"
		if gotErr != wantErr {
			t.Fatalf("backupPlan = %q, want %q", gotErr, wantErr)
		}
	})

	t.Run("memoizes multiple backup_plan entries within a single Read call only", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicy().Schema, map[string]interface{}{
			"backup_plan_details": []interface{}{
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Instance",
					"backup_plan":   "projects/policy-project/locations/us-central1/backupPlans/bp-1",
				},
				map[string]interface{}{
					"resource_type": "compute.googleapis.com/Disk",
					"backup_plan":   "projects/policy-project/locations/us-central1/backupPlans/bp-2",
				},
			},
		})
		res := map[string]interface{}{
			"backupPlanDetails": []interface{}{
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Instance",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-1",
				},
				map[string]interface{}{
					"resourceType": "compute.googleapis.com/Disk",
					"backupPlan":   "projects/123456789012/locations/us-central1/backupPlans/bp-2",
				},
			},
		}

		calls := 0
		err := normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(string) (string, error) {
			calls++
			return "policy-project", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 resolver call within single Read, got %d", calls)
		}
	})
}

func TestBackupDRResourceManagerBillingProject(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceBackupDRAutoProtectionPolicyBinding().Schema, map[string]interface{}{
		"project": "resource-project",
	})

	if got := backupDRResourceManagerBillingProject(d, &transport_tpg.Config{}); got != "resource-project" {
		t.Fatalf("billing project = %q, want %q", got, "resource-project")
	}

	if got := backupDRResourceManagerBillingProject(d, &transport_tpg.Config{
		UserProjectOverride: true,
		BillingProject:      "explicit-billing-project",
	}); got != "explicit-billing-project" {
		t.Fatalf("billing project = %q, want %q", got, "explicit-billing-project")
	}
}

func TestBackupDRAutoProtectionSchemaHasNoProjectNumberDiffSuppress(t *testing.T) {
	bindingScope := ResourceBackupDRAutoProtectionPolicyBinding().Schema["scope"]
	if bindingScope.DiffSuppressFunc != nil {
		t.Error("expected AutoProtectionPolicyBinding.scope not to use DiffSuppressFunc")
	}

	policyBPD := ResourceBackupDRAutoProtectionPolicy().Schema["backup_plan_details"].Elem.(*schema.Resource).Schema["backup_plan"]
	if policyBPD.DiffSuppressFunc != nil {
		t.Error("expected AutoProtectionPolicy.backup_plan_details.backup_plan not to use DiffSuppressFunc")
	}
}
