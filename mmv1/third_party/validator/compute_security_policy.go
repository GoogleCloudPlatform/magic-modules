package google

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ComputeSecurityPolicyAssetType string = "compute.googleapis.com/SecurityPolicy"

func resourceConverterComputeSecurityPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType: ComputeSecurityPolicyAssetType,
		Convert:   GetComputeSecurityPolicyCaiObject,
	}
}

func GetComputeSecurityPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//compute.googleapis.com/projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetComputeSecurityPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ComputeSecurityPolicyAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "SecurityPolicy",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetComputeSecurityPolicyApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandComputeSecurityPolicyName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	rulesProp, err := expandComputeSecurityPolicyRules(d.Get("rule"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("rule"); !isEmptyValue(reflect.ValueOf(rulesProp)) && (ok || !reflect.DeepEqual(v, rulesProp)) {
		obj["rule"] = rulesProp
	}

	return obj, nil
}

func expandComputeSecurityPolicyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRules(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDescription, err := expandComputeSecurityPolicyRulesDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !isEmptyValue(val) {
			transformed["description"] = transformedDescription
		}

		transformedPriority, err := expandComputeSecurityPolicyRulesPriority(original["priority"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPriority); val.IsValid() && !isEmptyValue(val) {
			transformed["priority"] = transformedPriority
		}

		transformedAction, err := expandComputeSecurityPolicyRulesAction(original["action"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAction); val.IsValid() && !isEmptyValue(val) {
			transformed["action"] = transformedAction
		}

		transformedPreview, err := expandComputeSecurityPolicyRulesPreview(original["preview"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPreview); val.IsValid() && !isEmptyValue(val) {
			transformed["preview"] = transformedPreview
		}

		transformedMatch, err := expandComputeSecurityPolicyRulesMatch(original["match"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMatch); val.IsValid() && !isEmptyValue(val) {
			transformed["match"] = transformedMatch
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeSecurityPolicyRulesDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesPriority(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesAction(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesPreview(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatch(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDescription, err := expandComputeSecurityPolicyRulesMatchDescription(original["description"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !isEmptyValue(val) {
		transformed["description"] = transformedDescription
	}

	transformedExpr, err := expandComputeSecurityPolicyRulesMatchExpr(original["expr"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpr); val.IsValid() && !isEmptyValue(val) {
		transformed["expr"] = transformedExpr
	}

	transformedVersionedExpr, err := expandComputeSecurityPolicyRulesMatchVersionedExpr(original["versioned_expr"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedVersionedExpr); val.IsValid() && !isEmptyValue(val) {
		transformed["versionedExpr"] = transformedVersionedExpr
	}

	transformedConfig, err := expandComputeSecurityPolicyRulesMatchConfig(original["config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedConfig); val.IsValid() && !isEmptyValue(val) {
		transformed["config"] = transformedConfig
	}

	return transformed, nil
}

func expandComputeSecurityPolicyRulesMatchDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchExpr(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedExpression, err := expandComputeSecurityPolicyRulesMatchExprExpression(original["expression"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpression); val.IsValid() && !isEmptyValue(val) {
		transformed["expression"] = transformedExpression
	}

	transformedTitle, err := expandComputeSecurityPolicyRulesMatchExprTitle(original["title"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTitle); val.IsValid() && !isEmptyValue(val) {
		transformed["title"] = transformedTitle
	}

	transformedDescription, err := expandComputeSecurityPolicyRulesMatchExprDescription(original["description"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !isEmptyValue(val) {
		transformed["description"] = transformedDescription
	}

	transformedLocation, err := expandComputeSecurityPolicyRulesMatchExprLocation(original["location"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocation); val.IsValid() && !isEmptyValue(val) {
		transformed["location"] = transformedLocation
	}

	return transformed, nil
}

func expandComputeSecurityPolicyRulesMatchExprExpression(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchExprTitle(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchExprDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchExprLocation(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchVersionedExpr(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSecurityPolicyRulesMatchConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSrcIpRanges, err := expandComputeSecurityPolicyRulesMatchConfigSrcIpRanges(original["src_ip_ranges"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSrcIpRanges); val.IsValid() && !isEmptyValue(val) {
		transformed["srcIpRanges"] = transformedSrcIpRanges
	}

	return transformed, nil
}

func expandComputeSecurityPolicyRulesMatchConfigSrcIpRanges(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}
