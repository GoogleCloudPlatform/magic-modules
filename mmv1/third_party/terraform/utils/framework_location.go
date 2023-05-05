package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocationDescriber interface {
	getLocationDescription(providerConfig *frameworkProvider) LocationDescription
}

type LocationDescription struct {
	LocationSchemaField types.String
	RegionSchemaField   types.String
	ZoneSchemaField     types.String

	ResourceLocation types.String
	ResourceRegion   types.String
	ResourceZone     types.String

	ProviderRegion types.String
	ProviderZone   types.String
}

func (ld *LocationDescription) getLocation() (types.String, error) {
	// Location from resource config
	if !ld.ResourceLocation.IsNull() && !ld.ResourceLocation.IsUnknown() {
		return ld.ResourceLocation, nil
	}

	// Location from region in resource config
	if !ld.ResourceRegion.IsNull() && !ld.ResourceRegion.IsUnknown() {
		return ld.ResourceRegion, nil
	}

	// Location from zone in resource config
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() {
		location := GetResourceNameFromSelfLink(ld.ResourceZone.ValueString()) // Zone could be a self link
		return types.StringValue(location), nil
	}

	// Location from zone in provider config
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() {
		return ld.ProviderZone, nil
	}

	return types.StringNull(), nil
}

func (ld *LocationDescription) getRegion() (types.String, error) {
	// Region from resource config
	if !ld.ResourceRegion.IsNull() && !ld.ResourceRegion.IsUnknown() {
		region := GetResourceNameFromSelfLink(ld.ResourceRegion.ValueString()) // Region could be a self link
		return types.StringValue(region), nil
	}
	// Region from zone in resource config
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() {
		region := getRegionFromZone(ld.ResourceZone.ValueString())
		return types.StringValue(region), nil
	}
	// Region from provider config
	if !ld.ProviderRegion.IsNull() && !ld.ProviderRegion.IsUnknown() {
		return ld.ProviderRegion, nil
	}
	// Region from zone in provider config
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() {
		region := getRegionFromZone(ld.ProviderZone.ValueString())
		return types.StringValue(region), nil
	}

	var err error
	if !ld.RegionSchemaField.IsNull() {
		err = fmt.Errorf("region could not be identified, please add `%s` in your resource or set `region` in your provider configuration block", ld.RegionSchemaField.ValueString())
	} else {
		err = errors.New("region could not be identified, please add `region` in your resource or provider configuration block")
	}
	return types.StringNull(), err
}

func (ld *LocationDescription) getZone() (types.String, error) {
	// TODO(SarahFrench): Make empty strings not ignored, see https://github.com/hashicorp/terraform-provider-google/issues/14447
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() && !ld.ResourceZone.Equal(types.StringValue("")) {
		// Zone could be a self link
		zone := GetResourceNameFromSelfLink(ld.ResourceZone.ValueString())
		return types.StringValue(zone), nil
	}
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() {
		return ld.ProviderZone, nil
	}

	var err error
	if !ld.ZoneSchemaField.IsNull() {
		err = fmt.Errorf("zone could not be identified, please add `%s` in your resource or `zone` in your provider configuration block", ld.ZoneSchemaField.ValueString())
	} else {
		err = errors.New("zone could not be identified, please add `zone` in your resource or provider configuration block")
	}
	return types.StringNull(), err
}
