package networkconnectivity

import (
	"testing"
)

func TestNetworkConnectivityParseHub(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		hub     string
		project string
		wantP   string
		wantN   string
	}{
		"full uri": {
			hub:     "projects/hub-proj/locations/global/hubs/my-hub",
			project: "spoke-proj",
			wantP:   "hub-proj",
			wantN:   "my-hub",
		},
		"name only": {
			hub:     "my-hub",
			project: "spoke-proj",
			wantP:   "spoke-proj",
			wantN:   "my-hub",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			gotP, gotN := networkConnectivityParseHub(tc.hub, tc.project)
			if gotP != tc.wantP || gotN != tc.wantN {
				t.Errorf("networkConnectivityParseHub(%q, %q) = (%q, %q), want (%q, %q)", tc.hub, tc.project, gotP, gotN, tc.wantP, tc.wantN)
			}
		})
	}
}

func TestNetworkConnectivitySpokeUpdatePending(t *testing.T) {
	t.Parallel()

	if networkConnectivitySpokeUpdatePending(nil) {
		t.Fatal("nil spoke should not be pending")
	}
	if networkConnectivitySpokeUpdatePending(map[string]interface{}{"state": "ACTIVE"}) {
		t.Fatal("active spoke should not be pending")
	}
	if !networkConnectivitySpokeUpdatePending(map[string]interface{}{
		"reasons": []interface{}{map[string]interface{}{"code": "UPDATE_PENDING_REVIEW"}},
	}) {
		t.Fatal("UPDATE_PENDING_REVIEW should be pending")
	}
}
