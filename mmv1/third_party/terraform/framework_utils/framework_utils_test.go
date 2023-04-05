package google

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCompileUserAgentString(t *testing.T) {
	cases := map[string]struct {
		Name              string // Name of the provider
		TerraformVersion  string
		ProviderVersion   string
		EnvValue          string
		ExpectedUserAgent string
	}{
		"the expected user agent is returned for given inputs": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9",
		},
		"the user agent can have values appended via an environment variable": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			EnvValue:          "I'm appended at the end!",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9 I'm appended at the end!",
		},
		"values appended via an environment variable have whitespace trimmed": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			EnvValue:          "              my surrounding white space is removed              ",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9 my surrounding white space is removed",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			t.Setenv(uaEnvVar, tc.EnvValue) // Use same global const as the CompileUserAgentString function

			// Act

			ua := CompileUserAgentString(ctx, tc.Name, tc.TerraformVersion, tc.ProviderVersion)

			// Assert
			if ua != tc.ExpectedUserAgent {
				t.Fatalf("Incorrect user agent output: got %s, want %s", ua, tc.ExpectedUserAgent)
			}
		})
	}
}

func TestGetProjectFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceProject types.String
		ProviderProject types.String
		ExpectedProject types.String
		ExpectedError   bool
	}{
		"project is pulled from the resource config value instead of the provider config value, even if both set": {
			ResourceProject: types.StringValue("foo"),
			ProviderProject: types.StringValue("bar"),
			ExpectedProject: types.StringValue("foo"),
		},
		"project is pulled from the provider config value when unset on the resource": {
			ResourceProject: types.StringNull(),
			ProviderProject: types.StringValue("bar"),
			ExpectedProject: types.StringValue("bar"),
		},
		"error when project is not set on the provider or the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			project := getProjectFramework(tc.ResourceProject, tc.ProviderProject, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if project != tc.ExpectedProject {
				t.Fatalf("Incorrect project: got %s, want %s", project, tc.ExpectedProject)
			}
		})
	}
}

func TestGetRegionFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceRegion types.String
		ResourceZone   types.String
		ProviderRegion types.String
		ProviderZone   types.String
		ExpectedRegion types.String
		ExpectedError  bool
	}{
		"region is pulled from the resource config's region value if available": {
			ResourceRegion: types.StringValue("foo"),
			ExpectedRegion: types.StringValue("foo"),
		},
		"region is pulled from the resource config's zone value if region is unset": {
			ResourceRegion: types.StringNull(),
			ResourceZone:   types.StringValue("foo-a"),
			ExpectedRegion: types.StringValue("foo"),
		},
		"region is pulled from the provider config's region value when region and zone are unset on the resource": {
			ResourceRegion: types.StringNull(),
			ResourceZone:   types.StringNull(),
			ProviderRegion: types.StringValue("bar"),
			ExpectedRegion: types.StringValue("bar"),
		},
		"region is pulled from the provider config's zone value when region is unset on the provider (and resource config lacks region/zone)": {
			ResourceRegion: types.StringNull(),
			ResourceZone:   types.StringNull(),
			ProviderRegion: types.StringNull(),
			ProviderZone:   types.StringValue("bar-a"),
			ExpectedRegion: types.StringValue("bar"),
		},
		"error when region and zone are not set on the provider nor the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			region := getRegionFramework(tc.ResourceRegion, tc.ResourceZone, tc.ProviderRegion, tc.ProviderZone, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}

func TestGetZoneFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceRegion types.String
		ResourceZone   types.String
		ProviderRegion types.String
		ProviderZone   types.String
		ExpectedZone   types.String
		ExpectedError  bool
	}{
		"zone is pulled from the resource config": {
			ResourceZone: types.StringValue("foo"),
			ExpectedZone: types.StringValue("foo"),
		},
		"zone is pulled from the resource config's region value if available": {
			ResourceZone: types.StringNull(),
			ProviderZone: types.StringValue("bar"),
			ExpectedZone: types.StringValue("bar"),
		},
		"error when zone is not set on the provider nor the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			region := getZoneFramework(tc.ResourceZone, tc.ProviderZone, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if region != tc.ExpectedZone {
				t.Fatalf("Incorrect zone: got %s, want %s", region, tc.ExpectedZone)
			}
		})
	}
}

func TestGetLocationFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceLocation types.String
		ResourceRegion   types.String
		ResourceZone     types.String
		ProviderZone     types.String
		ExpectedLocation types.String
		ExpectedError    bool
	}{
		"location is pulled from the resource config": {
			ResourceLocation: types.StringValue("resource-location"),
			ResourceRegion:   types.StringNull(),
			ResourceZone:     types.StringNull(),
			ProviderZone:     types.StringNull(),
			ExpectedLocation: types.StringValue("resource-location"),
		},
		"region is pulled from the resource config when location is not set": {
			ResourceLocation: types.StringNull(),
			ResourceRegion:   types.StringValue("resource-region"),
			ResourceZone:     types.StringNull(),
			ProviderZone:     types.StringNull(),
			ExpectedLocation: types.StringValue("resource-region"),
		},
		"zone is pulled from the resource config when location and region is not set": {
			ResourceLocation: types.StringNull(),
			ResourceRegion:   types.StringNull(),
			ResourceZone:     types.StringValue("resource-zone"),
			ProviderZone:     types.StringNull(),
			ExpectedLocation: types.StringValue("resource-zone"),
		},
		//"zone pulled from the resource config can be retrieved by splitting on slashes and selecting last element": {
		// This behaviour isn't implemented but is present in the SDK version of getLocation.
		// Do we reproduce it? Or leave it behind?
		//},
		"zone is pulled from the provider config when location/region/zone are not set in the resource config": {
			ResourceLocation: types.StringNull(),
			ResourceRegion:   types.StringNull(),
			ResourceZone:     types.StringNull(),
			ProviderZone:     types.StringValue("provider-zone"),
			ExpectedLocation: types.StringValue("provider-zone"),
		},
		"error when neither location or region set on resource, and zone is not set on the provider": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			location := getLocationFramework(tc.ResourceLocation, tc.ResourceRegion, tc.ProviderZone, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if location != tc.ExpectedLocation {
				t.Fatalf("Incorrect location: got %s, want %s", location, tc.ExpectedLocation)
			}
		})
	}
}
