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
	mapHCL = `
resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "datasetjlgukul2im"
  description   = "This is a test description"
  friendly_name = "test"
  location      = "EU"
  project       = "ci-test-project-nightly-beta"
  resource_tags = {
    "ci-test-project-nightly-beta/tf_test_tag_key1jlgukul2im" = "tf_test_tag_value1jlgukul2im"
    "ci-test-project-nightly-beta/tf_test_tag_key2jlgukul2im" = "tf_test_tag_value2jlgukul2im"
  }
}
`

	boolFieldWithFalseHCL = `
resource "google_compute_backend_bucket" "image_backend" {
  name        = "tf-test-image-backend-bucket"
  description = "Contains beautiful images"
  bucket_name = "tf-test-image-backend-bucket"
  enable_cdn  = false
  cdn_policy {
    cache_key_policy {
        query_string_whitelist = ["image-version"]
    }
  }
}
`

	boolFieldUnsetHCL = `
resource "google_compute_backend_bucket" "image_backend" {
  name        = "tf-test-image-backend-bucket"
  description = "Contains beautiful images"
  bucket_name = "tf-test-image-backend-bucket"
  cdn_policy {
    cache_key_policy {
        query_string_whitelist = ["image-version"]
    }
  }
}
`
)

func TestParseHCLBytes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		hcl       string
		exp       map[string]map[string]any
		expectErr bool
	}{
		{
			name: "basic",
			hcl:  basicHCL,
			exp: map[string]map[string]any{
				"google_project_service.project": {
					"service": "iam.googleapis.com",
				},
			},
		},
		{
			name: "nested blocks",
			hcl:  nestedBlocksHCL,
			exp: map[string]map[string]any{
				"google_storage_bucket.bucket": {
					"name":                         "my-bucket",
					"location":                     "US",
					"force_destroy":                true,
					"lifecycle_rule.action.type":   "Delete",
					"lifecycle_rule.condition.age": float64(30),
				},
			},
		},
		{
			name: "multiple resources",
			hcl:  multipleResourcesHCL,
			exp: map[string]map[string]any{
				"google_project_service.project": {
					"service": "iam.googleapis.com",
				},
				"google_storage_bucket.bucket": {
					"name": "my-bucket",
				},
			},
		},
		{
			name: "resource with a list of nested objects",
			hcl:  listOfNestedObjectsHCL,
			exp: map[string]map[string]any{
				"google_compute_firewall.default": {
					"allow.0.ports":    "80,8080,1000-2000", // "ports" appears in first element due to sorting
					"allow.0.protocol": "tcp",
					"allow.1.protocol": "icmp",
					"name":             "test-firewall",
					"network":          "google_compute_network.default.name",
					"source_tags":      "web",
				},
			},
		},
		{
			name: "resource with a list of multi-level nested objects",
			hcl:  listOfMultiLevelNestedObjectsHCL,
			exp: map[string]map[string]any{
				"google_compute_firewall.default": {
					"allow.0.a_second_level.0.a": false,
					"allow.0.a_second_level.1.b": true,
					"allow.0.protocol":           "icmp",
					"allow.1.ports":              "80,8080,1000-2000",
					"allow.1.protocol":           "tcp",
					"name":                       "test-firewall",
					"network":                    "google_compute_network.default.name",
					"source_tags":                "web",
				},
			},
		},
		{
			name: "resource with map",
			hcl:  mapHCL,
			exp: map[string]map[string]any{
				"google_bigquery_dataset.dataset": {
					"dataset_id":    "datasetjlgukul2im",
					"description":   "This is a test description",
					"friendly_name": "test",
					"location":      "EU",
					"project":       "ci-test-project-nightly-beta",
					"resource_tags": map[string]any{},
				},
			},
		},
		{
			name: "resource with false bool field",
			hcl:  boolFieldWithFalseHCL,
			exp: map[string]map[string]any{
				"google_compute_backend_bucket.image_backend": {
					"name":        "tf-test-image-backend-bucket",
					"description": "Contains beautiful images",
					"bucket_name": "tf-test-image-backend-bucket",
					"enable_cdn":  false,
					"cdn_policy.cache_key_policy.query_string_whitelist": "image-version",
				},
			},
		},
		{
			name: "resource with unset bool field",
			hcl:  boolFieldUnsetHCL,
			exp: map[string]map[string]any{
				"google_compute_backend_bucket.image_backend": {
					"name":        "tf-test-image-backend-bucket",
					"description": "Contains beautiful images",
					"bucket_name": "tf-test-image-backend-bucket",
					"cdn_policy.cache_key_policy.query_string_whitelist": "image-version",
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
