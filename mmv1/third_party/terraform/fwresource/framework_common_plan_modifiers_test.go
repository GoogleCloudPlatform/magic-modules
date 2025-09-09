package fwresource

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
		resourceSchema    schema.Schema
		providerConfig    *transport_tpg.Config
		expectedAttribute types.String
		expectError       bool
		errorContains     string
		req               resource.ModifyPlanRequest
		resp              *resource.ModifyPlanResponse
	}{
		"Prioritizes set config value": {
			providerConfig: &transport_tpg.Config{
				Project: "default-provider-project",
			},
			expectedAttribute: types.StringValue("state-project"),
			req:               GenerateModifyPlanRequestWithSetValue(testSchema, "project", types.StringValue("state-project")),
			resp:              GenerateModifyPlanResponse(testSchema),
		},
		"Falls back on provider default": {
			providerConfig: &transport_tpg.Config{
				Project: "default-provider-project",
			},
			expectedAttribute: types.StringValue("default-provider-project"),
			req:               GenerateModifyPlanRequest(testSchema),
			resp:              GenerateModifyPlanResponse(testSchema),
		},
		"Errors if there is no config value or provider default": {
			providerConfig: &transport_tpg.Config{},
			req:            GenerateModifyPlanRequest(testSchema),
			resp:           GenerateModifyPlanResponse(testSchema),
			expectError:    true,
			errorContains:  "required field is not set",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			DefaultProjectModify(ctx, tc.req, tc.resp, tc.providerConfig.Project)
			if tc.resp.Diagnostics.HasError() {
				if tc.expectError {
					// Check if the error message contains the expected substring.
					if tc.errorContains != "" {
						found := false
						for _, d := range tc.resp.Diagnostics.Errors() {
							if strings.Contains(d.Detail(), tc.errorContains) {
								found = true
								break
							}
						}
						if !found {
							t.Fatalf("expected error to contain %q, but it did not. Got: %v", tc.errorContains, tc.resp.Diagnostics.Errors())
						}
					}
					// Correctly handled an expected error.
					return
				}
				t.Fatalf("unexpected error: %v", tc.resp.Diagnostics)
			}

			if tc.expectError {
				t.Fatal("expected an error, but got none")
			}

			var finalAttribute types.String
			tc.resp.Plan.GetAttribute(ctx, path.Root("project"), &finalAttribute)
			if !finalAttribute.Equal(tc.expectedAttribute) {
				t.Fatalf("incorrect attributes parsed.\n- got:  %v\n- want: %v", finalAttribute, tc.expectedAttribute)
			}
		})
	}
}
