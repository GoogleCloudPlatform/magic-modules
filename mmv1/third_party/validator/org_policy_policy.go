package google

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterOrgPolicyPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		Convert:           GetV2OrgPoliciesCaiObject,
		MergeCreateUpdate: MergeV2OrgPolicies,
	}
}

func GetV2OrgPoliciesCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	assetNamePattern, assetType, err := getAssetNameAndTypeFromParent(d.Get("parent").(string))
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	name, err := tpgresource.AssetName(d, config, assetNamePattern)
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	if obj, err := GetV2OrgPoliciesApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name:          name,
			Type:          assetType,
			V2OrgPolicies: []*tpgresource.V2OrgPolicies{&obj},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}

}

func GetV2OrgPoliciesApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.V2OrgPolicies, error) {
	spec, err := expandSpecV2OrgPolicies(d.Get("spec").([]interface{}))
	if err != nil {
		return tpgresource.V2OrgPolicies{}, err
	}

	return tpgresource.V2OrgPolicies{
		Name:       d.Get("name").(string),
		PolicySpec: spec,
	}, nil
}

func MergeV2OrgPolicies(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.Resource = incoming.Resource
	return existing
}

func getAssetNameAndTypeFromParent(parent string) (assetName string, assetType string, err error) {
	const prefix = "cloudresourcemanager.googleapis.com/"
	if strings.Contains(parent, "projects") {
		return "//" + prefix + parent, prefix + "Project", nil
	} else if strings.Contains(parent, "folders") {
		return "//" + prefix + parent, prefix + "Folder", nil
	} else if strings.Contains(parent, "organizations") {
		return "//" + prefix + parent, prefix + "Organization", nil
	} else {
		return "", "", fmt.Errorf("Invalid parent address(%s) for an asset", parent)
	}
}

func expandSpecV2OrgPolicies(configured []interface{}) (*tpgresource.PolicySpec, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	specMap := configured[0].(map[string]interface{})

	policyRules, err := expandPolicyRulesSpec(specMap["rules"].([]interface{}))
	if err != nil {
		return &tpgresource.PolicySpec{}, err
	}

	return &tpgresource.PolicySpec{
		Etag:              specMap["etag"].(string),
		PolicyRules:       policyRules,
		InheritFromParent: specMap["inherit_from_parent"].(bool),
		Reset:             specMap["reset"].(bool),
	}, nil

}

func expandPolicyRulesSpec(configured []interface{}) ([]*tpgresource.PolicyRule, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	var policyRules []*tpgresource.PolicyRule
	for i := 0; i < len(configured); i++ {
		policyRule, err := expandPolicyRulePolicyRules(configured[i])
		if err != nil {
			return nil, err
		}
		policyRules = append(policyRules, policyRule)
	}

	return policyRules, nil

}

func expandPolicyRulePolicyRules(configured interface{}) (*tpgresource.PolicyRule, error) {
	policyRuleMap := configured.(map[string]interface{})

	values, err := expandValuesPolicyRule(policyRuleMap["values"].([]interface{}))
	if err != nil {
		return &tpgresource.PolicyRule{}, err
	}

	allowAll, err := convertStringToBool(policyRuleMap["allow_all"].(string))
	if err != nil {
		return &tpgresource.PolicyRule{}, err
	}

	denyAll, err := convertStringToBool(policyRuleMap["deny_all"].(string))
	if err != nil {
		return &tpgresource.PolicyRule{}, err
	}

	enforce, err := convertStringToBool(policyRuleMap["enforce"].(string))
	if err != nil {
		return &tpgresource.PolicyRule{}, err
	}

	condition, err := expandConditionPolicyRule(policyRuleMap["condition"].([]interface{}))
	if err != nil {
		return &tpgresource.PolicyRule{}, err
	}
	return &tpgresource.PolicyRule{
		Values:    values,
		AllowAll:  allowAll,
		DenyAll:   denyAll,
		Enforce:   enforce,
		Condition: condition,
	}, nil
}

func expandValuesPolicyRule(configured []interface{}) (*tpgresource.StringValues, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	valuesMap := configured[0].(map[string]interface{})
	return &tpgresource.StringValues{
		AllowedValues: convertInterfaceToStringArray(valuesMap["allowed_values"].([]interface{})),
		DeniedValues:  convertInterfaceToStringArray(valuesMap["denied_values"].([]interface{})),
	}, nil
}

func expandConditionPolicyRule(configured []interface{}) (*tpgresource.Expr, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	conditionMap := configured[0].(map[string]interface{})
	return &tpgresource.Expr{
		Expression:  conditionMap["expression"].(string),
		Title:       conditionMap["title"].(string),
		Description: conditionMap["description"].(string),
		Location:    conditionMap["location"].(string),
	}, nil
}

func convertStringToBool(val string) (bool, error) {
	if (val == "false") || (val == "FALSE") || (val == "") {
		return false, nil
	} else if (val == "true") || (val == "TRUE") {
		return true, nil
	}

	return false, fmt.Errorf("Invalid value for a boolean field: %s", val)
}

func convertInterfaceToStringArray(values []interface{}) []string {
	stringArray := make([]string, len(values))
	for i, v := range values {
		stringArray[i] = v.(string)
	}
	return stringArray
}
