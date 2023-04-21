package google

func resourceConverterFolderOrgPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderOrgPolicyCaiObject,
		MergeCreateUpdate: MergeFolderOrgPolicy,
	}
}

func GetFolderOrgPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/{{folder}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetFolderOrgPolicyApiObject(d, config); err == nil {
		return []Asset{{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Folder",
			OrgPolicy: []*OrgPolicy{&obj},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func MergeFolderOrgPolicy(existing, incoming Asset) Asset {
	existing.OrgPolicy = append(existing.OrgPolicy, incoming.OrgPolicy...)
	return existing
}

func GetFolderOrgPolicyApiObject(d TerraformResourceData, config *Config) (OrgPolicy, error) {

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
