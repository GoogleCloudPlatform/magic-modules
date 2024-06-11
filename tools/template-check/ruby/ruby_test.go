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
			fileText:        "some text\n<% unless version == 'ga' -%>\nmore text",
			expectedResults: nil,
		},
		"disallow targeting beta": {
			fileText:        "some text\n<% unless version == 'beta' -%>\nmore text",
			expectedResults: []string{"2: <% unless version == 'beta' -%>"},
		},
		"disallow 'if not'": {
			fileText:        "some text\n<% if version != 'ga' -%>\nmore text",
			expectedResults: []string{"2: <% if version != 'ga' -%>"},
		},
		"disallow single '='": {
			fileText:        "some text\n<% unless version = 'ga' -%>\nmore text",
			expectedResults: []string{"2: <% unless version = 'ga' -%>"},
		},
		"disallow leaving trailing line break": {
			fileText:        "some text\n<% unless version == 'ga' %>\nmore text",
			expectedResults: []string{"2: <% unless version == 'ga' %>"},
		},
		"one valid, one invalid": {
			fileText:        "some text\n<% unless version == 'beta' -%>\nmore text\n<% unless version == 'ga' -%>",
			expectedResults: []string{"2: <% unless version == 'beta' -%>"},
		},
		"multiple invalid": {
			fileText:        "some text\n<% unless version == 'beta' -%>\nmore text\n\n\n<% if version == \"beta\" -%>",
			expectedResults: []string{"2: <% unless version == 'beta' -%>", "6: <% if version == \"beta\" -%>"},
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
