package google

import (
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

const ServiceAccountAssetType string = "iam.googleapis.com/ServiceAccount"

func resourceConverterServiceAccount() ResourceConverter {
	return ResourceConverter{
		AssetType: ServiceAccountAssetType,
		Convert:   GetServiceAccountCaiObject,
	}
}

func GetServiceAccountCaiObject(d TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/projects/{{project}}/serviceAccounts/{{unique_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetServiceAccountApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ServiceAccountAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://iam.googleapis.com/$discovery/rest",
				DiscoveryName:        "ServiceAccount",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetServiceAccountApiObject(d TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	descriptionProp, err := expandServiceAccountDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	emailProp, err := expandServiceAccountDescription(d.Get("email"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("email"); !tpgresource.IsEmptyValue(reflect.ValueOf(emailProp)) && (ok || !reflect.DeepEqual(v, emailProp)) {
		obj["email"] = emailProp
	}

	displayNameProp, err := expandServiceAccountDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	nameProp, err := expandServiceAccountName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	disabledProp, err := expandServiceAccountDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	uniqueIdProp, err := expandServiceAccountUniqueId(d.Get("unique_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("unique_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(uniqueIdProp)) && (ok || !reflect.DeepEqual(v, uniqueIdProp)) {
		obj["uniqueId"] = uniqueIdProp
	}

	projectProp, err := expandServiceAccountProject(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectProp)) && (ok || !reflect.DeepEqual(v, projectProp)) {
		obj["projectId"] = projectProp
	}

	accountIdProp, err := expandServiceAccountId(d.Get("account_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("account_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(accountIdProp)) && (ok || !reflect.DeepEqual(v, accountIdProp)) {
		accountId := accountIdProp
		if _, ok := obj["email"]; !ok {
			// Generating email when the service account is being created (email not present)
			obj["email"] = fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
		}
	}
	return obj, nil
}

func expandServiceAccountId(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountDescription(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountDisplayName(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountName(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountEmail(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountDisabled(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountUniqueId(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountProject(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
