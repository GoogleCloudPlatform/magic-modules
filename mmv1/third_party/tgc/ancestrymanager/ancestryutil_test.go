package ancestrymanager

import (
	"testing"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
)

func TestAssetParent(t *testing.T) {
	tests := []struct {
		name      string
		ancestors []string
		tfData    tpgresource.TerraformResourceData
		cai       *resources.Asset
		want      string
		wantErr   bool
	}{
		{
			name:      "project",
			ancestors: []string{"projects/123", "folders/456", "organizations/789"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Project"},
			want:      "//cloudresourcemanager.googleapis.com/folders/456",
		},
		{
			// new project without org_id and folder_id cannot derive ancestor list
			name:      "project with no other ancestors",
			ancestors: []string{"projects/123"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Project"},
			want:      "//cloudresourcemanager.googleapis.com/projects/123",
		},
		{
			name:      "new project with only unknown org",
			ancestors: []string{"organizations/unknown"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Project"},
			want:      "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
		{
			name:      "folder",
			ancestors: []string{"folders/456", "organizations/789"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Folder"},
			want:      "//cloudresourcemanager.googleapis.com/organizations/789",
		},
		{
			name:      "new folder with unknown org",
			ancestors: []string{"organizations/unknown"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Folder"},
			want:      "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
		{
			name:      "organization",
			ancestors: []string{"organizations/789"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Organization"},
			want:      "",
		},
		{
			name:      "other resource",
			ancestors: []string{"projects/123", "folders/456", "organizations/789"},
			cai:       &resources.Asset{Type: "storage.googleapis.com/Bucket"},
			want:      "//cloudresourcemanager.googleapis.com/projects/123",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := assetParent(test.cai, test.ancestors)
			if err != nil {
				t.Fatalf("AssetParent(%v, %v) = %v, want = nil", test.cai, test.ancestors, err)
			}
			if got != test.want {
				t.Errorf("AssetParent(%v, %v) = %v, want = %v", test.cai, test.ancestors, got, test.want)
			}
		})
	}
}

func TestAssetParent_Fail(t *testing.T) {
	tests := []struct {
		name      string
		ancestors []string
		tfData    tpgresource.TerraformResourceData
		cai       *resources.Asset
	}{
		{
			name:      "project",
			ancestors: []string{},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Project"},
		},
		{
			name:      "folder",
			ancestors: []string{"folders/456"},
			cai:       &resources.Asset{Type: "cloudresourcemanager.googleapis.com/Folder"},
		},
		{
			name:      "other resource",
			ancestors: []string{},
			cai:       &resources.Asset{Type: "storage.googleapis.com/Bucket"},
		},
		{
			name:      "invalid ancestor",
			ancestors: []string{"abc/3"},
		},
		{
			name:      "empty project ID",
			ancestors: []string{"projects/"},
		},
		{
			name:      "empty folder ID",
			ancestors: []string{"folders/"},
		},
		{
			name:      "empty organization ID",
			ancestors: []string{"organizations/"},
		},
		{
			name:      "asset not provided",
			ancestors: []string{"organizations/123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := assetParent(test.cai, test.ancestors)
			if err == nil {
				t.Fatalf("AssetParent(%v, %v) = nil, want = error", test.cai, test.ancestors)
			}
		})
	}
}

func TestConvertToAncestryPath(t *testing.T) {
	cases := []struct {
		name           string
		input          []string
		expectedOutput string
	}{
		{
			name:           "Empty",
			input:          []string{},
			expectedOutput: "",
		},
		{
			name:           "ProjectOrganization",
			input:          []string{"project/my-prj", "organization/my-org"},
			expectedOutput: "organization/my-org/project/my-prj",
		},
		{
			name:           "convert to existing ancestry path style",
			input:          []string{"projects/my-prj", "folders/123", "organizations/my-org"},
			expectedOutput: "organization/my-org/folder/123/project/my-prj",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := ConvertToAncestryPath(c.input)
			if output != c.expectedOutput {
				t.Errorf("expected output %q, got %q", c.expectedOutput, output)
			}
		})
	}
}
