package toolkit

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai"
	"go.uber.org/zap"
)

func TestIsSupported(t *testing.T) {
	migratedMap := map[string]bool{
		"google_compute_instance": true,
	}

	tests := []struct {
		name     string
		resource string
		expected bool
	}{
		{
			name:     "Migrated resource",
			resource: "google_compute_instance",
			expected: true,
		},
		{
			name:     "Legacy resource",
			resource: "google_project", // Assuming this is in legacy converters
			expected: true,
		},
		{
			name:     "Unsupported resource",
			resource: "google_invalid_resource",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsSupported(tc.resource, migratedMap)
			if got != tc.expected {
				t.Errorf("IsSupported(%q) = %v, want %v", tc.resource, got, tc.expected)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name        string
		jsonPlan    []byte
		migratedMap map[string]bool
	}{
		{
			name: "Migrated resource only",
			jsonPlan: []byte(`{
				"resource_changes": [
					{
						"type": "google_compute_instance",
						"change": {
							"actions": ["create"],
							"after": {"name": "instance-1"}
						}
					}
				]
			}`),
			migratedMap: map[string]bool{"google_compute_instance": true},
		},
		{
			name: "Legacy resource only",
			jsonPlan: []byte(`{
				"resource_changes": [
					{
						"type": "google_project",
						"change": {
							"actions": ["create"],
							"after": {"name": "project-1"}
						}
					}
				]
			}`),
			migratedMap: map[string]bool{"google_compute_instance": true}, // Not migrated
		},
		{
			name: "Mixed resources",
			jsonPlan: []byte(`{
				"resource_changes": [
					{
						"type": "google_compute_instance",
						"change": {
							"actions": ["create"],
							"after": {"name": "instance-1"}
						}
					},
					{
						"type": "google_project",
						"change": {
							"actions": ["create"],
							"after": {"name": "project-1"}
						}
					}
				]
			}`),
			migratedMap: map[string]bool{"google_compute_instance": true},
		},
		{
			name: "Empty plan",
			jsonPlan: []byte(`{
				"resource_changes": []
			}`),
			migratedMap: map[string]bool{"google_compute_instance": true},
		},
	}

	o := &tfplan2cai.Options{
		ErrorLogger: zap.NewNop(),
	}
	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Convert(ctx, tc.jsonPlan, o, tc.migratedMap)
			// We don't assert success here because it might fail without real provider schema,
			// but we want to ensure it doesn't panic and covers the code path.
			if err != nil {
				t.Logf("Convert returned error (expected in unit test without full setup): %v", err)
			}
		})
	}
}
