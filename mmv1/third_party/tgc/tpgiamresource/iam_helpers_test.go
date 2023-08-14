package tpgiamresource

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	"github.com/stretchr/testify/assert"
)

func TestMergeBindings(t *testing.T) {
	cases := []struct {
		name string
		// Inputs
		existing []tpgresource.IAMBinding
		incoming []tpgresource.IAMBinding
		// Expected outputs
		expectedAdditive      []tpgresource.IAMBinding
		expectedAuthoritative []tpgresource.IAMBinding
	}{
		{
			name:                  "EmptyAddEmpty",
			existing:              []tpgresource.IAMBinding{},
			incoming:              []tpgresource.IAMBinding{},
			expectedAdditive:      []tpgresource.IAMBinding{},
			expectedAuthoritative: []tpgresource.IAMBinding{},
		},
		{
			name:     "EmptyAddOne",
			existing: []tpgresource.IAMBinding{},
			incoming: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAuthoritative: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneAddEmpty",
			existing: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []tpgresource.IAMBinding{},
			expectedAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedAuthoritative: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneAddOne",
			existing: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
			expectedAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b"},
				},
			},
			expectedAuthoritative: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
		},
		{
			name: "GrandFinale",
			existing: []tpgresource.IAMBinding{
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
			incoming: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b", "member-c"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-b", "member-c"},
				},
			},
			expectedAdditive: []tpgresource.IAMBinding{
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
			expectedAuthoritative: []tpgresource.IAMBinding{
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
		existing []tpgresource.IAMBinding
		incoming []tpgresource.IAMBinding
		// Expected outputs
		expectedDeleteAdditive      []tpgresource.IAMBinding
		expectedDeleteAuthoritative []tpgresource.IAMBinding
	}{
		{
			name:                        "EmptyDeleteEmpty",
			existing:                    []tpgresource.IAMBinding{},
			incoming:                    []tpgresource.IAMBinding{},
			expectedDeleteAdditive:      nil,
			expectedDeleteAuthoritative: nil,
		},
		{
			name:     "EmptyDeleteOne",
			existing: []tpgresource.IAMBinding{},
			incoming: []tpgresource.IAMBinding{
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
			existing: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			incoming: []tpgresource.IAMBinding{},
			expectedDeleteAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedDeleteAuthoritative: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
		},
		{
			name: "OneDeleteOne",
			existing: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b"},
				},
			},
			incoming: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-b"},
				},
			},
			expectedDeleteAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a"},
				},
			},
			expectedDeleteAuthoritative: nil,
		},
		{
			name: "GrandFinale",
			existing: []tpgresource.IAMBinding{
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
			incoming: []tpgresource.IAMBinding{
				{
					Role:    "role-a",
					Members: []string{"member-a", "member-b", "member-c"},
				},
				{
					Role:    "role-b",
					Members: []string{"member-b", "member-c"},
				},
			},
			expectedDeleteAdditive: []tpgresource.IAMBinding{
				{
					Role:    "role-b",
					Members: []string{"member-d"},
				},
				{
					Role:    "role-c",
					Members: []string{"member-c"},
				},
			},
			expectedDeleteAuthoritative: []tpgresource.IAMBinding{
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
