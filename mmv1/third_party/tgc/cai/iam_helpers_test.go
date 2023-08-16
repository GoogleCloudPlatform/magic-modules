package cai

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		t.Run(c.name+"/MergeAdditiveBindings", func(t *testing.T) {
			assert.EqualValues(t,
				c.expectedAdditive,
				MergeAdditiveBindings(c.existing, c.incoming),
			)
		})
		t.Run(c.name+"/MergeAuthoritativeBindings", func(t *testing.T) {
			assert.EqualValues(t,
				c.expectedAuthoritative,
				MergeAuthoritativeBindings(c.existing, c.incoming),
			)
		})
	}
}

func TestMergeDeleteBindings(t *testing.T) {
	cases := []struct {
		name string
		// Inputs
		existing []IAMBinding
		incoming []IAMBinding
		// Expected outputs
		expectedDeleteAdditive      []IAMBinding
		expectedDeleteAuthoritative []IAMBinding
	}{
		{
			name:                        "EmptyDeleteEmpty",
			existing:                    []IAMBinding{},
			incoming:                    []IAMBinding{},
			expectedDeleteAdditive:      nil,
			expectedDeleteAuthoritative: nil,
		},
		{
			name:     "EmptyDeleteOne",
			existing: []IAMBinding{},
			incoming: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedDeleteAdditive:      nil,
			expectedDeleteAuthoritative: nil,
		},
		{
			name: "OneDeleteEmpty",
			existing: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []IAMBinding{},
			expectedDeleteAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedDeleteAuthoritative: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneDeleteOne",
			existing: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b"},
				},
			},
			incoming: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
			expectedDeleteAdditive: []IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedDeleteAuthoritative: nil,
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
			expectedDeleteAdditive: []IAMBinding{
				{
					Role:    "role-b",
					Members: []string{"member-d"},
				},
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
			expectedDeleteAuthoritative: []IAMBinding{
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name+"/MergeDeleteAdditiveBindings", func(t *testing.T) {
			assert.EqualValues(t,
				c.expectedDeleteAdditive,
				MergeDeleteAdditiveBindings(c.existing, c.incoming),
			)
		})
		t.Run(c.name+"/MergeDeleteAuthoritativeBindings", func(t *testing.T) {
			assert.EqualValues(t,
				c.expectedDeleteAuthoritative,
				MergeDeleteAuthoritativeBindings(c.existing, c.incoming),
			)
		})
	}
}
