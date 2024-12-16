package main

import (
	"regexp"
	"testing"
)

func Test_main_happyPaths(t *testing.T) {
	testCases := map[string]struct {
		providerServiceFile string
		teamcityServiceFile string
		expectError         bool
		errorRegex          *regexp.Regexp
		missingServiceRegex *regexp.Regexp
	}{
		"everything matches": {
			providerServiceFile: "./test-fixtures/everything-ok/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/everything-ok/services_ga.kt",
		},
		"something missing in TeamCity config present in provider code": {
			providerServiceFile: "./test-fixtures/mismatch-teamcity/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/mismatch-teamcity/services_ga.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("TeamCity service file is missing services present in the provider"),
			missingServiceRegex: regexp.MustCompile("[pubsub]"),
		},
		"something missing in provider code present in TeamCity config": {
			providerServiceFile: "./test-fixtures/mismatch-provider/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/mismatch-provider/services_ga.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("Provider codebase is missing services present in the TeamCity service file"),
			missingServiceRegex: regexp.MustCompile("[compute]"),
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			err := compareServices(tc.teamcityServiceFile, tc.providerServiceFile)
			if err != nil && !tc.expectError {
				t.Fatalf("saw an unexpected error: %s", err)
			}
			if err == nil && tc.expectError {
				t.Fatalf("expected an error but saw none")
			}

			if err == nil {
				// Stop handling of non-error test cases
				return
			}

			if !tc.errorRegex.MatchString(err.Error()) {
				t.Fatalf("expected error to contain a match for regex `%s`, got error string: `%s`", tc.errorRegex.String(), err)
			}
			if !tc.missingServiceRegex.MatchString(err.Error()) {
				t.Fatalf("expected error to contain a match for regex `%s`, got error string: `%s`", tc.errorRegex.String(), err)
			}
		})
	}
}

func Test_main_unhappyPaths(t *testing.T) {
	testCases := map[string]struct {
		providerServiceFile string
		teamcityServiceFile string
		expectError         bool
		errorRegex          *regexp.Regexp
	}{
		"cannot find provider service file": {
			providerServiceFile: "./test-fixtures/doesnt-exist.txt",
			teamcityServiceFile: "./test-fixtures/everything-ok/services_ga.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("error opening provider service list file: open ./test-fixtures/doesnt-exist.txt"),
		},
		"cannot find TeamCity service file": {
			providerServiceFile: "./test-fixtures/everything-ok/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/everything-ok/doesnt-exist.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("error opening TeamCity service list file: open ./test-fixtures/everything-ok/doesnt-exist.kt"),
		},
		"empty TeamCity service file": {
			providerServiceFile: "./test-fixtures/everything-ok/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/empty-files/services_ga.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("could not find any services in the TeamCity service list file ./test-fixtures/empty-files/services_ga.kt"),
		},
		"empty provider service file": {
			providerServiceFile: "./test-fixtures/empty-files/ga-services.txt",
			teamcityServiceFile: "./test-fixtures/everything-ok/services_ga.kt",
			expectError:         true,
			errorRegex:          regexp.MustCompile("could not find any services in the provider service list file ./test-fixtures/empty-files/ga-services.txt"),
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			err := compareServices(tc.teamcityServiceFile, tc.providerServiceFile)
			if err != nil && !tc.expectError {
				t.Fatalf("saw an unexpected error: %s", err)
			}
			if err == nil && tc.expectError {
				t.Fatalf("expected an error but saw none")
			}

			if !tc.errorRegex.MatchString(err.Error()) {
				t.Fatalf("expected error to contain a match for regex `%s`, got error string: `%s`", tc.errorRegex.String(), err)
			}
		})
	}
}

func Test_listDifference(t *testing.T) {
	testCases := map[string]struct {
		a          []string
		b          []string
		expectDiff bool
		errorRegex *regexp.Regexp
	}{
		"detects when lists match": {
			a: []string{"a", "c", "b"},
			b: []string{"a", "b", "c"},
		},
		"detects when items from list A is missing items present in list B - 1 missing": {
			a:          []string{"a", "b"},
			b:          []string{"a", "c", "b"},
			expectDiff: true,
			errorRegex: regexp.MustCompile("[c]"),
		},
		"detects when items from list A is missing items present in list B - 2 missing": {
			a:          []string{"b"},
			b:          []string{"a", "c", "b"},
			expectDiff: true,
			errorRegex: regexp.MustCompile("[a, c]"),
		},
		"doesn't detect differences if list A is a superset of list B": {
			a:          []string{"a", "b", "c"},
			b:          []string{"a", "c"},
			expectDiff: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			err := listDifference(tc.a, tc.b)
			if !tc.expectDiff && (err != nil) {
				t.Fatalf("saw an unexpected diff error: %s", err)
			}
			if tc.expectDiff && (err == nil) {
				t.Fatalf("expected a diff error but saw none")
			}
			if !tc.expectDiff && err == nil {
				// Stop assertions in no error cases
				return
			}
			if !tc.errorRegex.MatchString(err.Error()) {
				t.Fatalf("expected diff error to contain a match for regex %s, error string: %s", tc.errorRegex.String(), err)
			}
		})
	}
}
