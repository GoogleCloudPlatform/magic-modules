package ephemeral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Length types.Int64  `tfsdk:"length"`
	Result types.String `tfsdk:"result"`
}

func (p *googleEphemeralServiceAccountAccessToken) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a test string",
		Attributes: map[string]schema.Attribute{
			"length": schema.Int64Attribute{
				Description: "The amount of test in string desired. The minimum value for length is 1.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"result": schema.StringAttribute{
				Description: "The generated test string.",
				Computed:    true,
				Sensitive:   true,
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
