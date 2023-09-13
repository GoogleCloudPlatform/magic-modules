package bigtable

import (
	"testing"
)

func TestGetUnavailableClusterZones(t *testing.T) {
	cases := map[string]struct {
		clusterZones     []string
		unavailableZones []string
		want             []string
	}{
		"not unavailalbe": {
			clusterZones:     []string{"us-central1", "eu-west1"},
			unavailableZones: []string{"us-central2", "eu-west2"},
			want:             nil,
		},
		"unavailable one to one": {
			clusterZones:     []string{"us-central2"},
			unavailableZones: []string{"us-central2"},
			want:             []string{"us-central2"},
		},
		"unavailable one to many": {
			clusterZones:     []string{"us-central2"},
			unavailableZones: []string{"us-central2", "us-central1"},
			want:             []string{"us-central2"},
		},
		"unavailable many to one": {
			clusterZones:     []string{"us-central2", "us-central1"},
			unavailableZones: []string{"us-central2"},
			want:             []string{"us-central2"},
		},
		"unavailable many to many": {
			clusterZones:     []string{"us-central2", "us-central1"},
			unavailableZones: []string{"us-central2", "us-central1"},
			want:             []string{"us-central2", "us-central1"},
		},
	}

	for tn, tc := range cases {
		clusters := []map[string]string{}
		for _, zone := range tc.clusterZones {
			clusters.append(map[string]string{"zone": tc.clusterZone})
		}
		if got := getUnavailableClusterZones(clusterstc.unavailableZones); got != tc.want {
			t.Errorf("bad: %s, got %q, want %q", tn, got, tc.want)
		}
	}
}
