package transport_test

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestBasePathDefault(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials: transport_tpg.TestFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
	}
	cases := map[string]struct {
		BasePath       string
		RepPath        string
		BasePathKey    string
		Config         *transport_tpg.Config
		Location       string
		ExpectedOutput string
	}{
		"Default to global path": {
			BasePath:       "https://clouddeploy.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://clouddeploy.googleapis.com/v1/",
		},
		"Overridden path takes priority": {
			BasePath:       "https://override.{{location}}.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://override.us-central1.googleapis.com/v1/",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			basePath, _ := transport_tpg.BasePath(tc.BasePath, tc.RepPath, tc.BasePathKey, tc.Config, tc.Location)

			if basePath != tc.ExpectedOutput {
				t.Fatalf("want %s,  got %s", tc.ExpectedOutput, basePath)
			}
		})
	}
}

func TestBasePathPreferGlobal(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials:           transport_tpg.TestFakeCredentialsPath,
		Project:               "my-gce-project",
		Region:                "us-central1",
		PreferGlobalEndpoints: true,
	}
	cases := map[string]struct {
		BasePath       string
		RepPath        string
		BasePathKey    string
		Config         *transport_tpg.Config
		Location       string
		ExpectedOutput string
	}{
		"Default to global path": {
			BasePath:       "https://clouddeploy.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://clouddeploy.googleapis.com/v1/",
		},
		"Overridden path takes priority": {
			BasePath:       "https://override.{{location}}.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://override.us-central1.googleapis.com/v1/",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			basePath, _ := transport_tpg.BasePath(tc.BasePath, tc.RepPath, tc.BasePathKey, tc.Config, tc.Location)

			if basePath != tc.ExpectedOutput {
				t.Fatalf("want %s,  got %s", tc.ExpectedOutput, basePath)
			}
		})
	}
}

func TestBasePathPreferRegional(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials:             transport_tpg.TestFakeCredentialsPath,
		Project:                 "my-gce-project",
		Region:                  "us-central1",
		PreferRegionalEndpoints: true,
	}
	cases := map[string]struct {
		BasePath       string
		RepPath        string
		BasePathKey    string
		Config         *transport_tpg.Config
		Location       string
		ExpectedOutput string
	}{
		"Default to regional path": {
			BasePath:       "https://clouddeploy.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://www.clouddeploy.us-central1.rep.googleapis.com/v1/",
		},
		"Overridden path takes priority": {
			BasePath:       "https://override.{{location}}.googleapis.com/v1/",
			RepPath:        "https://www.clouddeploy.{{location}}.rep.googleapis.com/v1/",
			BasePathKey:    "Clouddeploy",
			Config:         config,
			Location:       "us-central1",
			ExpectedOutput: "https://override.us-central1.googleapis.com/v1/",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			basePath, _ := transport_tpg.BasePath(tc.BasePath, tc.RepPath, tc.BasePathKey, tc.Config, tc.Location)

			if basePath != tc.ExpectedOutput {
				t.Fatalf("want %s,  got %s", tc.ExpectedOutput, basePath)
			}
		})
	}
}
