package fwtransport

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephemeraltypes "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

func TestBuildReplacementFunc_EphemeralOpenRequest(t *testing.T) {
	testSchema := ephemeraltypes.Schema{
		Attributes: map[string]ephemeraltypes.Attribute{
			"name": ephemeraltypes.StringAttribute{Required: true},
		},
	}

	cases := map[string]struct {
		linkTmpl    string
		configValue map[string]tftypes.Value
		expected    string
	}{
		"replaces resource-specific variable from ephemeral config": {
			linkTmpl: "projects/my-project/instances/{{name}}",
			configValue: map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, "my-instance"),
			},
			expected: "projects/my-project/instances/my-instance",
		},
		"leaves unresolvable variables empty": {
			linkTmpl: "projects/my-project/instances/{{unknown}}",
			configValue: map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, "my-instance"),
			},
			expected: "projects/my-project/instances/",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			ctx := context.Background()

			config := tfsdk.Config{
				Schema: testSchema,
				Raw: tftypes.NewValue(tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name": tftypes.String,
					},
				}, tc.configValue),
			}

			req := ephemeral.OpenRequest{Config: config}

			var diags diag.Diagnostics
			re := regexp.MustCompile("{{([%[:word:]]+)}}")
			data := DefaultVars{}

			f := BuildReplacementFunc(ctx, re, req, &diags, data, nil, tc.linkTmpl, false)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics error: %v", diags)
			}

			result := re.ReplaceAllStringFunc(tc.linkTmpl, f)
			if result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}
