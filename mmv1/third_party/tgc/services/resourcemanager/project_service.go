package resourcemanager

import (
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ServiceUsageAssetType string = "serviceusage.googleapis.com/Service"

func ResourceConverterServiceUsage() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ServiceUsageAssetType,
		Convert:   GetServiceUsageCaiObject,
	}
}

func GetServiceUsageCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//serviceusage.googleapis.com/projects/{{project}}/services/{{service}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetServiceUsageApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ServiceUsageAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/serviceusage/v1/rest",
				DiscoveryName:        "Service",
				Data:                 obj,
			}},
		}, nil
	}
	return []cai.Asset{}, err
}

func GetServiceUsageApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	parentProjectProp, err := expandServiceUsageParentProject(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	}
	obj["parent"] = parentProjectProp

	serviceNameProp, err := expandServiceUsageServiceName(d.Get("service"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceNameProp)) && (ok || !reflect.DeepEqual(v, serviceNameProp)) {
		obj["name"] = serviceNameProp
	}

	obj["state"] = "ENABLED"

	return obj, nil
}

func expandServiceUsageParentProject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil || v.(string) == "" {
		// It does not try to construct anything from empty.
		return "", nil
	}
	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	return fmt.Sprintf("projects/%s", v.(string)), nil
}

func expandServiceUsageServiceName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
