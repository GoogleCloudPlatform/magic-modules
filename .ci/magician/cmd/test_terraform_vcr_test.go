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
		wantNotRun       []string
	}{
		"no diff": {
			gaDiff:   "",
			betaDiff: "",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"no added tests": {
			gaDiff:   "+// some change",
			betaDiff: "+// some change",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"tests added and run": {
			gaDiff:   "+func TestAccOne(t *testing.T) {",
			betaDiff: "+func TestAccTwo(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"tests removed and run": {
			gaDiff:   "-func TestAccOne(t *testing.T) {",
			betaDiff: "-func TestAccTwo(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"tests added and not run": {
			gaDiff:   "+func TestAccThree(t *testing.T) {",
			betaDiff: "+func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{"TestAccThree", "TestAccFour"},
		},
		"tests removed and not run": {
			gaDiff:   "-func TestAccThree(t *testing.T) {",
			betaDiff: "-func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"tests added but commented out": {
			gaDiff:   "+//func TestAccThree(t *testing.T) {",
			betaDiff: "+//func TestAccFour(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{},
		},
		"tests added and not run are deduped": {
			gaDiff:   "+func TestAccThree(t *testing.T) {",
			betaDiff: "+func TestAccThree(t *testing.T) {",
			result: &vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRun: []string{"TestAccThree"},
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			notRun := notRunTests(tc.gaDiff, tc.betaDiff, tc.result)
			assert.Equal(t, tc.wantNotRun, notRun)
		})
	}
}
