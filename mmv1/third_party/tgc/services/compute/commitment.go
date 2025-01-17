package compute

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const commitmentAssetType string = "compute.googleapis.com/Commitment"

func ResourceConverterCommitment() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: commitmentAssetType,
		Convert:   GetCommitmentCaiObject,
	}
}

func GetCommitmentCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/regions/{{region}}/commitments")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetCommitmentApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: commitmentAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Commitment",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetCommitmentApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandCommitmentName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	planProp, err := expandCommitmentPlan(d.Get("plan"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("plan"); !tpgresource.IsEmptyValue(reflect.ValueOf(planProp)) && (ok || !reflect.DeepEqual(v, planProp)) {
		obj["plan"] = planProp
	}

	descriptionProp, err := expandCommitmentDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	resourcesProp, err := expandCommitmentResources(d.Get("resources"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("resources"); !tpgresource.IsEmptyValue(reflect.ValueOf(resourcesProp)) && (ok || !reflect.DeepEqual(v, resourcesProp)) {
		obj["resources"] = resourcesProp
	}

	typeProp, err := expandCommitmentType(d.Get("type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}

	categoryProp, err := expandCommitmentCategory(d.Get("category"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("category"); !tpgresource.IsEmptyValue(reflect.ValueOf(categoryProp)) && (ok || !reflect.DeepEqual(v, categoryProp)) {
		obj["category"] = categoryProp
	}

	licenseResourceProp, err := expandCommitmentLicenseResource(d.Get("license_resource"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("license_resource"); !tpgresource.IsEmptyValue(reflect.ValueOf(licenseResourceProp)) && (ok || !reflect.DeepEqual(v, licenseResourceProp)) {
		obj["licenseResource"] = licenseResourceProp
	}

	autoRenewProp, err := expandCommitmentAutoRenew(d.Get("auto_renew"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("auto_renew"); !tpgresource.IsEmptyValue(reflect.ValueOf(autoRenewProp)) && (ok || !reflect.DeepEqual(v, autoRenewProp)) {
		obj["autoRenew"] = autoRenewProp
	}

	regionProp, err := expandCommitmentRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !tpgresource.IsEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}

	projectProp, err := expandCommitmentProject(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectProp)) && (ok || !reflect.DeepEqual(v, projectProp)) {
		obj["project"] = projectProp
	}

	idProp, err := expandCommitmentId(d.Get("commitment_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("commitment_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	idIdentifierProp, err := expandCommitmentIdIdentifier(d.Get("id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idIdentifierProp)) && (ok || !reflect.DeepEqual(v, idIdentifierProp)) {
		obj["id"] = idIdentifierProp
	}

	creationTimestampProp, err := expandCommitmentCreationTimestamp(d.Get("creation_timestamp"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("creation_timestamp"); !tpgresource.IsEmptyValue(reflect.ValueOf(idIdentifierProp)) && (ok || !reflect.DeepEqual(v, idIdentifierProp)) {
		obj["creationTimestamp"] = creationTimestampProp
	}

	statusProp, err := expandCommitmentStatus(d.Get("status"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("status"); !tpgresource.IsEmptyValue(reflect.ValueOf(statusProp)) && (ok || !reflect.DeepEqual(v, statusProp)) {
		obj["status"] = statusProp
	}

	statusMessageProp, err := expandCommitmentStatusMessage(d.Get("status_message"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("status_message"); !tpgresource.IsEmptyValue(reflect.ValueOf(statusMessageProp)) && (ok || !reflect.DeepEqual(v, statusMessageProp)) {
		obj["statusMessage"] = statusMessageProp
	}

	startTimestampProp, err := expandCommitmentStartTimestamp(d.Get("start_timestamp"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("start_timestamp"); !tpgresource.IsEmptyValue(reflect.ValueOf(statusMessageProp)) && (ok || !reflect.DeepEqual(v, statusMessageProp)) {
		obj["startTimestamp"] = startTimestampProp
	}

	endTimestampProp, err := expandCommitmentEndTimestamp(d.Get("end_timestamp"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("end_timestamp"); !tpgresource.IsEmptyValue(reflect.ValueOf(endTimestampProp)) && (ok || !reflect.DeepEqual(v, endTimestampProp)) {
		obj["endTimestamp"] = endTimestampProp
	}

	selfLinkProp, err := expandCommitmentSelfLink(d.Get("self_link"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("self_link"); !tpgresource.IsEmptyValue(reflect.ValueOf(endTimestampProp)) && (ok || !reflect.DeepEqual(v, endTimestampProp)) {
		obj["selfLink"] = selfLinkProp
	}

	return obj, nil
}

func expandCommitmentLicenseResource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedlicense, err := expandCommitmentLicense(original["license"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedlicense); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["license"] = transformedlicense
	}

	transformedAmount, err := expandCommitmentAmount(original["amount"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAmount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["amount"] = transformedAmount
	}

	transformedCoresPerLicense, err := expandCommitmentCoresPerLicense(original["cores_per_license"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCoresPerLicense); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["coresPerLicense"] = transformedCoresPerLicense
	}

	return transformed, nil
}

func expandCommitmentSelfLink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentEndTimestamp(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentStartTimestamp(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentStatusMessage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentStatus(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentCreationTimestamp(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentIdIdentifier(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/commitments/{{name}}")
	if err != nil {
		return nil, err
	}

	return v, nil
}

func expandCommitmentId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentProject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentRegion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentAutoRenew(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentCoresPerLicense(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentLicense(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentCategory(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentResources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedType, err := expandCommitmentType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		transformedAmount, err := expandCommitmentAmount(original["amount"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAmount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["amount"] = transformedAmount
		}

		transformedAcceleratorType, err := expandCommitmentAcceleratorType(original["accelerator_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorType"] = transformedAcceleratorType
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandCommitmentAcceleratorType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentAmount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentPlan(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandCommitmentName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
