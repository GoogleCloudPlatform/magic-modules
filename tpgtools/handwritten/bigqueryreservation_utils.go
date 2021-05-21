package google

func getProjectField(d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("project")
	if !ok {
		return "", nil
	}

	return res.(string), nil
}
