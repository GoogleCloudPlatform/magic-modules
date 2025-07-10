package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	basicHCL = `
resource "google_project_service" "project" {
  service = "iam.googleapis.com"
}
`
	nestedBlocksHCL = `
resource "google_storage_bucket" "bucket" {
  name          = "my-bucket"
  location      = "US"
  force_destroy = true

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = 30
    }
  }
}
`
	multipleResourcesHCL = `
resource "google_project_service" "project" {
  service = "iam.googleapis.com"
}

resource "google_storage_bucket" "bucket" {
  name = "my-bucket"
}
`
	listOfNestedObjectsHCL = `
resource "google_compute_firewall" "default" {
  name    = "test-firewall"
  network = google_compute_network.default.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }

  source_tags = ["web"]
}
`
	listOfMultiLevelNestedObjectsHCL = `
resource "google_compute_firewall" "default" {
  name    = "test-firewall"
  network = google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }

  allow {
    protocol = "icmp"
    a_second_level {
      b = true
    }
    a_second_level {
      a = false
    }
  }

  source_tags = ["web"]
}
`
)

func TestParseHCLBytes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		hcl       string
		exp       map[string]map[string]struct{}
		expectErr bool
	}{
		{
			name: "basic",
			hcl:  basicHCL,
			exp: map[string]map[string]struct{}{
				"google_project_service.project": {
					"service": {},
				},
			},
		},
		{
			name: "nested blocks",
			hcl:  nestedBlocksHCL,
			exp: map[string]map[string]struct{}{
				"google_storage_bucket.bucket": {
					"name":                         {},
					"location":                     {},
					"force_destroy":                {},
					"lifecycle_rule.action.type":   {},
					"lifecycle_rule.condition.age": {},
				},
			},
		},
		{
			name: "multiple resources",
			hcl:  multipleResourcesHCL,
			exp: map[string]map[string]struct{}{
				"google_project_service.project": {
					"service": {},
				},
				"google_storage_bucket.bucket": {
					"name": {},
				},
			},
		},
		{
			name: "resource with a list of nested objects",
			hcl:  listOfNestedObjectsHCL,
			exp: map[string]map[string]struct{}{
				"google_compute_firewall.default": {
					"allow.0.ports":    {}, // "ports" appears in first element due to sorting
					"allow.0.protocol": {},
					"allow.1.protocol": {},
					"name":             {},
					"network":          {},
					"source_tags":      {},
				},
			},
		},
		{
			name: "resource with a list of multi-level nested objects",
			hcl:  listOfMultiLevelNestedObjectsHCL,
			exp: map[string]map[string]struct{}{
				"google_compute_firewall.default": {
					"allow.0.a_second_level.0.a": {},
					"allow.0.a_second_level.1.b": {},
					"allow.0.protocol":           {},
					"allow.1.ports":              {},
					"allow.1.protocol":           {},
					"name":                       {},
					"network":                    {},
					"source_tags":                {},
				},
			},
		},
		{
			name:      "invalid hcl",
			hcl:       `resource "google_project_service" "project" {`,
			expectErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseHCLBytes([]byte(tc.hcl), "test.hcl")
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.exp, got); diff != "" {
				t.Errorf("unexpected diff (-want +got): %s", diff)
			}
		})
	}
}
