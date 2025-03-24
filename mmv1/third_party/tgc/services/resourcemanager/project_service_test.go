package resourcemanager

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
)

func TestParseServiceForExistingProject(t *testing.T) {
	cases := []struct {
		name                  string
		data                  tpgresource.TerraformResourceData
		expectedType          string
		expectedAssetName     string
		expectedParentProject string
		expectedService       string
		expectedState         string
	}{
		{
			name: "resource has service and project",
			data: &mockTerraformResourceData{
				m: map[string]interface{}{
					"project": "test-project",
					"service": "iamcredentials.googleapis.com",
				},
			},
			expectedType:          "serviceusage.googleapis.com/Service",
			expectedAssetName:     "//serviceusage.googleapis.com/projects/test-project/services/iamcredentials.googleapis.com",
			expectedParentProject: "projects/test-project",
			expectedServiceName:   "iamcredentials.googleapis.com",
			expectedState:         "ENABLED",
		},
		{
			name: "resource has service but missing project",
			data: &mockTerraformResourceData{
				m: map[string]interface{}{
					"service": "iamcredentials.googleapis.com",
				},
			},
			expectedType:          "serviceusage.googleapis.com/Service",
			expectedAssetName:     "//serviceusage.googleapis.com/projects/default_project/services/iamcredentials.googleapis.com",
			expectedParentProject: "",
			expectedServiceName:   "iamcredentials.googleapis.com",
			expectedState:         "ENABLED",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{}
			config.Project = "default_project"
			asset, err := GetServiceUsageCaiObject(tt.data, config)
			if err != nil {
				t.Errorf("%s failed to convert project_service resource:%q", tt.name, err)
			}
			if asset.Type != tt.expectedType {
				t.Errorf("Case: %s. Error converting asset type. Expected: %q got: %q", tt.name, tt.expectedType, asset.Type)
			}
			if asset.Name != tt.expectedAssetName {
				t.Errorf("Case: %s. Error converting asset name. Expected: %q got: %q", tt.name, tt.expectedAssetName, asset.Name)
			}
			if asset.Resource.Data["parent"] != tt.expectedParentProject {
				t.Errorf("Case: %s. Error converting asset parent project. Expected: %q got: %q", tt.name, tt.expectedParentProject, asset.Resource.Data["parent"])
			}
			if asset.Resource.Data["name"] != tt.expectedServiceName {
				t.Errorf("Case: %s. Error converting asset service name. Expected: %q got: %q", tt.name, tt.expectedServiceName, asset.Resource.Data["name"])
			}
			if asset.Resource.Data["state"] != tt.expectedState {
				t.Errorf("Case: %s. Error converting asset state. Expected: %q got: %q", tt.name, tt.expectedService, asset.Resource.Data["state"])
			}
		})
	}
}

type mockTerraformResourceData struct {
	m map[string]interface{}
	tpgresource.TerraformResourceData
}

func (d *mockTerraformResourceData) GetOkExists(k string) (interface{}, bool) {
	v, ok := d.m[k]
	return v, ok
}

func (d *mockTerraformResourceData) GetOk(k string) (interface{}, bool) {
	v, ok := d.m[k]
	return v, ok
}

func (d *mockTerraformResourceData) Get(k string) interface{} {
	v, ok := d.m[k]
	if !ok {
		return nil
	}
	return v
}
