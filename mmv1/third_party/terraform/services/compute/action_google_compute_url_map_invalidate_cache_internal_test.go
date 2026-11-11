package compute

import (
	"testing"
)

func TestParseComputeUrlMapRef(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		urlMap, project string
		wantProject     string
		wantName        string
		wantErr         bool
	}{
		"name only": {
			urlMap:   "my-map",
			wantName: "my-map",
		},
		"name with project attribute": {
			urlMap:      "my-map",
			project:     "p",
			wantProject: "p",
			wantName:    "my-map",
		},
		"global partial URI": {
			urlMap:   "global/urlMap/my-map",
			wantName: "my-map",
		},
		"project partial URI": {
			urlMap:      "projects/p/global/urlMap/my-map",
			wantProject: "p",
			wantName:    "my-map",
		},
		"self link with matching project": {
			urlMap:      "https://www.googleapis.com/compute/v1/projects/p/global/urlMap/my-map",
			project:     "p",
			wantProject: "p",
			wantName:    "my-map",
		},
		"conflicting project": {
			urlMap:  "projects/p/global/urlMap/my-map",
			project: "other",
			wantErr: true,
		},
		"regional URL map": {
			urlMap:  "projects/p/regions/us-central1/urlMap/my-map",
			wantErr: true,
		},
		"wrong resource type": {
			urlMap:  "projects/p/global/backendServices/my-service",
			wantErr: true,
		},
		"trailing path segments": {
			urlMap:  "projects/p/global/urlMap/my-map/extra",
			wantErr: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseComputeUrlMapRef(tc.urlMap, tc.project)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Project != tc.wantProject || got.Name != tc.wantName {
				t.Errorf("expected project %q and name %q, got project %q and name %q", tc.wantProject, tc.wantName, got.Project, got.Name)
			}
		})
	}
}

func TestValidateComputeUrlMapInvalidationPath(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path    string
		wantErr bool
	}{
		"exact path":           {path: "/index.html"},
		"prefix wildcard":      {path: "/images/*"},
		"everything":           {path: "/*"},
		"root":                 {path: "/"},
		"missing leading /":    {path: "images/*", wantErr: true},
		"empty":                {path: "", wantErr: true},
		"query string":         {path: "/index.html?v=2", wantErr: true},
		"fragment":             {path: "/index.html#top", wantErr: true},
		"wildcard not at end":  {path: "/images/*.png", wantErr: true},
		"two wildcards at end": {path: "/images/**", wantErr: true},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := ValidateComputeUrlMapInvalidationPath(tc.path)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q", tc.path)
				return
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %s", tc.path, err)
			}
		})
	}
}
