package compute

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ComputeInstanceGroupAssetType string = "compute.googleapis.com/InstanceGroup"

func ResourceConverterComputeInstanceGroup() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ComputeInstanceGroupAssetType,
		Convert:   GetComputeInstanceGroupCaiObject,
	}
}

func GetComputeInstanceGroupCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/instanceGroups/{{group}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetComputeInstanceGroupApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ComputeInstanceGroupAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "InstanceGroup",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetComputeInstanceGroupApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandComputeInstanceGroupName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandComputeInstanceGroupDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	namedPortsProp, err := expandComputeInstanceGroupNamedPorts(d.Get("named_ports"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("named_ports"); !tpgresource.IsEmptyValue(reflect.ValueOf(namedPortsProp)) && (ok || !reflect.DeepEqual(v, namedPortsProp)) {
		obj["namedPorts"] = namedPortsProp
	}
	networkProp, err := expandComputeInstanceGroupNetwork(d.Get("network"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("network"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	fingerprintProp, err := expandComputeInstanceGroupFingerprint(d.Get("fingerprint"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("fingerprint"); !tpgresource.IsEmptyValue(reflect.ValueOf(fingerprintProp)) && (ok || !reflect.DeepEqual(v, fingerprintProp)) {
		obj["fingerprint"] = fingerprintProp
	}
	zoneProp, err := expandComputeInstanceGroupZone(d.Get("zone"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("zone"); !tpgresource.IsEmptyValue(reflect.ValueOf(zoneProp)) && (ok || !reflect.DeepEqual(v, zoneProp)) {
		url, err := tpgresource.ReplaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/zones/")
		if err != nil {
			return nil, err
		}
		url = url + zoneProp.(string)

		obj["zone"] = url
	}
	selfLinkProp, err := expandComputeInstanceGroupSelfLink(d.Get("self_link"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("self_link"); !tpgresource.IsEmptyValue(reflect.ValueOf(selfLinkProp)) && (ok || !reflect.DeepEqual(v, selfLinkProp)) {
		obj["selfLink"] = selfLinkProp
	}
	sizeProp, err := expandComputeInstanceGroupSize(d.Get("size"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("size"); !tpgresource.IsEmptyValue(reflect.ValueOf(sizeProp)) && (ok || !reflect.DeepEqual(v, sizeProp)) {
		obj["size"] = sizeProp
	}
	regionProp, err := expandComputeInstanceGroupRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !tpgresource.IsEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}
	subnetworkProp, err := expandComputeInstanceGroupSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("subnetwork"); !tpgresource.IsEmptyValue(reflect.ValueOf(subnetworkProp)) && (ok || !reflect.DeepEqual(v, subnetworkProp)) {
		obj["subnetwork"] = subnetworkProp
	}

	return obj, nil
}

func expandComputeInstanceGroupName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupNamedPorts(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupPort(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupFingerprint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupZone(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupSelfLink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupSize(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupRegion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGroupSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
