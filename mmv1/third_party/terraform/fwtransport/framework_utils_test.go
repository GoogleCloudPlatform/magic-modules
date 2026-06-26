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
	cases := map[string]struct {
		LinkTmpl    string
		ConfigAttrs map[string]tftypes.Value
		WantMatch   *regexp.Regexp
	}{
		"replaces resource-specific variable from ephemeral config": {
			LinkTmpl: "https://example.com/{{secret_id}}",
			ConfigAttrs: map[string]tftypes.Value{
				"secret_id": tftypes.NewValue(tftypes.String, "my-secret"),
			},
			WantMatch: regexp.MustCompile(`https://example\.com/my-secret`),
		},
		"leaves unresolvable variables empty": {
			LinkTmpl:    "https://example.com/{{nonexistent}}",
			ConfigAttrs: map[string]tftypes.Value{},
			WantMatch:   regexp.MustCompile(`https://example\.com/`),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			ctx := context.Background()

			attrs := map[string]tftypes.Value{}
			attrTypes := map[string]tftypes.Type{}
			schemaAttrs := map[string]ephemeraltypes.Attribute{}
			for k, v := range tc.ConfigAttrs {
				attrs[k] = v
				attrTypes[k] = v.Type()
				schemaAttrs[k] = ephemeraltypes.StringAttribute{}
			}

			rawConfig := tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, attrs)
			config := tfsdk.Config{
				Schema: ephemeraltypes.Schema{Attributes: schemaAttrs},
				Raw:    rawConfig,
			}

			req := ephemeral.OpenRequest{Config: config}

			var diags diag.Diagnostics
			f := BuildReplacementFunc(tc.LinkTmpl, nil, nil, req, ctx, &diags, false)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			result := replaceVars(tc.LinkTmpl, f)
			if !tc.WantMatch.MatchString(result) {
				t.Errorf("got %q, want match %s", result, tc.WantMatch)
			}
		})
	}
}
