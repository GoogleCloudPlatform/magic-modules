package fwresource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestDefaultProjectModify(t *testing.T) {
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	cases := map[string]struct {
		resourceSchema     schema.Schema
		providerConfig     *transport_tpg.Config
		expectedAttribute  types.String
		expectError        bool
		errorContains      string
		state              *tfsdk.State
		plan               *tfsdk.Plan
	}{
		"Prioritizes state value": {
			resourceSchema: testSchema,
			providerConfig: &transport_tpg.Config{
				Project: "default-provider-project",
			},
			expectedAttribute: types.StringValue("state-project"),
			state: tfsdk.State{
				Schema: testSchema,
			},
			plan: tfsdk.Plan{
				Schema: testSchema,
			},
		},
		"Falls back on config default": {
			resourceSchema: testSchema,
			providerConfig: &transport_tpg.Config{
				Project: "default-provider-project",
			},
			expectedAttribute: types.StringValue("default-provider-project"),
			state: tfsdk.State{
				Schema: testSchema,
			},
			plan: tfsdk.Plan{
				Schema: testSchema,
			},
		},
		"Errors if there is no state value or config default": {
			resourceSchema: testSchema,
			providerConfig: &transport_tpg.Config{},
			state: tfsdk.State{
				Schema: testSchema,
			},
			plan: tfsdk.Plan{
				Schema: testSchema,
			},
			expectError:    true,
			errorContains:  "required field is not set",
		},
	}
	cases["Prioritizes state value"].state.SetAttribute(context.Background(), path.Root("project"), types.StringValue("state-project"))

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := resource.ModifyPlanRequest{
				State: tc.state,
			}
			resp := resource.ModifyPlanResponse{
				Plan: tc.plan,
			}
			DefaultProjectModify(ctx, req, &resp, providerConfig.Project)
			if resp.Diagnostics.HasError() {
				if tc.expectError {
					// Check if the error message contains the expected substring.
					if tc.errorContains != "" {
						found := false
						for _, d := range diags.Errors() {
							if strings.Contains(d.Detail(), tc.errorContains) {
								found = true
								break
							}
						}
						if !found {
							t.Fatalf("expected error to contain %q, but it did not. Got: %v", tc.errorContains, diags.Errors())
						}
					}
					// Correctly handled an expected error.
					return
				}
				t.Fatalf("unexpected error: %v", diags)
			}

			if tc.expectError {
				t.Fatal("expected an error, but got none")
			}

			var finalAttribute types.String
			resp.Plan.GetAttribute(ctx, path.Root("project"), &finalAttribute)
			if !finalAttribute.Equal(tc.expectedAttribute) {
				t.Fatalf("incorrect attributes parsed.\n- got:  %v\n- want: %v", parsedAttributes, tc.expectedAttributes)
			}
		})
	}
}
