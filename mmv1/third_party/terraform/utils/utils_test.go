package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestGetProject(t *testing.T) {
	cases := map[string]struct {
		ResourceProject string
		ProviderProject string
		ExpectedProject string
		ExpectedError   bool
	}{
		"project is pulled from resource config instead of provider config": {
			ResourceProject: "foo",
			ProviderProject: "bar",
			ExpectedProject: "foo",
		},
		"project is pulled from provider config when not set on resource": {
			ProviderProject: "bar",
			ExpectedProject: "bar",
		},
		"error returned when project not set on either provider or resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if tc.ProviderProject != "" {
				config.Project = tc.ProviderProject
			}

			// Create resource config
			// Here use ResourceComputeDisk schema as example
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeDisk().Schema, emptyConfigMap)
			if tc.ResourceProject != "" {
				if err := d.Set("project", tc.ResourceProject); err != nil {
					t.Fatalf("Cannot set project: %s", err)
				}
			}

			// Act
			project, err := tpgresource.GetProject(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if project != tc.ExpectedProject {
				t.Fatalf("Incorrect project: got %s, want %s", project, tc.ExpectedProject)
			}
		})
	}
}

func TestGetZone(t *testing.T) {
	cases := map[string]struct {
		ResourceZone  string
		ProviderZone  string
		ExpectedZone  string
		ExpectedError bool
	}{
		"zone is pulled from resource config instead of provider config": {
			ResourceZone: "foo",
			ProviderZone: "bar",
			ExpectedZone: "foo",
		},
		"zone value from resource can be a self link": {
			ResourceZone: "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a",
			ExpectedZone: "us-central1-a",
		},
		"zone is pulled from provider config when not set on resource": {
			ProviderZone: "bar",
			ExpectedZone: "bar",
		},
		"error returned when zone not set on either provider or resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if tc.ProviderZone != "" {
				config.Zone = tc.ProviderZone
			}

			// Create resource config
			// Here use ResourceComputeDisk schema as example - because it has a zone field in schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeDisk().Schema, emptyConfigMap)
			if tc.ResourceZone != "" {
				if err := d.Set("zone", tc.ResourceZone); err != nil {
					t.Fatalf("Cannot set zone: %s", err)
				}
			}

			// Act
			zone, err := tpgresource.GetZone(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if zone != tc.ExpectedZone {
				t.Fatalf("Incorrect zone: got %s, want %s", zone, tc.ExpectedZone)
			}
		})
	}
}

func TestGetRegion(t *testing.T) {
	cases := map[string]struct {
		ResourceRegion string
		ProviderRegion string
		ProviderZone   string
		ExpectedRegion string
		ExpectedZone   string
		ExpectedError  bool
	}{
		"region is pulled from resource config instead of provider config": {
			ResourceRegion: "foo",
			ProviderRegion: "bar",
			ProviderZone:   "lol-a",
			ExpectedRegion: "foo",
		},
		"region pulled from resource config can be a self link": {
			ResourceRegion: "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1",
			ExpectedRegion: "us-central1",
		},
		"region is pulled from region on provider config when region unset in resource config": {
			ProviderRegion: "bar",
			ProviderZone:   "lol-a",
			ExpectedRegion: "bar",
		},
		"region is pulled from zone on provider config when region unset in both resource and provider config": {
			ProviderZone:   "lol-a",
			ExpectedRegion: "lol",
		},
		"error returned when region not set on resource and neither region or zone set on provider": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if tc.ProviderRegion != "" {
				config.Region = tc.ProviderRegion
			}
			if tc.ProviderZone != "" {
				config.Zone = tc.ProviderZone
			}

			// Create resource config
			// Here use ResourceComputeSubnetwork schema as example - because it has a region field in schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeSubnetwork().Schema, emptyConfigMap)
			if tc.ResourceRegion != "" {
				if err := d.Set("region", tc.ResourceRegion); err != nil {
					t.Fatalf("Cannot set region: %s", err)
				}
			}

			// Act
			region, err := tpgresource.GetRegion(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}

func TestCheckGCSName(t *testing.T) {
	valid63 := RandString(t, 63)
	cases := map[string]bool{
		// Valid
		"foobar":       true,
		"foobar1":      true,
		"12345":        true,
		"foo_bar_baz":  true,
		"foo-bar-baz":  true,
		"foo-bar_baz1": true,
		"foo--bar":     true,
		"foo__bar":     true,
		"foo-goog":     true,
		"foo.goog":     true,
		valid63:        true,
		fmt.Sprintf("%s.%s.%s", valid63, valid63, valid63): true,

		// Invalid
		"goog-foobar":     false,
		"foobar-google":   false,
		"-foobar":         false,
		"foobar-":         false,
		"_foobar":         false,
		"foobar_":         false,
		"fo":              false,
		"foo$bar":         false,
		"foo..bar":        false,
		RandString(t, 64): false,
		fmt.Sprintf("%s.%s.%s.%s", valid63, valid63, valid63, valid63): false,
	}

	for bucketName, valid := range cases {
		err := tpgresource.CheckGCSName(bucketName)
		if valid && err != nil {
			t.Errorf("The bucket name %s was expected to pass validation and did not pass.", bucketName)
		} else if !valid && err == nil {
			t.Errorf("The bucket name %s was NOT expected to pass validation and passed.", bucketName)
		}
	}
}
