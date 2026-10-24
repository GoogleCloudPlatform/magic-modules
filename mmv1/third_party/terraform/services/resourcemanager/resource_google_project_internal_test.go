package resourcemanager

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestPopulateGoogleProjectResourceData_DoesNotCopyAllLabelsToTerraformLabelsWhenUnset(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceGoogleProject().Schema, map[string]interface{}{
		"project_id": "example-project",
		"name":       "Example Project",
	})

	project := &cloudresourcemanager.Project{
		ProjectNumber:  123456789,
		Name:           "Example Project",
		LifecycleState: "ACTIVE",
		Labels: map[string]string{
			"firebase":     "enabled",
			"earth-engine": "",
		},
	}

	if err := populateGoogleProjectResourceData(d, project, "example-project", &transport_tpg.Config{}); err != nil {
		t.Fatalf("populateGoogleProjectResourceData() returned error: %v", err)
	}

	if got := d.Get("labels").(map[string]interface{}); len(got) != 0 {
		t.Fatalf("expected labels to remain empty when unset in config, got %#v", got)
	}

	if got := d.Get("terraform_labels").(map[string]interface{}); len(got) != 0 {
		t.Fatalf("expected terraform_labels to remain empty when unset in config, got %#v", got)
	}

	if got := d.Get("effective_labels").(map[string]interface{}); !reflect.DeepEqual(got, map[string]interface{}{
		"firebase":     "enabled",
		"earth-engine": "",
	}) {
		t.Fatalf("expected effective_labels to contain all project labels, got %#v", got)
	}

	if got := d.Get("lifecycle_state").(string); got != "ACTIVE" {
		t.Fatalf("expected lifecycle_state to be ACTIVE, got %q", got)
	}

	if got := d.Get("recover_soft_deleted").(bool); got != false {
		t.Fatalf("expected recover_soft_deleted default to be false, got %t", got)
	}
}

func TestValidateSoftDeletedProjectPreflight(t *testing.T) {
	pid := "test-preflight-project"
	softDeletedOrgProj := &cloudresourcemanager.Project{
		ProjectId:      pid,
		LifecycleState: "DELETE_REQUESTED",
		Parent: &cloudresourcemanager.ResourceId{
			Type: "organization",
			Id:   "12345",
		},
	}

	cases := []struct {
		name               string
		recoverSoftDeleted bool
		configuredOrgId    string
		existing           *cloudresourcemanager.Project
		wantErrSubstrs     []string
	}{
		{
			name:               "nil_project",
			recoverSoftDeleted: true,
			configuredOrgId:    "12345",
			existing:           nil,
		},
		{
			name:               "active_project",
			recoverSoftDeleted: false,
			configuredOrgId:    "12345",
			existing: &cloudresourcemanager.Project{
				ProjectId:      pid,
				LifecycleState: "ACTIVE",
			},
		},
		{
			name:               "delete_in_progress",
			recoverSoftDeleted: true,
			configuredOrgId:    "12345",
			existing: &cloudresourcemanager.Project{
				ProjectId:      pid,
				LifecycleState: "DELETE_IN_PROGRESS",
			},
			wantErrSubstrs: []string{"DELETE_IN_PROGRESS"},
		},
		{
			name:               "delete_requested_without_recover_flag",
			recoverSoftDeleted: false,
			configuredOrgId:    "12345",
			existing:           softDeletedOrgProj,
			wantErrSubstrs:     []string{"recover_soft_deleted = true", "and import it into Terraform state"},
		},
		{
			name:               "delete_requested_matching_org",
			recoverSoftDeleted: true,
			configuredOrgId:    "12345",
			existing:           softDeletedOrgProj,
		},
		{
			name:               "delete_requested_conflicting_org",
			recoverSoftDeleted: true,
			configuredOrgId:    "99999",
			existing:           softDeletedOrgProj,
			wantErrSubstrs:     []string{"conflicts with configured org_id"},
		},
		{
			name:               "delete_requested_folder_parent",
			recoverSoftDeleted: true,
			configuredOrgId:    "12345",
			existing: &cloudresourcemanager.Project{
				ProjectId:      pid,
				LifecycleState: "DELETE_REQUESTED",
				Parent: &cloudresourcemanager.ResourceId{
					Type: "folder",
					Id:   "45678",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSoftDeletedProjectPreflight(pid, tc.recoverSoftDeleted, tc.configuredOrgId, tc.existing)
			if len(tc.wantErrSubstrs) == 0 {
				if err != nil {
					t.Fatalf("validateSoftDeletedProjectPreflight() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validateSoftDeletedProjectPreflight() expected error containing %v, got nil", tc.wantErrSubstrs)
			}
			for _, sub := range tc.wantErrSubstrs {
				if !strings.Contains(err.Error(), sub) {
					t.Fatalf("validateSoftDeletedProjectPreflight() error = %q, want substring %q", err.Error(), sub)
				}
			}
		})
	}
}

func TestParseFolderId(t *testing.T) {
	cases := []struct {
		name  string
		input interface{}
		want  string
	}{
		{name: "nil_value", input: nil, want: ""},
		{name: "non_string_int", input: 12345, want: ""},
		{name: "prefixed_folder_id", input: "folders/12345", want: "12345"},
		{name: "bare_folder_id", input: "12345", want: "12345"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ParseFolderId(tc.input); got != tc.want {
				t.Fatalf("ParseFolderId(%#v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRemainingTimeoutAndSleepWithDeadline(t *testing.T) {
	cases := []struct {
		name           string
		deadlineOffset time.Duration
		sleepDuration  time.Duration
		wantRemErr     bool
		wantSleepErr   bool
	}{
		{
			name:           "expired_deadline",
			deadlineOffset: -1 * time.Second,
			sleepDuration:  100 * time.Millisecond,
			wantRemErr:     true,
			wantSleepErr:   true,
		},
		{
			name:           "sleep_exceeds_remaining_deadline",
			deadlineOffset: 50 * time.Millisecond,
			sleepDuration:  5 * time.Second,
			wantRemErr:     false,
			wantSleepErr:   true,
		},
		{
			name:           "valid_future_deadline",
			deadlineOffset: 1 * time.Minute,
			sleepDuration:  1 * time.Millisecond,
			wantRemErr:     false,
			wantSleepErr:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deadline := time.Now().Add(tc.deadlineOffset)
			rem, err := remainingTimeout(deadline)
			if (err != nil) != tc.wantRemErr {
				t.Fatalf("remainingTimeout() err = %v, wantRemErr = %t", err, tc.wantRemErr)
			}
			if !tc.wantRemErr && rem <= 0 {
				t.Fatalf("remainingTimeout() = %v, want > 0", rem)
			}

			start := time.Now()
			sleepErr := sleepWithDeadline(time.Now().Add(tc.deadlineOffset), tc.sleepDuration)
			if (sleepErr != nil) != tc.wantSleepErr {
				t.Fatalf("sleepWithDeadline() err = %v, wantSleepErr = %t", sleepErr, tc.wantSleepErr)
			}
			if elapsed := time.Since(start); elapsed > 1*time.Second {
				t.Fatalf("sleepWithDeadline() took %v, expected prompt completion", elapsed)
			}
		})
	}
}

func TestResourceGoogleProjectSoftDeleteCustomDiff_offline(t *testing.T) {
	r := ResourceGoogleProject()
	r.CustomizeDiff = resourceGoogleProjectSoftDeleteCustomDiff
	meta := &transport_tpg.Config{}

	cases := []struct {
		name              string
		state             *terraform.InstanceState
		cfg               map[string]interface{}
		wantErrSubstr     string
		wantLifecycleDiff bool
	}{
		{
			name:  "create_recover_false",
			state: nil,
			cfg: map[string]interface{}{
				"project_id": "test-pid",
				"name":       "Test Project",
			},
		},
		{
			name:  "create_recover_true_offline_client",
			state: nil,
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Test Project",
				"recover_soft_deleted": true,
			},
		},
		{
			name: "inplace_update_active",
			state: &terraform.InstanceState{
				ID: "projects/test-pid",
				Attributes: map[string]string{
					"project_id": "test-pid",
					"name":       "Test Project",
				},
			},
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Updated Project Name",
				"recover_soft_deleted": true,
			},
		},
		{
			name: "inplace_recover_delete_requested",
			state: &terraform.InstanceState{
				ID: "projects/test-pid",
				Attributes: map[string]string{
					"project_id":           "test-pid",
					"name":                 "Test Project",
					"lifecycle_state":      "DELETE_REQUESTED",
					"recover_soft_deleted": "true",
				},
			},
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Updated Project Name",
				"recover_soft_deleted": true,
			},
			wantLifecycleDiff: true,
		},
		{
			name: "inplace_delete_requested_recover_false_errors",
			state: &terraform.InstanceState{
				ID: "projects/test-pid",
				Attributes: map[string]string{
					"project_id":           "test-pid",
					"name":                 "Test Project",
					"lifecycle_state":      "DELETE_REQUESTED",
					"recover_soft_deleted": "false",
				},
			},
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Test Project",
				"recover_soft_deleted": false,
			},
			wantErrSubstr: "recover_soft_deleted = true",
		},
		{
			name: "inplace_delete_requested_conflicting_org_errors",
			state: &terraform.InstanceState{
				ID: "projects/test-pid",
				Attributes: map[string]string{
					"project_id":           "test-pid",
					"name":                 "Test Project",
					"lifecycle_state":      "DELETE_REQUESTED",
					"org_id":               "11111",
					"recover_soft_deleted": "true",
				},
			},
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Test Project",
				"org_id":               "22222",
				"recover_soft_deleted": true,
			},
			wantErrSubstr: "conflicts with configured org_id",
		},
		{
			name: "inplace_delete_requested_toggle_recover_to_false_succeeds",
			state: &terraform.InstanceState{
				ID: "projects/test-pid",
				Attributes: map[string]string{
					"project_id":           "test-pid",
					"name":                 "Test Project",
					"lifecycle_state":      "DELETE_REQUESTED",
					"recover_soft_deleted": "true",
				},
			},
			cfg: map[string]interface{}{
				"project_id":           "test-pid",
				"name":                 "Test Project",
				"recover_soft_deleted": false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rc := terraform.NewResourceConfigRaw(tc.cfg)
			diff, err := r.SimpleDiff(context.Background(), tc.state, rc, meta)
			if tc.wantErrSubstr != "" {
				if err == nil {
					t.Fatalf("SimpleDiff() expected error containing %q, got nil", tc.wantErrSubstr)
				}
				if !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Fatalf("SimpleDiff() error = %q, want substring %q", err.Error(), tc.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("SimpleDiff() unexpected error: %v", err)
			}
			if diff == nil || diff.Empty() {
				t.Fatalf("SimpleDiff() expected non-empty diff")
			}
			if tc.wantLifecycleDiff {
				attrChange, ok := diff.Attributes["lifecycle_state"]
				if !ok {
					t.Fatalf("expected diff to include lifecycle_state change, got %#v", diff.Attributes)
				}
				if attrChange.Old != "DELETE_REQUESTED" || attrChange.New != "ACTIVE" {
					t.Fatalf("lifecycle_state diff Old=%q New=%q, want Old=DELETE_REQUESTED New=ACTIVE", attrChange.Old, attrChange.New)
				}
			}
		})
	}
}
