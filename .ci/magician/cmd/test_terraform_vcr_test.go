package cmd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"magician/vcr"
)

func TestModifiedPackagesFromDiffs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		diffs    []string
		packages map[string]struct{}
		all      bool
	}{
		{
			name:     "one-package",
			diffs:    []string{"google-beta/services/servicename/resource.go"},
			packages: map[string]struct{}{"servicename": struct{}{}},
			all:      false,
		},
		{
			name: "multiple-packages",
			diffs: []string{
				"google-beta/services/serviceone/resource.go",
				"google-beta/services/servicetwo/test-fixtures/fixture.txt",
				"google-beta/services/servicethree/resource_test.go",
			},
			packages: map[string]struct{}{
				"serviceone":   struct{}{},
				"servicetwo":   struct{}{},
				"servicethree": struct{}{},
			},
			all: false,
		},
		{
			name:     "all-packages",
			diffs:    []string{"google-beta/provider/provider.go"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "all-packages-go-mod",
			diffs:    []string{"scripts/go.mod"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "all-packages-go-sum",
			diffs:    []string{"go.sum"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "no-packages",
			diffs:    []string{"website/docs/d/notebooks_runtime_iam_policy.html.markdown"},
			packages: map[string]struct{}{},
			all:      false,
		},
	} {
		if packages, all := modifiedPackages(tc.diffs); !reflect.DeepEqual(packages, tc.packages) {
			t.Errorf("Unexpected packages found for test %s: %v, expected %v", tc.name, packages, tc.packages)
		} else if all != tc.all {
			t.Errorf("Unexpected value for all packages for test %s: %v, expected %v", tc.name, all, tc.all)
		}
	}
}

func TestNotRunTests(t *testing.T) {
	cases := map[string]struct {
		gaDiff, betaDiff string
		result           *vcr.Result
		wantNotRunBeta   []string
		wantNotRunGa     []string
	}{
		"no diff": {
			gaDiff:   "",
			betaDiff: "",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"no added tests": {
			gaDiff:   "+// some change",
			betaDiff: "+// some change",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and passed": {
			gaDiff:   "+func TestAccTwo(t *testing.T) {",
			betaDiff: "+func TestAccTwo(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccTwo"},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"multiple tests added and passed": {
			gaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			betaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			result: &vcr.Result{
				PassedTests: []string{"TestAccTwo", "TestAccThree"},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and failed": {
			gaDiff:   "+func TestAccTwo(t *testing.T) {",
			betaDiff: "+func TestAccTwo(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"tests removed and run": {
			gaDiff:   "-func TestAccOne(t *testing.T) {",
			betaDiff: "-func TestAccTwo(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and not run": {
			gaDiff:   "+func TestAccThree(t *testing.T) {",
			betaDiff: "+func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{"TestAccFour"},
			wantNotRunGa:   []string{"TestAccThree"},
		},
		"multiple tests added and not run": {
			gaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			betaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccFour"},
			},
			wantNotRunBeta: []string{"TestAccThree", "TestAccTwo"},
			wantNotRunGa:   []string{},
		},
		"tests removed and not run": {
			gaDiff:   "-func TestAccThree(t *testing.T) {",
			betaDiff: "-func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"tests added but commented out": {
			gaDiff:   "+//func TestAccThree(t *testing.T) {",
			betaDiff: "+//func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"multiline diffs": {
			gaDiff: `diff --git a/google/services/alloydb/resource_alloydb_backup_generated_test.go b/google/services/alloydb/resource_alloydb_backup_generated_test.go
+func TestAccAlloydbBackup_alloydbBackupFullTestNewExample(t *testing.T) {
+func TestAccCloudRunService_cloudRunServiceMulticontainerExample(t *testing.T) {`,
			betaDiff: `diff --git a/google-beta/services/alloydb/resource_alloydb_backup_generated_test.go b/google-beta/services/alloydb/resource_alloydb_backup_generated_test.go
+func TestAccAlloydbBackup_alloydbBackupFullTestNewExample(t *testing.T) {`,
			result: &vcr.Result{
				PassedTests: []string{},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{"TestAccAlloydbBackup_alloydbBackupFullTestNewExample"},
			wantNotRunGa:   []string{"TestAccCloudRunService_cloudRunServiceMulticontainerExample"},
		},
		"always count GA-only added tests": {
			gaDiff:   "+func TestAccOne(t *testing.T) {",
			betaDiff: "",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{"TestAccOne"},
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			notRunBeta, notRunGa := notRunTests(tc.gaDiff, tc.betaDiff, tc.result)
			assert.Equal(t, tc.wantNotRunBeta, notRunBeta)
			assert.Equal(t, tc.wantNotRunGa, notRunGa)
		})
	}
}
