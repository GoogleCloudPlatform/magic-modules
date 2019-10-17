package google

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func expandContainerEnabledObject(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if val := reflect.ValueOf(v); !val.IsValid() || isEmptyValue(val) {
		return nil, nil
	}
	transformed := map[string]interface{}{
		"enabled": v,
	}
	return transformed, nil
}

func expandContainerClusterEnableLegacyAbac(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return expandContainerEnabledObject(v, d, config)
}

func expandContainerClusterEnableBinaryAuthorization(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return expandContainerEnabledObject(v, d, config)
}

func expandContainerMaxPodsConstraint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if val := reflect.ValueOf(v); !val.IsValid() || isEmptyValue(val) {
		return nil, nil
	}
	transformed := map[string]interface{}{
		"maxPodsPerNode": v,
	}
	return transformed, nil
}

func expandContainerClusterDefaultMaxPodsPerNode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return expandContainerMaxPodsConstraint(v, d, config)
}

func expandContainerNodePoolMaxPodsPerNode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return expandContainerMaxPodsConstraint(v, d, config)
}

func expandContainerClusterNetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	fv, err := ParseNetworkFieldValue(v.(string), d, config)
	if err != nil {
		return nil, err
	}
	return fv.RelativeLink(), nil
}

func expandContainerClusterSubnetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	fv, err := ParseNetworkFieldValue(v.(string), d, config)
	if err != nil {
		return nil, err
	}
	return fv.RelativeLink(), nil
}

func canonicalizeServiceScopesFromSet(scopesSet *schema.Set) (interface{}, error) {
	scopes := make([]string, scopesSet.Len())
	for i, scope := range scopesSet.List() {
		scopes[i] = canonicalizeServiceScope(scope.(string))
	}
	return scopes, nil
}

func expandContainerClusterNodeConfigOauthScopes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	scopesSet := v.(*schema.Set)
	return canonicalizeServiceScopesFromSet(scopesSet)
}

func expandContainerNodePoolNodeConfigOauthScopes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	scopesSet := v.(*schema.Set)
	return canonicalizeServiceScopesFromSet(scopesSet)
}
