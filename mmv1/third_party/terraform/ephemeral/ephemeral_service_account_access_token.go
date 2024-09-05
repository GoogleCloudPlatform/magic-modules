package ephemeral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoogleEphemeralServiceAccountAccessToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountAccessToken{}
}

type googleEphemeralServiceAccountAccessToken struct{}

func (p *googleEphemeralServiceAccountAccessToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_access_access_token"
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
				// ValidateFunc: verify.ValidateRegexp("(" + strings.Join(verify.PossibleServiceAccountNames, "|") + ")"),
			},
			"access_token": schema.StringAttribute{
				Sensitive: true,
				Computed:  true,
			},
			"lifetime": schema.StringAttribute{
				Optional: true,
				// ValidateFunc: verify.ValidateDuration(), // duration <=3600s; TODO: support validateDuration(min,max)
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
	var data ephemeralServiceAccountAccessTokenModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result := ""
	for range data.Length.ValueInt64() {
		result += "test"
	}

	data.Result = types.StringValue(string(result))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
