package tgcresource

import (
	"testing"
)

func TestParseFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		fieldName string
		expected  string
	}{
		{
			name:      "Success_StandardURL",
			url:       "projects/my-project-id/locations/us-central1/repositories/my-repo",
			fieldName: "repositories",
			expected:  "my-repo",
		},
		{
			name:      "Success_MultipleFragments",
			url:       "projects/123/zones/us-west1-a/instances/test-vm/disks/boot-disk",
			fieldName: "instances",
			expected:  "test-vm",
		},
		{
			name:      "Success_FirstFragment",
			url:       "projects/123/zones/us-west1-a/instances/test-vm",
			fieldName: "projects",
			expected:  "123",
		},
		{
			name:      "Success_LastFragment",
			url:       "projects/123/zones/us-west1-a/instances/test-vm",
			fieldName: "instances",
			expected:  "test-vm",
		},
		{
			name:      "Success_LeadingAndTrailingSlashes",
			url:       "/regions/us-east1/addresses/my-ip/",
			fieldName: "addresses",
			expected:  "my-ip",
		},
		{
			name:      "Fail_FieldNotFound",
			url:       "projects/123/regions/us-east1/networks/default",
			fieldName: "subnetworks",
			expected:  "",
		},
		{
			name:      "Fail_EmptyURL",
			url:       "",
			fieldName: "projects",
			expected:  "",
		},
		{
			name:      "Fail_EmptyField",
			url:       "projects/123/regions/us-east1",
			fieldName: "",
			expected:  "",
		},
		{
			name:      "Fail_ValueMissing",
			url:       "projects/123/regions/",
			fieldName: "regions",
			expected:  "",
		},
		{
			name:      "Fail_OnlySlashes",
			url:       "///",
			fieldName: "projects",
			expected:  "",
		},
		{
			name:      "Fail_FieldNameInValue", // Ensure it only matches the field key
			url:       "projects/123/regions/projects",
			fieldName: "regions",
			expected:  "projects", // The function correctly returns 'projects' here, which is expected behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseFieldValue(tt.url, tt.fieldName)

			if actual != tt.expected {
				t.Errorf("ParseFieldValue(%q, %q) returned %q, want %q",
					tt.url, tt.fieldName, actual, tt.expected)
			}
		})
	}
}
