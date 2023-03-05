package google

func (u *PrivatecaCaPoolIamUpdater) SetProject(project string) {
	u.project = project
}

func (u *PrivatecaCaPoolIamUpdater) SetLocation(location string) {
	u.location = location
}

func (u *PrivatecaCaPoolIamUpdater) SetCaPool(caPool string) {
	u.caPool = caPool
}

func (u *PrivatecaCaPoolIamUpdater) SetTerraformResourceData(d TerraformResourceData) {
	u.d = d
}
