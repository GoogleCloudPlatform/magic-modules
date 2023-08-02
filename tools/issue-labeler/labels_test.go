package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

func TestLabels(t *testing.T) {
	file, err := os.ReadFile("enrolled_teams.yml")
	if err != nil {
		glog.Exitf("Error reading enrolled teams yaml: %v", err)
	}
	enrolledTeams := make(map[string][]string)
	err = yaml.Unmarshal(file, &enrolledTeams)
	if err != nil {
		glog.Exitf("Error unmarshalling enrolled teams yaml: %v", err)
	}
	for _, tc := range []struct {
		issueBody      string
		expectedLabels []string
	}{
		{
			issueBody: `### New or Affected Resource(s):
			google_gke_hub_feature
			google_storage_hmac_key
			#`,
			expectedLabels: []string{"service/gkehub", "service/storage"},
		},
		{
			issueBody: `### New or Affected Resource(s):
			#`,
		},
	} {
		if actualLabels := serviceLabels(tc.issueBody, enrolledTeams); !reflect.DeepEqual(actualLabels, tc.expectedLabels) {
			t.Errorf("unexpected labels for issue body %s: %v, expected %v", tc.issueBody, actualLabels, tc.expectedLabels)
		}
	}
}
