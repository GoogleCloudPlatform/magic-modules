package google

import (
	"google.golang.org/api/orgpolicy/v2"
)

func resourceConverterOrgPolicyV2CustomConstraint() ResourceConverter {
	return ResourceConverter{
		AssetType: "orgpolicy.googleapis.com/CustomConstraint",
		Convert:   GetOrgPolicyV2CustomConstraintCaiObject,
	}
}

func GetOrgPolicyV2CustomConstraintCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//orgpolicy.googleapis.com/v2/organizations/{{org_id}}/customConstraints/{{custom_constraint}}")

	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetOrgPolicyV2CustomConstraintApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "orgpolicy.googleapis.com/CustomConstraint",
			Resource: &AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://orgpolicy.googleapis.com/$discovery/rest?version=v2",
				DiscoveryName:        "GoogleCloudOrgpolicyV2CustomConstraint",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetOrgPolicyV2CustomConstraintApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {

	// ASK-In what cases return err ?

	// Get the custom constraint name
	customConstraint := d.Get("name").(string)

	// Create a custom constraint, setting the name.
	cc := &orgpolicy.GoogleCloudOrgpolicyV2CustomConstraint{
		Name: customConstraint,
	}

	if v, ok := d.GetOk("description"); ok {
		cc.Description = v.(string)
	}

	if v, ok := d.GetOk("display_name"); ok {
		cc.DisplayName = v.(string)
	}

	if v, ok := d.GetOk("condition"); ok {
		cc.Condition = v.(string)
	}

	if v, ok := d.GetOk("resource_types"); ok {
		cc.ResourceTypes = expandResourceTypes(v.([]interface{}))
	}

	if v, ok := d.GetOk("action_type"); ok {
		cc.ActionType = v.(string)
	}

	if v, ok := d.GetOk("method_types"); ok {
		cc.MethodTypes = expandMethodTypes(v.([]interface{}))
	}

	m, err := jsonMap(cc)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func expandResourceTypes(v interface{}) []string {
	if v == nil {
		return nil
	}
	vs := v.([]interface{})

	if len(vs) < 1 {
		return nil
	}

	countResourceTypes := len(vs)

	var resourceTypes = make([]string, countResourceTypes)

	for i := 0; i < countResourceTypes; i++ {
		resourceTypes[i] = vs[i].(string)
	}

	return resourceTypes
}

func expandMethodTypes(v interface{}) []string {
	if v == nil {
		return nil
	}

	vs := v.([]interface{})

	if len(vs) < 1 {
		return nil
	}

	countMethodTypes := len(vs)

	var methodTypes = make([]string, countMethodTypes)

	for i := 0; i < countMethodTypes; i++ {
		methodTypes[i] = vs[i].(string)
	}

	return methodTypes
}
