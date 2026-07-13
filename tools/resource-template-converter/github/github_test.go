package github

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mmv1/products/accesscontextmanager/AccessLevel.yaml", "products/accesscontextmanager/AccessLevel.yaml"},
		{"products/accesscontextmanager/AccessLevel.yaml", "products/accesscontextmanager/AccessLevel.yaml"},
		{"/abs/path/mmv1/products/accesscontextmanager/AccessLevel.yaml", "/abs/path/mmv1/products/accesscontextmanager/AccessLevel.yaml"},
	}

	for _, tc := range tests {
		got := NormalizePath(tc.input)
		if got != tc.expected {
			t.Errorf("NormalizePath(%q) = %q, expected %q", tc.input, got, tc.expected)
		}
	}
}
