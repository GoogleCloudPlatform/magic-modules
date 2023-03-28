package google

import (
	"fmt"
	"strings"
)

func resourceConverterOrgPolicyPolicy() ResourceConverter {
	return ResourceConverter{
		Convert:           GetOrgPolicyPolicyCaiObject,
		MergeCreateUpdate: MergeOrgPolicyPolicy,
	}
}

func GetOrgPolicyPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	assetNamePattern, assetType, err := getAssetNameAndTypeFromParent(d.Get("parent").(string))
	if err != nil {
		return []Asset{}, err
	}

	name, err := assetName(d, config, assetNamePattern)
	if err != nil {
		return []Asset{}, err
	}

	if obj, err := GetOrgPolicyPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name:            name,
			Type:            assetType,
			OrgPolicyPolicy: []*OrgPolicyPolicy{&obj},
		}}, nil
	} else {
		return []Asset{}, err
	}

}

func GetOrgPolicyPolicyApiObject(d TerraformResourceData, config *Config) (OrgPolicyPolicy, error) {
	spec, err := expandSpecOrgPolicyPolicy(d.Get("spec").([]interface{}))
	if err != nil {
		return OrgPolicyPolicy{}, err
	}

	return OrgPolicyPolicy{
		Name: d.Get("name").(string),
		PolicySpec: spec,
	}, nil
}

func MergeOrgPolicyPolicy(existing, incoming Asset) Asset {
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

func expandSpecOrgPolicyPolicy(configured []interface{}) (*PolicySpec, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	specMap := configured[0].(map[string]interface{})

	policyRules, err := expandPolicyRulesSpec(specMap["rules"].([]interface{}))
	if err != nil {
		return &PolicySpec{}, err
	}

	return &PolicySpec{
		Etag:              specMap["etag"].(string),
		PolicyRules:       policyRules,
		InheritFromParent: specMap["inherit_from_parent"].(bool),
		Reset:             specMap["reset"].(bool),
	}, nil

}

func expandPolicyRulesSpec(configured []interface{}) ([]*PolicyRule, error) {
	if len(configured) == 0 || configured[0] == nil {
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
