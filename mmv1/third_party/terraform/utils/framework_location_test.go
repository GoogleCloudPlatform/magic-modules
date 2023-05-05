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
				// A resource would not have all 3 fields set, but if they were all present zone is used
				ResourceZone:     types.StringValue("resource-zone"),
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceLocation: types.StringValue("resource-location"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone"),
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
		"returns the value of the zone field in provider config when zone is set to an empty string in resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue(""),
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
				// A resource would not have all 3 fields set, but if they were all present region is used first
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceLocation: types.StringValue("resource-location"),
				ResourceZone:     types.StringValue("resource-zone"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone"),
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
				RegionSchemaField: types.StringValue("foobar"),
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

func TestLocationDescription_getLocation(t *testing.T) {
	cases := map[string]struct {
		ld               LocationDescription
		ExpectedLocation types.String
		ExpectedError    bool
	}{
		"returns the value of the location field in resource config": {
			ld: LocationDescription{
				// A resource would not have all 3 fields set, but if they were all present location is used first
				ResourceLocation: types.StringValue("resource-location"),
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceZone:     types.StringValue("resource-zone"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone"),
			},
			ExpectedLocation: types.StringValue("resource-location"),
		},
		"does not shorten the location value when it is set as a self link in the resource config": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location"),
			},
			ExpectedLocation: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location"),
		},
		"returns the region value set in the resource config when location is not in the schema": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("resource-region"),
			},
			ExpectedLocation: types.StringValue("resource-region"),
		},
		"does not shorten the region value when it is set as a self link in the resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/resource-region"),
			},
			ExpectedLocation: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/resource-region"),
		},
		"returns the zone value set in the resource config when neither location nor region in the schema": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("resource-zone"),
			},
			ExpectedLocation: types.StringValue("resource-zone"),
		},
		"shortens zone values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/resource-zone"),
			},
			ExpectedLocation: types.StringValue("resource-zone"),
		},
		"returns the zone value from the provider config when none of location/region/zone are set in the resource config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("provider-zone"),
			},
			ExpectedLocation: types.StringValue("provider-zone"),
		},
		"does not shorten the zone value when it is set as a self link in the provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone"),
			},
			ExpectedLocation: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone"),
		},
		"does not use the region value set in the provider config": {
			ld: LocationDescription{
				ProviderRegion: types.StringValue("provider-zone"),
			},
			ExpectedError: true,
		},
		"returns an error when none of location/region/zone are set on the resource, and neither region or zone is set on the provider": {
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when location value can't be found": {
			ld: LocationDescription{
				LocationSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region, err := tc.ld.getLocation()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.LocationSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.LocationSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.LocationSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if region != tc.ExpectedLocation {
				t.Fatalf("Incorrect location: got %s, want %s", region, tc.ExpectedLocation)
			}
		})
	}
}
