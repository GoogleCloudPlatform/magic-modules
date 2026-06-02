package gkehub

func alsoExpandEmptyBundlesInMap(c *Client, f map[string]FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles, res *FeatureMembership) (map[string]any, error) {
	if len(f) == 0 {
		return nil, nil
	}

	items := make(map[string]any)
	for k, v := range f {
		i, err := alsoExpandEmptyBundles(c, &v, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}
	return items, nil
}

func alsoExpandEmptyBundles(c *Client, f *FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles, res *FeatureMembership) (map[string]any, error) {
	m := make(map[string]any)
	if v := f.ExemptedNamespaces; v != nil {
		m["exemptedNamespaces"] = v
	}
	return m, nil
}
