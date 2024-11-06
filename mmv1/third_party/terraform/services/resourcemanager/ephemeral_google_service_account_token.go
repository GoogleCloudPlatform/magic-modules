// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/iamcredentials/v1"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountAccessToken{}

func GoogleEphemeralServiceAccountAccessToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountAccessToken{}
}

type googleEphemeralServiceAccountAccessToken struct {
	providerConfig *fwtransport.FrameworkProviderConfig
}

func (p *googleEphemeralServiceAccountAccessToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_token"
}

type ephemeralServiceAccountAccessTokenModel struct {
	TargetServiceAccount types.String `tfsdk:"target_service_account"`
	AccessToken          types.String `tfsdk:"access_token"`
	Scopes               types.Set    `tfsdk:"scopes"`
	Delegates            types.Set    `tfsdk:"delegates"`
	Lifetime             types.String `tfsdk:"lifetime"`
}

func (p *googleEphemeralServiceAccountAccessToken) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_service_account": schema.StringAttribute{
				Required: true,
				// Validators: verify.ValidateRegexp("(" + strings.Join(verify.PossibleServiceAccountNames, "|") + ")"),
			},
			"access_token": schema.StringAttribute{
				Sensitive: true,
				Computed:  true,
			},
			"lifetime": schema.StringAttribute{
				Optional: true,
				// Validators: verify.ValidateDuration(), // duration <=3600s; TODO: support validateDuration(min,max)
				// Default:      "3600s",
			},
			"scopes": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				// ValidateFunc is not yet supported on lists or sets.
			},
			"delegates": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountAccessToken) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*fwtransport.FrameworkProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fwtransport.FrameworkProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	p.providerConfig = pd
}

func (p *googleEphemeralServiceAccountAccessToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountAccessTokenModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Lifetime.IsNull() {
		data.Lifetime = types.StringValue("3600s")
	}

	service := p.providerConfig.NewIamCredentialsClient(p.providerConfig.UserAgent)
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", data.TargetServiceAccount.ValueString())

	ScopesSetValue, diags := data.Scopes.ToSetValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var delegates []string
	if !data.Delegates.IsNull() {
		DelegatesSetValue, diags := data.Delegates.ToSetValue(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		delegates = StringSet(DelegatesSetValue)
	}

	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime:  data.Lifetime.ValueString(),
		Delegates: delegates,
		Scope:     tpgresource.CanonicalizeServiceScopes(StringSet(ScopesSetValue)),
	}

	at, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, tokenRequest).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating access token",
			fmt.Sprintf("Error generating access token: %s", err),
		)
		return
	}

	data.AccessToken = types.StringValue(at.AccessToken)
	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}

func StringSet(d basetypes.SetValue) []string {

	StringSlice := make([]string, 0)
	for _, v := range d.Elements() {
		StringSlice = append(StringSlice, v.(basetypes.StringValue).ValueString())
	}
	return StringSlice
}
