package google

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLocationDescription_getZone(t *testing.T) {
	cases := map[string]struct {
		ld            LocationDescription
		ExpectedZone  types.String
		ExpectedError bool
	}{
		"returns the value of the zone field in resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("resource-zone"),
				ProviderZone: types.StringValue("provider-zone"),
			},
			ExpectedZone: types.StringValue("resource-zone"),
		},
		"shortens zone values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a"),
			},
			ExpectedZone: types.StringValue("us-central1-a"),
		},
		"returns the value of the zone field in provider config when zone is unset in resource config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("provider-zone"),
			},
			ExpectedZone: types.StringValue("provider-zone"),
		},
		"returns an error when a zone value can't be found": {
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when a zone value can't be found": {
			ld: LocationDescription{
				ZoneSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			zone, err := tc.ld.getZone()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.ZoneSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.ZoneSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.ZoneSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if zone != tc.ExpectedZone {
				t.Fatalf("Incorrect zone: got %s, want %s", zone, tc.ExpectedZone)
			}
		})
	}
}

func TestLocationDescription_getRegion(t *testing.T) {
	cases := map[string]struct {
		ld             LocationDescription
		ExpectedRegion types.String
		ExpectedError  bool
	}{
		"returns the value of the region field in resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("resource-region"),
				ProviderRegion: types.StringValue("provider-region"),
			},
			ExpectedRegion: types.StringValue("resource-region"),
		},
		"shortens region values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			},
			ExpectedRegion: types.StringValue("us-central1"),
		},
		"returns a region derived from the zone field in resource config when region is unset": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("provider-zone-a"),
			},
			ExpectedRegion: types.StringValue("provider-zone"),
		},
		"does not shorten region values when derived from a zone self link set in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a"),
			},
			ExpectedRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1"), // Value isn't sortened from URI to name
		},
		"returns the value of the region field in provider config when region/zone is unset in resource config": {
			ld: LocationDescription{
				ProviderRegion: types.StringValue("provider-region"),
			},
			ExpectedRegion: types.StringValue("provider-region"),
		},
		"returns a region derived from the zone field in provider config when region unset in both resource and provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("provider-zone-a"),
			},
			ExpectedRegion: types.StringValue("provider-zone"),
		},
		"returns an error when zone values can't be found": {
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when region value can't be found": {
			ld: LocationDescription{
				ZoneSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region, err := tc.ld.getRegion()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.RegionSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.RegionSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.RegionSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}
