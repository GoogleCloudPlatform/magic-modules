package compute

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ComputeTargetPoolAssetType string = "compute.googleapis.com/TargetPool"

func ResourceConverterComputeTargetPool() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ComputeTargetPoolAssetType,
		Convert:   GetComputeTargetPoolCaiObject,
	}
}

func GetComputeTargetPoolCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/regions/{{region}}/targetPools/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetComputeTargetPoolApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ComputeTargetPoolAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "TargetPool",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetComputeTargetPoolApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandComputeTargetPoolName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	descriptionProp, err := expandComputeTargetPoolDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	regionProp, err := expandComputeTargetPoolRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !tpgresource.IsEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}

	healthChecksProp, err := expandComputeTargetPoolHealthChecks(d.Get("health_checks"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("health_checks"); !tpgresource.IsEmptyValue(reflect.ValueOf(healthChecksProp)) && (ok || !reflect.DeepEqual(v, healthChecksProp)) {
		obj["healthChecks"] = healthChecksProp
	}

	instancesProp, err := expandComputeTargetPoolInstances(d.Get("instances"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("instances"); !tpgresource.IsEmptyValue(reflect.ValueOf(instancesProp)) && (ok || !reflect.DeepEqual(v, instancesProp)) {
		obj["instances"] = instancesProp
	}

	sessionAffinityProp, err := expandComputeTargetPoolSessionAffinity(d.Get("session_affinity"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("session_affinity"); !tpgresource.IsEmptyValue(reflect.ValueOf(sessionAffinityProp)) && (ok || !reflect.DeepEqual(v, sessionAffinityProp)) {
		obj["sessionAffinity"] = sessionAffinityProp
	}

	failoverRatioProp, err := expandComputeTargetPoolFailoverRatio(d.Get("failover_ratio"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("failover_ratio"); !tpgresource.IsEmptyValue(reflect.ValueOf(failoverRatioProp)) && (ok || !reflect.DeepEqual(v, failoverRatioProp)) {
		obj["failoverRatio"] = failoverRatioProp
	}

	backupPoolProp, err := expandComputeTargetPoolBackupPool(d.Get("backup_pool"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("backup_pool"); !tpgresource.IsEmptyValue(reflect.ValueOf(backupPoolProp)) && (ok || !reflect.DeepEqual(v, backupPoolProp)) {
		obj["backupPool"] = backupPoolProp
	}

	selfLinkProp, err := expandComputeTargetPoolSelfLink(d.Get("self_link"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("self_link"); !tpgresource.IsEmptyValue(reflect.ValueOf(selfLinkProp)) && (ok || !reflect.DeepEqual(v, selfLinkProp)) {
		obj["selfLink"] = selfLinkProp
	}

	securityPolicyProp, err := expandComputeTargetPoolSecurityPolicy(d.Get("security_policy"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("security_policy"); !tpgresource.IsEmptyValue(reflect.ValueOf(securityPolicyProp)) && (ok || !reflect.DeepEqual(v, securityPolicyProp)) {
		obj["securityPolicy"] = securityPolicyProp
	}

	return obj, nil
}

func expandComputeTargetPoolName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolRegion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolHealthChecks(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeTargetPoolSessionAffinity(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolFailoverRatio(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolBackupPool(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolSelfLink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetPoolSecurityPolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
