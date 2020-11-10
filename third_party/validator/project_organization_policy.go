package google

import (
	"fmt"
	"strings"

	"google3/third_party/golang/hashicorp/terraform_plugin_sdk/helper/schema/schema"
)

func GetProjectOrgPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetProjectOrgPolicyApiObject(d, config); err == nil {
		return Asset{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Project",
			OrgPolicy: []*OrgPolicy{&obj},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetProjectOrgPolicyApiObject(d TerraformResourceData, config *Config) (OrgPolicy, error) {

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return OrgPolicy{}, err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return OrgPolicy{}, err
	}

	policy := OrgPolicy{
		Constraint:     canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
		ListPolicy:     listPolicy,
		RestoreDefault: restoreDefault,
	}

	return policy, nil
}

func expandListOrganizationPolicy(configured []interface{}) (*Policy_ListPolicy, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	listPolicyMap := configured[0].(map[string]interface{})

	allow := listPolicyMap["allow"].([]interface{})
	deny := listPolicyMap["deny"].([]interface{})

	var allValues int32
	var allowedValues []string
	var deniedValues []string
	if len(allow) > 0 {
		allowMap := allow[0].(map[string]interface{})
		all := allowMap["all"].(bool)
		values := allowMap["values"].(*schema.Set)

		if all {
			allValues = 1
		} else {
			allowedValues = convertStringArr(values.List())
		}
	}

	if len(deny) > 0 {
		denyMap := deny[0].(map[string]interface{})
		all := denyMap["all"].(bool)
		values := denyMap["values"].(*schema.Set)

		if all {
			allValues = 0
		} else {
			deniedValues = convertStringArr(values.List())
		}
	}

	listPolicy := configured[0].(map[string]interface{})
	return &Policy_ListPolicy{
		AllValues:         Policy_ListPolicy_AllValues(allValues),
		AllowedValues:     allowedValues,
		DeniedValues:      deniedValues,
		SuggestedValue:    listPolicy["suggested_value"].(string),
		InheritFromParent: listPolicy["inherit_from_parent"].(bool),
	}, nil
}

func expandRestoreOrganizationPolicy(configured []interface{}) (*Policy_RestoreDefault, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	restoreDefaultMap := configured[0].(map[string]interface{})
	defaultValue := restoreDefaultMap["default"].(bool)

	if defaultValue {
		return &Policy_RestoreDefault{}, nil
	}

	return &Policy_RestoreDefault{}, fmt.Errorf("Invalid value for restore_policy. Expecting default = true")
}

func expandBooleanOrganizationPolicy(configured []interface{}) *Policy_BooleanPolicy {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	booleanPolicy := configured[0].(map[string]interface{})
	return &Policy_BooleanPolicy{
		Enforced: booleanPolicy["enforced"].(bool),
	}
}

func canonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}
