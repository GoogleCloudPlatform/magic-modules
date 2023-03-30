package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGetRegionFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceRegion types.String
		ProviderRegion types.String
		ExpectedRegion types.String
		ExpectedError  bool
	}{
		"region is pulled from the resource config value instead of the provider config value, even if both set": {
			ResourceRegion: types.StringValue("foo"),
			ProviderRegion: types.StringValue("bar"),
			ExpectedRegion: types.StringValue("foo"),
		},
		"region is pulled from the provider config value when unset on the resource": {
			ResourceRegion: types.StringNull(),
			ProviderRegion: types.StringValue("bar"),
			ExpectedRegion: types.StringValue("bar"),
		},
		"error when region is not set on the provider or the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			region := getRegionFramework(tc.ResourceRegion, tc.ProviderRegion, &diags)

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
