package networkconnectivity

import (
	"testing"
)

func TestNetworkConnectivitySpokeIsAccepted(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"ACTIVE":          true,
		"UPDATING":        true,
		"UPDATE_FAILED":   true,
		"UPDATE_REJECTED": true,
		"INACTIVE":        false,
		"CREATING":        false,
		"OBSOLETE":        false,
		"DELETING":        false,
		"":                false,
	}

	for state, want := range cases {
		if got := NetworkConnectivitySpokeIsAccepted(state); got != want {
			t.Errorf("NetworkConnectivitySpokeIsAccepted(%q) = %v, want %v", state, got, want)
		}
	}
}

func TestNetworkConnectivitySpokeHasPendingUpdate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		spoke map[string]interface{}
		want  bool
	}{
		"nil": {},
		"active with no reasons": {
			spoke: map[string]interface{}{"state": "ACTIVE"},
		},
		"update pending review": {
			spoke: map[string]interface{}{
				"state": "ACTIVE",
				"reasons": []interface{}{
					map[string]interface{}{"code": "UPDATE_PENDING_REVIEW"},
				},
			},
			want: true,
		},
		"field paths pending update": {
			spoke: map[string]interface{}{
				"fieldPathsPendingUpdate": []interface{}{"linkedVpcNetwork.includeExportRanges"},
			},
			want: true,
		},
		"other reason": {
			spoke: map[string]interface{}{
				"reasons": []interface{}{
					map[string]interface{}{"code": "REJECTED"},
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if got := NetworkConnectivitySpokeHasPendingUpdate(tc.spoke); got != tc.want {
				t.Errorf("NetworkConnectivitySpokeHasPendingUpdate(%v) = %v, want %v", tc.spoke, got, tc.want)
			}
		})
	}
}

func TestNetworkConnectivityHubSpokeAcceptanceId(t *testing.T) {
	t.Parallel()

	got := networkConnectivityHubSpokeAcceptanceId(
		"projects/hub-proj",
		"projects/hub-proj/locations/global/hubs/my-hub",
		"projects/spoke-proj/locations/global/spokes/my-spoke",
	)
	want := "projects/hub-proj/locations/global/hubs/my-hub/spokeAcceptances/projects/spoke-proj/locations/global/spokes/my-spoke"
	if got != want {
		t.Errorf("networkConnectivityHubSpokeAcceptanceId() = %q, want %q", got, want)
	}
}
