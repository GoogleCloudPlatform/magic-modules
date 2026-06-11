package vcr

import (
	"magician/provider"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCollectResults(t *testing.T) {
	for _, test := range []struct {
		name     string
		output   string
		expected Result
	}{
		{
			name: "no compound tests",
			output: `--- FAIL: TestAccServiceOneResourceOne (100.00s)
--- PASS: TestAccServiceOneResourceTwo (100.00s)
--- PASS: TestAccServiceTwoResourceOne (100.00s)
--- PASS: TestAccServiceTwoResourceTwo (100.00s)
`,
			expected: Result{
				PassedTests: []string{"TestAccServiceOneResourceTwo", "TestAccServiceTwoResourceOne", "TestAccServiceTwoResourceTwo"},
				FailedTests: []string{"TestAccServiceOneResourceOne"},
			},
		},
		{
			name: "compound tests",
			output: `--- FAIL: TestAccServiceOneResourceOne (100.00s)
--- FAIL: TestAccServiceOneResourceTwo (100.00s)
    --- PASS: TestAccServiceOneResourceTwo/test_one (100.00s)
    --- FAIL: TestAccServiceOneResourceTwo/test_two (100.00s)
--- PASS: TestAccServiceTwoResourceOne (100.00s)
    --- PASS: TestAccServiceTwoResourceOne/test_one (100.00s)
    --- PASS: TestAccServiceTwoResourceOne/test_two (100.00s)
--- PASS: TestAccServiceTwoResourceTwo (100.00s)
`,
			expected: Result{
				PassedTests: []string{
					"TestAccServiceTwoResourceOne",
					"TestAccServiceTwoResourceTwo",
				},
				FailedTests: []string{"TestAccServiceOneResourceOne", "TestAccServiceOneResourceTwo"},
				PassedSubtests: []string{
					"TestAccServiceOneResourceTwo__test_one",
					"TestAccServiceTwoResourceOne__test_one",
					"TestAccServiceTwoResourceOne__test_two",
				},
				FailedSubtests: []string{"TestAccServiceOneResourceTwo__test_two"},
			},
		},
		{
			name: "build failure",
			output: `FAIL	github.com/hashicorp/terraform-provider-google-beta/google-beta/services/corebilling [build failed]
--- PASS: TestAccServiceTwoResourceOne (100.00s)
`,
			expected: Result{
				PassedTests:   []string{"TestAccServiceTwoResourceOne"},
				BuildFailures: []string{"corebilling"},
			},
		},
		{
			name: "build failure in middle",
			output: `Error replaying tests:
error running go: exit status 1
stdout:
FAIL	github.com/hashicorp/terraform-provider-google-beta/google-beta/services/corebilling [build failed]
FAIL
stderr:
go: downloading ...
`,
			expected: Result{
				BuildFailures: []string{"corebilling"},
			},
		},
	} {
		if diff := cmp.Diff(test.expected, collectResult(test.output)); diff != "" {
			t.Errorf("collectResult(%q) got unexpected diff (-want +got):\n%s", test.output, diff)
		}
	}

}

func TestAsyncCassetteUploadPath(t *testing.T) {
	vt := &Tester{
		cassetteBucket: "ci-vcr-cassettes",
	}

	wantNightly := "gs://ci-vcr-cassettes/beta/fixtures/"
	if got := vt.asyncCassetteUploadPath("main", provider.Beta); got != wantNightly {
		t.Errorf("asyncCassetteUploadPath(\"main\", Beta) = %q; want %q", got, wantNightly)
	}

	wantFallback := "gs://ci-vcr-cassettes/beta/refs/heads//fixtures/"
	if got := vt.asyncCassetteUploadPath("", provider.Beta); got != wantFallback {
		t.Errorf("asyncCassetteUploadPath(\"\", Beta) = %q; want %q", got, wantFallback)
	}

	wantPR := "gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/"
	if got := vt.asyncCassetteUploadPath("auto-pr-123", provider.Beta); got != wantPR {
		t.Errorf("asyncCassetteUploadPath(\"auto-pr-123\", Beta) = %q; want %q", got, wantPR)
	}
}
