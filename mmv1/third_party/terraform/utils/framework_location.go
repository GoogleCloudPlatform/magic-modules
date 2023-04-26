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
	if !ld.ProviderRegion.IsNull() {
		return ld.ProviderRegion, nil
	}
	// Region from zone in provider config
	if !ld.ProviderZone.IsNull() {
		region := getRegionFromZone(ld.ProviderZone.ValueString())
		return types.StringValue(region), nil
	}

	var err error
	if !ld.RegionSchemaField.IsNull() {
		err = fmt.Errorf("region could not be identified, please add `%s` in your resource or set `region` in your provider configuration block", ld.ZoneSchemaField.ValueString())
	} else {
		err = errors.New("region could not be identified, please add `region` in your resource or provider configuration block")
	}
	return types.StringNull(), err
}

func (ld *LocationDescription) getZone() (types.String, error) {
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() {
		// Zone could be a self link
		zone := GetResourceNameFromSelfLink(ld.ResourceZone.ValueString())
		return types.StringValue(zone), nil
	}
	if !ld.ProviderZone.IsNull() {
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
