package compute

import (
	"testing"
)

func TestComputeBackendServiceBackendHash_preference(t *testing.T) {
	t.Parallel()

	backend := func(preference interface{}) map[string]interface{} {
		b := map[string]interface{}{
			"group":           "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-f/instanceGroups/my-group",
			"balancing_mode":  "UTILIZATION",
			"capacity_scaler": 1.0,
			"description":     "",
			"max_utilization": 0.8,
		}
		if preference != nil {
			b["preference"] = preference
		}
		return b
	}

	cases := map[string]struct {
		a, b      map[string]interface{}
		wantEqual bool
	}{
		"PREFERRED and DEFAULT hash differently": {
			a:         backend("PREFERRED"),
			b:         backend("DEFAULT"),
			wantEqual: false,
		},
		"PREFERRED and unset hash differently": {
			a:         backend("PREFERRED"),
			b:         backend(""),
			wantEqual: false,
		},
		"DEFAULT and empty string hash the same": {
			a:         backend("DEFAULT"),
			b:         backend(""),
			wantEqual: true,
		},
		"DEFAULT and missing key hash the same": {
			a:         backend("DEFAULT"),
			b:         backend(nil),
			wantEqual: true,
		},
		"same preference hashes the same": {
			a:         backend("PREFERRED"),
			b:         backend("PREFERRED"),
			wantEqual: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ha := resourceGoogleComputeBackendServiceBackendHash(tc.a)
			hb := resourceGoogleComputeBackendServiceBackendHash(tc.b)
			if (ha == hb) != tc.wantEqual {
				t.Errorf("hash(%v) = %d, hash(%v) = %d, want equal: %v", tc.a["preference"], ha, tc.b["preference"], hb, tc.wantEqual)
			}
		})
	}
}
