package google

import (
	"reflect"
	"testing"
)

func TestStringCamelize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		term        string
		firstLetter string
		expected    string
	}{
		{
			description: "Camelize string with lowercase first letter",
			term:        "AuthorizedOrgsDesc",
			firstLetter: "lower",
			expected:    "authorizedOrgsDesc",
		},
		{
			description: "Camelize string with uppercase first letter",
			term:        "AuthorizedOrgsDesc",
			firstLetter: "upper",
			expected:    "AuthorizedOrgsDesc",
		},
		{
			description: "Camelize snakecase string with lowercase first letter",
			term:        "Authorized_Orgs_Desc",
			firstLetter: "lower",
			expected:    "authorizedOrgsDesc",
		},
		{
			description: "Camelize snakecase string with uppercase first letter",
			term:        "Authorized_Orgs_Desc",
			firstLetter: "upper",
			expected:    "AuthorizedOrgsDesc",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := Camelize(tc.term, tc.firstLetter), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}
