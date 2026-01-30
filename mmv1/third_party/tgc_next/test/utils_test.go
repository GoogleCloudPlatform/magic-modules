package test

import (
	"testing"
)

func TestGetSubTestName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple nested test",
			input:    "TestAccMonitoringAlertPolicy/basic",
			expected: "basic",
		},
		{
			name:     "Deeply nested test",
			input:    "TestAccMonitoringAlertPolicy/TestAccMonitoringAlertPolicy/basic",
			expected: "TestAccMonitoringAlertPolicy/basic",
		},
		{
			name:     "No slash",
			input:    "TestAccMonitoringAlertPolicy",
			expected: "", // Or should it be empty? User said "after the first /". If no slash, probably empty string or original string? Let's assume empty if no slash, or we can treat as "no subtest".
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Slash at end",
			input:    "TestAccMonitoringAlertPolicy/",
			expected: "",
		},
		{
			name:     "Slash at start",
			input:    "/basic",
			expected: "basic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSubTestName(tt.input)
			if got != tt.expected {
				t.Errorf("GetSubTestName(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
