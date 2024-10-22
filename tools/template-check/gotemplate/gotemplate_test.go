package ruby

import (
	"strings"
	"testing"
)

func TestCheckVersionGuards(t *testing.T) {
	cases := map[string]struct {
		fileText        string
		expectedResults []string
	}{
		"allow standard format targeting ga": {
			fileText:        `some text
							{{- if ne $.TargetVersionName "ga" }}
							more text`,
			expectedResults: nil,
		},
		"disallow targeting beta": {
			fileText:        `some text
							{{- if ne $.TargetVersionName "beta" }}
							more text`,
			expectedResults: []string{"2: {{- if ne $.TargetVersionName "beta" }}"},
		},
		"one valid, one invalid": {
			fileText:        `some text
							{{- if ne $.TargetVersionName "beta" }}
							more text
							{{- if ne $.TargetVersionName "ga" }}`,
			expectedResults: []string{"2: {{- if ne $.TargetVersionName "beta" }}"},
		},
		"multiple invalid": {
			fileText:        `some text
							{{- if ne $.TargetVersionName "beta" }}
							more text
							
							
							{{- if eq $.TargetVersionName "beta" }}`,
			expectedResults: []string{`2: {{- if ne $.TargetVersionName "beta" }}`, `6: {{- if eq $.TargetVersionName "beta" }}`},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			results := CheckVersionGuards(strings.NewReader(tc.fileText))
			if len(results) != len(tc.expectedResults) {
				t.Errorf("Expected length %d, got %d", len(tc.expectedResults), len(results))
				return
			}
			for i, result := range results {
				if result != tc.expectedResults[i] {
					t.Errorf("Expected %v, got %v", tc.expectedResults[i], result)
				}
			}
		})
	}
}
