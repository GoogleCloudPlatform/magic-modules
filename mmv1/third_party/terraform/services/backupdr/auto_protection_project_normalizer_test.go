package backupdr

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestBackupDRAutoProtectionPolicyBindingScopeForceNewIfChanged(t *testing.T) {
	baseState := func(scope string) *terraform.InstanceState {
		return &terraform.InstanceState{
			ID: "projects/mgmt-project/locations/us-central1/autoProtectionPolicies/p1/bindings/b1",
			Attributes: map[string]string{
				"project":                           "mgmt-project",
				"location":                          "us-central1",
				"auto_protection_policy_id":         "p1",
				"auto_protection_policy_binding_id": "b1",
				"scope":                             scope,
				"wait_for_delete":                   "false",
				"deletion_policy":                   "DELETE",
			},
		}
	}

	baseConfig := func(scope string, waitForDelete bool) *terraform.ResourceConfig {
		return terraform.NewResourceConfigRaw(map[string]interface{}{
			"project":                           "mgmt-project",
			"location":                          "us-central1",
			"auto_protection_policy_id":         "p1",
			"auto_protection_policy_binding_id": "b1",
			"scope":                             scope,
			"wait_for_delete":                   waitForDelete,
			"deletion_policy":                   "DELETE",
		})
	}

	diffWithResolvers := func(
		t *testing.T,
		state *terraform.InstanceState,
		cfg *terraform.ResourceConfig,
		resolveIDFromNumber func(string) (string, error),
		resolveNumberFromID func(string) (string, error),
	) (*terraform.InstanceDiff, error) {
		t.Helper()
		res := ResourceBackupDRAutoProtectionPolicyBinding()
		res.CustomizeDiff = func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
			return backupDRAutoProtectionPolicyBindingScopeForceNewIfChangedWithResolvers(
				diff,
				resolveIDFromNumber,
				resolveNumberFromID,
			)
		}
		return res.Diff(context.Background(), state, cfg, nil)
	}

	t.Run("upgrade path: numeric state to matching project ID updates in-place without ForceNew", func(t *testing.T) {
		numCalls := 0
		diff, err := diffWithResolvers(
			t,
			baseState("projects/123456789012"),
			baseConfig("projects/target-project-a", false),
			func(projectNumber string) (string, error) {
				numCalls++
				if projectNumber != "123456789012" {
					t.Fatalf("unexpected projectNumber = %q", projectNumber)
				}
				return "target-project-a", nil
			},
			func(string) (string, error) {
				t.Fatal("resolveNumberFromID should not be called")
				return "", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if numCalls != 1 {
			t.Fatalf("expected 1 resolveIDFromNumber call, got %d", numCalls)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place state update (RequiresNew == false), got RequiresNew == true: %#v", diff.Attributes["scope"])
		}
		if attr := diff.Attributes["scope"]; attr == nil || attr.Old != "projects/123456789012" || attr.New != "projects/target-project-a" {
			t.Fatalf("unexpected scope diff attr: %#v", attr)
		}
	})

	t.Run("numeric state to different project ID forces replacement", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			baseState("projects/123456789012"),
			baseConfig("projects/other-project", false),
			func(string) (string, error) {
				return "target-project-a", nil
			},
			nil,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !diff.RequiresNew() {
			t.Fatal("expected ForceNew replacement when numeric state and new project ID refer to different projects")
		}
		if attr := diff.Attributes["scope"]; attr == nil || !attr.RequiresNew {
			t.Fatalf("expected scope.RequiresNew == true, got %#v", attr)
		}
	})

	t.Run("project ID state to matching numeric project updates in-place without ForceNew", func(t *testing.T) {
		idCalls := 0
		diff, err := diffWithResolvers(
			t,
			baseState("projects/target-project-a"),
			baseConfig("projects/123456789012", false),
			func(string) (string, error) {
				t.Fatal("resolveIDFromNumber should not be called")
				return "", nil
			},
			func(projectID string) (string, error) {
				idCalls++
				if projectID != "target-project-a" {
					t.Fatalf("unexpected projectID = %q", projectID)
				}
				return "123456789012", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if idCalls != 1 {
			t.Fatalf("expected 1 resolveNumberFromID call, got %d", idCalls)
		}
		if diff.RequiresNew() {
			t.Fatal("expected in-place update when switching project ID to its own project number")
		}
	})

	t.Run("project ID state to different numeric project resolves existing state project and forces replacement", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			baseState("projects/target-project-a"),
			baseConfig("projects/999999999999", false),
			func(string) (string, error) {
				t.Fatal("should resolve old state project ID, not unverified new numeric project")
				return "", nil
			},
			func(projectID string) (string, error) {
				if projectID != "target-project-a" {
					t.Fatalf("unexpected projectID = %q", projectID)
				}
				return "123456789012", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !diff.RequiresNew() {
			t.Fatal("expected ForceNew replacement when switching project ID to a different project number")
		}
	})

	t.Run("different project IDs force replacement without CRM calls", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			baseState("projects/target-project-a"),
			baseConfig("projects/target-project-b", false),
			func(string) (string, error) {
				t.Fatal("resolver should not be called")
				return "", nil
			},
			func(string) (string, error) {
				t.Fatal("resolver should not be called")
				return "", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !diff.RequiresNew() {
			t.Fatal("expected ForceNew replacement when changing between different project IDs")
		}
	})

	t.Run("different numeric project numbers force replacement without CRM calls", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			baseState("projects/111111111111"),
			baseConfig("projects/222222222222", false),
			func(string) (string, error) {
				t.Fatal("resolver should not be called")
				return "", nil
			},
			func(string) (string, error) {
				t.Fatal("resolver should not be called")
				return "", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !diff.RequiresNew() {
			t.Fatal("expected ForceNew replacement when changing between different project numbers")
		}
	})

	t.Run("CRM lookup error fails diff instead of forcing destructive replacement", func(t *testing.T) {
		_, err := diffWithResolvers(
			t,
			baseState("projects/123456789012"),
			baseConfig("projects/target-project-a", false),
			func(string) (string, error) {
				return "", errors.New("403 permission denied")
			},
			nil,
		)
		if err == nil {
			t.Fatal("expected error when CRM lookup fails during scope diff")
		}
		if !strings.Contains(err.Error(), "projects/123456789012") {
			t.Fatalf("expected actionable error message, got: %v", err)
		}
	})

	t.Run("unchanged scope with other field update skips CRM and does not ForceNew", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			baseState("projects/target-project-a"),
			baseConfig("projects/target-project-a", true),
			func(string) (string, error) {
				t.Fatal("resolver should not be called when scope is unchanged")
				return "", nil
			},
			func(string) (string, error) {
				t.Fatal("resolver should not be called when scope is unchanged")
				return "", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatal("expected in-place update for wait_for_delete change")
		}
		if _, hasScope := diff.Attributes["scope"]; hasScope {
			t.Fatalf("expected no diff on scope, got %#v", diff.Attributes["scope"])
		}
	})

	t.Run("new resource skips CRM and does not error", func(t *testing.T) {
		diff, err := diffWithResolvers(
			t,
			nil,
			baseConfig("projects/target-project-a", false),
			func(string) (string, error) {
				t.Fatal("resolver should not be called on create")
				return "", nil
			},
			func(string) (string, error) {
				t.Fatal("resolver should not be called on create")
				return "", nil
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff == nil {
			t.Fatal("expected non-nil create diff")
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

func TestBackupDRAutoProtectionSchemaInvariants(t *testing.T) {
	bindingSchema := ResourceBackupDRAutoProtectionPolicyBinding().Schema
	if bindingSchema["scope"].DiffSuppressFunc != nil {
		t.Error("expected AutoProtectionPolicyBinding.scope not to use DiffSuppressFunc")
	}
	if bindingSchema["scope"].ForceNew {
		t.Error("expected AutoProtectionPolicyBinding.scope ForceNew to be managed via CustomizeDiff")
	}
	for _, field := range []string{"project", "location", "auto_protection_policy_id", "auto_protection_policy_binding_id"} {
		if !bindingSchema[field].ForceNew {
			t.Errorf("expected AutoProtectionPolicyBinding.%s to remain ForceNew: true", field)
		}
	}

	policyBPD := ResourceBackupDRAutoProtectionPolicy().Schema["backup_plan_details"].Elem.(*schema.Resource).Schema["backup_plan"]
	if policyBPD.DiffSuppressFunc != nil {
		t.Error("expected AutoProtectionPolicy.backup_plan_details.backup_plan not to use DiffSuppressFunc")
	}
}
