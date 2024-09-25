// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package ephemeral

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/iamcredentials/v1"
)

var (
	_ ephemeral.EphemeralResourceWithConfigure = &googleEphemeralServiceAccountAccessToken{}
	_ ephemeral.EphemeralResource              = &googleEphemeralServiceAccountAccessToken{}
)

func GoogleEphemeralServiceAccountAccessToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountAccessToken{}
}

type googleEphemeralServiceAccountAccessToken struct {
	providerConfig *fwtransport.FrameworkProviderConfig
}

func (p *googleEphemeralServiceAccountAccessToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_access_access_token"
}

func (p *googleEphemeralServiceAccountAccessToken) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*fwtransport.FrameworkProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fwtransport.FrameworkProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing project, region, zone and tokenSource
	p.providerConfig = data
}

type ephemeralServiceAccountAccessTokenModel struct {
	TargetServiceAccount types.Int64  `tfsdk:"target_service_account"`
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

func (p *googleEphemeralServiceAccountAccessToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	config := p.providerConfig
	// userAgent, err := GenerateUserAgentString(config)
	// if err != nil {
	// 	resp.Diagnostics.Append(diag.NewErrorDiagnostic("Couldnt generate User Agent String", err.Error()))
	// 	return
	// }

	service := config.NewIamCredentialsClient(config.UserAgent)
	var data ephemeralServiceAccountAccessTokenModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Lifetime.IsNull() {
		data.Lifetime = types.StringValue("3600s")
	}

	name := fmt.Sprintf("projects/-/serviceAccounts/%s", data.TargetServiceAccount.String())
	DelegatesSetValue, _ := data.Delegates.ToSetValue(ctx)
	ScopesSetValue, _ := data.Scopes.ToSetValue(ctx)
	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime:  data.Lifetime.String(),
		Delegates: StringSet(DelegatesSetValue),
		Scope:     tpgresource.CanonicalizeServiceScopes(StringSet(ScopesSetValue)),
	}

	at, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, tokenRequest).Do()
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Couldnt generate token", err.Error()))
		return
	}
	data.AccessToken = basetypes.NewStringValue(at.AccessToken)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func StringSet(d basetypes.SetValue) []string {

	StringSlice := make([]string, 0)
	for _, v := range d.Elements() {
		StringSlice = append(StringSlice, v.String())
	}
	return StringSlice
}

// TODO: investigation will be needed to determine how to best approach generating
// a userAgentString with module name
// func GenerateUserAgentString(currentUserAgent string) (string, error) {
// 	var m transport_tpg.ProviderMeta

// 	err := d.GetProviderMeta(&m)
// 	if err != nil {
// 		return currentUserAgent, err
// 	}

// 	if m.ModuleName != "" {
// 		return strings.Join([]string{currentUserAgent, m.ModuleName}, " "), nil
// 	}

// 	return currentUserAgent, nil
// }
