package cmd

import (
	"reflect"
	"testing"
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
		if packages, all, _ := modifiedPackagesFromDiffs(tc.diffs); !reflect.DeepEqual(packages, tc.packages) {
			t.Errorf("Unexpected packages found for test %s: %v, expected %v", tc.name, packages, tc.packages)
		} else if all != tc.all {
			t.Errorf("Unexpected value for all packages for test %s: %v, expected %v", tc.name, all, tc.all)
		}
	}
}
