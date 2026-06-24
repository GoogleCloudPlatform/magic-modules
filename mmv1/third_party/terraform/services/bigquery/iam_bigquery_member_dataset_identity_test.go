package bigquery

import "testing"

func TestBigqueryDatasetIamMemberResource_HasIdentity(t *testing.T) {
	resource := BigqueryDatasetIamMemberResource()
	if resource.Identity == nil {
		t.Fatalf("expected google_bigquery_dataset_iam_member resource identity to be configured")
	}
}

func TestBigqueryDatasetIamParentResourceIdentityParser(t *testing.T) {
	tests := []struct {
		name      string
		project   string
		datasetID string
		want      string
	}{
		{
			name:      "split fields",
			project:   "my-project",
			datasetID: "my_dataset",
			want:      "projects/my-project/datasets/my_dataset",
		},
		{
			name:      "dataset id as full resource id",
			datasetID: "projects/my-project/datasets/my_dataset",
			want:      "projects/my-project/datasets/my_dataset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := BigqueryDatasetIamMemberResource()
			rd := resource.TestResourceData()
			identity, err := rd.Identity()
			if err != nil {
				t.Fatalf("creating identity: %v", err)
			}
			if tt.project != "" {
				if err := identity.Set("project", tt.project); err != nil {
					t.Fatalf("setting identity project: %v", err)
				}
			}
			if err := identity.Set("dataset_id", tt.datasetID); err != nil {
				t.Fatalf("setting identity dataset_id: %v", err)
			}

			got, err := BigqueryDatasetIamParentResourceIdentityParser(rd, identity, nil)
			if err != nil {
				t.Fatalf("parsing identity: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected parsed id: got %q, want %q", got, tt.want)
			}
		})
	}
}
