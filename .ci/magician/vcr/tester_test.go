package vcr

import (
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
	} {
		if diff := cmp.Diff(test.expected, collectResult(test.output)); diff != "" {
			t.Errorf("collectResult(%q) got unexpected diff (-want +got):\n%s", test.output, diff)
		}
	}

}
