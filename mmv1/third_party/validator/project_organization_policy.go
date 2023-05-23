package google

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConverterProjectOrgPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectOrgPolicyCaiObject,
		MergeCreateUpdate: MergeProjectOrgPolicy,
	}
}

func GetProjectOrgPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetProjectOrgPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Project",
			OrgPolicy: []*OrgPolicy{&obj},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func MergeProjectOrgPolicy(existing, incoming Asset) Asset {
	existing.OrgPolicy = append(existing.OrgPolicy, incoming.OrgPolicy...)
	return existing
}

func GetProjectOrgPolicyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (OrgPolicy, error) {

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

func expandListOrganizationPolicy(configured []interface{}) (*ListPolicy, error) {
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
			allowedValues = tpgresource.ConvertStringArr(values.List())
		}
	}

	if len(deny) > 0 {
		denyMap := deny[0].(map[string]interface{})
		all := denyMap["all"].(bool)
		values := denyMap["values"].(*schema.Set)

		if all {
			allValues = 0
		} else {
			deniedValues = tpgresource.ConvertStringArr(values.List())
		}
	}

	listPolicy := configured[0].(map[string]interface{})
	return &ListPolicy{
		AllValues:         ListPolicyAllValues(allValues),
		AllowedValues:     allowedValues,
		DeniedValues:      deniedValues,
		SuggestedValue:    listPolicy["suggested_value"].(string),
		InheritFromParent: listPolicy["inherit_from_parent"].(bool),
	}, nil
}

func expandRestoreOrganizationPolicy(configured []interface{}) (*RestoreDefault, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	restoreDefaultMap := configured[0].(map[string]interface{})
	defaultValue := restoreDefaultMap["default"].(bool)

	if defaultValue {
		return &RestoreDefault{}, nil
	}

	return &RestoreDefault{}, fmt.Errorf("Invalid value for restore_policy. Expecting default = true")
}

func expandBooleanOrganizationPolicy(configured []interface{}) *BooleanPolicy {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	booleanPolicy := configured[0].(map[string]interface{})
	return &BooleanPolicy{
		Enforced: booleanPolicy["enforced"].(bool),
	}
}

func canonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}
