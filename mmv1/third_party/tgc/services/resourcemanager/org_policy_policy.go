package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterOrgPolicyPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		Convert:           GetV2OrgPoliciesCaiObject,
		MergeCreateUpdate: MergeV2OrgPolicies,
	}
}

func GetV2OrgPoliciesCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	assetNamePattern, assetType, err := getAssetNameAndTypeFromParent(d.Get("parent").(string))
	if err != nil {
		return []cai.Asset{}, err
	}

	name, err := cai.AssetName(d, config, assetNamePattern)
	if err != nil {
		return []cai.Asset{}, err
	}

	if obj, err := GetV2OrgPoliciesApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name:          name,
			Type:          assetType,
			V2OrgPolicies: []*cai.V2OrgPolicies{&obj},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}

}

func GetV2OrgPoliciesApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.V2OrgPolicies, error) {
	spec, err := expandSpecV2OrgPolicies(d.Get("spec").([]interface{}))
	if err != nil {
		return cai.V2OrgPolicies{}, err
	}

	return cai.V2OrgPolicies{
		Name:       d.Get("name").(string),
		PolicySpec: spec,
	}, nil
}

func MergeV2OrgPolicies(existing, incoming cai.Asset) cai.Asset {
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

func expandSpecV2OrgPolicies(configured []interface{}) (*cai.PolicySpec, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	specMap := configured[0].(map[string]interface{})

	policyRules, err := expandPolicyRulesSpec(specMap["rules"].([]interface{}))
	if err != nil {
		return &cai.PolicySpec{}, err
	}

	return &cai.PolicySpec{
		Etag:              specMap["etag"].(string),
		PolicyRules:       policyRules,
		InheritFromParent: specMap["inherit_from_parent"].(bool),
		Reset:             specMap["reset"].(bool),
	}, nil

}

func expandPolicyRulesSpec(configured []interface{}) ([]*cai.PolicyRule, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	var policyRules []*cai.PolicyRule
	for i := 0; i < len(configured); i++ {
		policyRule, err := expandPolicyRulePolicyRules(configured[i])
		if err != nil {
			return nil, err
		}
		policyRules = append(policyRules, policyRule)
	}

	return policyRules, nil

}

func expandPolicyRulePolicyRules(configured interface{}) (*cai.PolicyRule, error) {
	policyRuleMap := configured.(map[string]interface{})

	values, err := expandValuesPolicyRule(policyRuleMap["values"].([]interface{}))
	if err != nil {
		return &cai.PolicyRule{}, err
	}

	allowAll, err := convertStringToBool(policyRuleMap["allow_all"].(string))
	if err != nil {
		return &cai.PolicyRule{}, err
	}

	denyAll, err := convertStringToBool(policyRuleMap["deny_all"].(string))
	if err != nil {
		return &cai.PolicyRule{}, err
	}

	enforce, err := convertStringToBool(policyRuleMap["enforce"].(string))
	if err != nil {
		return &cai.PolicyRule{}, err
	}

	condition, err := expandConditionPolicyRule(policyRuleMap["condition"].([]interface{}))
	if err != nil {
		return &cai.PolicyRule{}, err
	}
	return &cai.PolicyRule{
		Values:    values,
		AllowAll:  allowAll,
		DenyAll:   denyAll,
		Enforce:   enforce,
		Condition: condition,
	}, nil
}

func expandValuesPolicyRule(configured []interface{}) (*cai.StringValues, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	valuesMap := configured[0].(map[string]interface{})
	return &cai.StringValues{
		AllowedValues: cai.ConvertInterfaceToStringArray(valuesMap["allowed_values"].([]interface{})),
		DeniedValues:  cai.ConvertInterfaceToStringArray(valuesMap["denied_values"].([]interface{})),
	}, nil
}

func expandConditionPolicyRule(configured []interface{}) (*cai.Expr, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}
	conditionMap := configured[0].(map[string]interface{})
	return &cai.Expr{
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
