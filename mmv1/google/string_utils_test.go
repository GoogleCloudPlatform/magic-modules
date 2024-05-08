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
			term:        "authorizedOrgsDesc",
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

func TestStringPlural(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		term        string
		expected    string
	}{
		{
			description: "Plural normal string",
			term:        "apple",
			expected:    "apples",
		},
		{
			description: "Plural string ending with ies",
			term:        "policies",
			expected:    "policies",
		},
		{
			description: "Plural snakecase string ending with es",
			term:        "indices",
			expected:    "indices",
		},
		{
			description: "Plural snakecase string ending with ex",
			term:        "index",
			expected:    "indices",
		},
		{
			description: "Plural snakecase string ending with esh",
			term:        "mesh",
			expected:    "meshes",
		},
		{
			description: "Plural snakecase string ending with y",
			term:        "policy",
			expected:    "policies",
		},
		{
			description: "Plural snakecase string ending with ey",
			term:        "key",
			expected:    "keys",
		},
		{
			description: "Plural snakecase string ending with ay",
			term:        "gateway",
			expected:    "gateways",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := Plural(tc.term), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestStringFirstSentence(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		term        string
		expected    string
	}{
		{
			description: "sentence end with period",
			term:        "Lorem ipsum. Dolor sit amet. Elit",
			expected:    "Lorem ipsum.",
		},
		{
			description: "sentence end with question mark",
			term:        "Lorem ipsum? Dolor sit amet. Elit",
			expected:    "Lorem ipsum?",
		},
		{
			description: "sentence end with exclamation mark",
			term:        "Lorem ipsum! Dolor sit amet. Elit",
			expected:    "Lorem ipsum!",
		},
		{
			description: "no period returns full string",
			term:        "Lorem ipsum dolor",
			expected:    "Lorem ipsum dolor",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := FirstSentence(tc.term), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}
