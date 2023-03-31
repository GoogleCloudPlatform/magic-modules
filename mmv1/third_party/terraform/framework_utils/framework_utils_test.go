package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGetRegionFramework(t *testing.T) {
	cases := map[string]struct {
		Resource       ResourceFoobar
		Provider       ProviderFoobar
		ExpectedRegion types.String
		ExpectedError  bool
	}{
		"region is pulled from the resource config's region value if available": {
			Resource: ResourceFoobar{
				Region: types.StringValue("foo"),
			},
			ExpectedRegion: types.StringValue("foo"),
		},
		"region is pulled from the resource config's zone value if region is unset": {
			Resource: ResourceFoobar{
				Region: types.StringNull(),
				Zone:   types.StringValue("foo-a"),
			},
			ExpectedRegion: types.StringValue("foo"),
		},
		"region is pulled from the provider config's region value when region and zone are unset on the resource": {
			Resource: ResourceFoobar{
				Region: types.StringNull(),
				Zone:   types.StringNull(),
			},
			Provider: ProviderFoobar{
				Region: types.StringValue("bar"),
			},
			ExpectedRegion: types.StringValue("bar"),
		},
		"region is pulled from the provider config's zone value when region is unset on the provider (and resource config lacks region/zone)": {
			Resource: ResourceFoobar{
				Region: types.StringNull(),
				Zone:   types.StringNull(),
			},
			Provider: ProviderFoobar{
				Region: types.StringValue("bar"),
				Zone:   types.StringValue("bar-a"),
			},
			ExpectedRegion: types.StringValue("bar"),
		},
		"error when region and zone are not set on the provider nor the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			region := getRegionFramework(tc.Resource, tc.Provider, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}
