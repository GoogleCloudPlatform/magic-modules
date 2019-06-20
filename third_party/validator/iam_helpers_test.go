package google

import (
	"fmt"
	"testing"
)

func TestMergeBindings(t *testing.T) {
	cases := []struct {
		name string
		// Inputs
		existing []IAMBinding
		incoming []IAMBinding
		// Expected outputs
		expectedAdditive      []IAMBinding
		expectedAuthoritative []IAMBinding
	}{
		{
			name:                  "EmptyAddEmpty",
			existing:              []IAMBinding{},
			incoming:              []IAMBinding{},
			expectedAdditive:      []IAMBinding{},
			expectedAuthoritative: []IAMBinding{},
		},
		{
			name:     "EmptyAddOne",
			existing: []IAMBinding{},
			incoming: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAuthoritative: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneAddEmpty",
			existing: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []IAMBinding{},
			expectedAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAuthoritative: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneAddOne",
			existing: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
			expectedAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b"},
				},
			},
			expectedAuthoritative: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
		},
		{
			name: "GrandFinale",
			existing: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-c", "member-d"},
				},
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
			incoming: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b", "member-c"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-b", "member-c"},
				},
			},
			expectedAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b", "member-c"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-b", "member-c", "member-d"},
				},
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
			expectedAuthoritative: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b", "member-c"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-b", "member-c"},
				},
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertEqual(t,
				c.expectedAdditive,
				mergeAdditiveBindings(c.existing, c.incoming),
				"mergeAdditiveBindings",
			)
			assertEqual(t,
				c.expectedAuthoritative,
				mergeAuthoritativeBindings(c.existing, c.incoming),
				"mergeAuthoritativeBindings",
			)
		})
	}
}

// assertEqual compares two values by converting them to strings.
func assertEqual(t *testing.T, exp, got interface{}, name string) {
	expS, gotS := fmt.Sprintf("%+v", exp),
		fmt.Sprintf("%+v", got)

	if expS != gotS {
		t.Errorf("%s: expected %s, got %s", name, expS, gotS)
	}
}
