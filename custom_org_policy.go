package google

import (
	"fmt"
	"strings"
)

func resourceConverterCustomOrgPolicy() ResourceConverter {
	return ResourceConverter{
		Convert:           GetCustomOrgPolicyCaiObject,
		MergeCreateUpdate: MergeCustomOrgPolicy,
	}
}

func GetCustomOrgPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {

	assetNamePattern, assetType, err := getAssetNameAndTypeFromParent(d.Get("parent").(string))
	if err != nil {
		return []Asset{}, err
	}

	name, err := assetName(d, config, assetNamePattern)
	if err != nil {
		return []Asset{}, err
	}

	if obj, err := GetCustomOrgPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name:            name,
			Type:            assetType,
			CustomOrgPolicy: []*CustomOrgPolicy{&obj},
		}}, nil
	} else {
		return []Asset{}, err
	}

}

func GetCustomOrgPolicyApiObject(d TerraformResourceData, config *Config) (CustomOrgPolicy, error) {
	spec, err := expandSpecCustomOrgPolicy(d.Get("spec").([]interface{}))
	if err != nil {
		return CustomOrgPolicy{}, err
	}

	return CustomOrgPolicy{
		Name: d.Get("name").(string),
		Spec: spec,
	}, nil
}

func MergeCustomOrgPolicy(existing, incoming Asset) Asset {
	existing.Resource = incoming.Resource
	return existing
}

func getAssetNameAndTypeFromParent(parent string) (assetName string, assetType string, err error) {
	const prefix = "cloudresourcemanager.googleapis.com/"
	if strings.Contains(parent, "projects") {
		return prefix + "projects/{{project_id}}", prefix + "Project", nil
	} else if strings.Contains(parent, "folders") {
		return prefix + "folders/{{folder_id}}", prefix + "Folder", nil
	} else if strings.Contains(parent, "organizations") {
		return prefix + "organizations/{{organization_id}}", prefix + "Organization", nil
	} else {
		return "", "", fmt.Errorf("Invalid parent address(%s) for an asset", parent)
	}
}

func expandSpecCustomOrgPolicy(configured []interface{}) (*Spec, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	specMap := configured[0].(map[string]interface{})

	policyRules, err := expandPolicyRulesSpec(specMap["rules"].([]interface{}))
	if err != nil {
		return &Spec{}, err
	}

	return &Spec{
		Etag:       specMap["etag"].(string),
		Rules:      policyRules,
		InheritFromParent: specMap["inherit_from_parent"].(bool),
		Reset: specMap["reset"].(bool),
	}, nil

}

func expandPolicyRulesSpec(configured []interface{}) ([]*PolicyRule, error) {
	if configured[0] == nil {
		return nil, nil
	}

	var policyRules []*PolicyRule
	for i := 0; i < len(configured); i++ {
		policyRule, err := expandPolicyRulePolicyRules(configured[i])
		if err != nil {
			return nil, err
		}
		policyRules = append(policyRules, policyRule)
	}

	return policyRules, nil

}

func expandPolicyRulePolicyRules(configured interface{}) (*PolicyRule, error) {
	policyRuleMap := configured.(map[string]interface{})

	values, err := expandValuesPolicyRule(policyRuleMap["values"].([]interface{}))
	if err != nil {
		return &PolicyRule{}, err
	}

	allowAll, err := convertStringToBool(policyRuleMap["allow_all"].(string))
	if err != nil {
		return &PolicyRule{}, err
	}

	denyAll, err := convertStringToBool(policyRuleMap["deny_all"].(string))
	if err != nil {
		return &PolicyRule{}, err
	}

	enforce, err := convertStringToBool(policyRuleMap["enforce"].(string))
	if err != nil {
		return &PolicyRule{}, err
	}

	condition, err := expandConditionPolicyRule(policyRuleMap["condition"].([]interface{}))
	if err != nil {
		return &PolicyRule{}, err
	}
	return &PolicyRule{
		Values:    values,
		AllowAll:  allowAll,
		DenyAll:   denyAll,
		Enforce:   enforce,
		Condition: condition,
	}, nil
}

func expandValuesPolicyRule(configured []interface{}) (*StringValues, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	valuesMap := configured[0].(map[string]interface{})
	return &StringValues{
		AllowedValues: convertInterfaceToStringArray(valuesMap["allowed_values"].([]interface{})),
		DeniedValues:  convertInterfaceToStringArray(valuesMap["denied_values"].([]interface{})),
	}, nil
}

func expandConditionPolicyRule(configured []interface{}) (*Expr, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	conditionMap := configured[0].(map[string]interface{})
	return &Expr{
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
