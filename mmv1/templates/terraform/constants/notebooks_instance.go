var NotebooksInstanceProvidedScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
}

func NotebooksInstanceScopesDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
  old, new := d.GetChange("service_account_scopes")
	oldValue := old.([]interface{})
	newValue := new.([]interface{})
	oldValueList := []string{}
	newValueList := []string{}

	for _, item := range oldValue {
		oldValueList = append(oldValueList,item.(string))
	}

	for _, item := range newValue {
		newValueList = append(newValueList,item.(string))
	}
	newValueList= append(newValueList,NotebooksInstanceProvidedScopes...)

	sort.Strings(oldValueList)
	sort.Strings(newValueList)
	if reflect.DeepEqual(oldValueList, newValueList) {
		return true
	}
	return false
}

func NotebooksInstanceKmsDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if strings.HasPrefix(old, new) {
		return true
	}
	return false
}
