package resourcemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ ephemeral.EphemeralResource = &googleEphemeralSecretManagerSecretVersion{}

func GoogleEphemeralSecretManagerSecretVersion() ephemeral.EphemeralResource {
	return &googleEphemeralSecretManagerSecretVersion{}
}

type googleEphemeralSecretManagerSecretVersion struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralSecretManagerSecretVersion) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_manager_secret_version"
}

type ephemeralSecretManagerSecretVersionModel struct {
	Project     types.String `tfsdk:"project"`
	Secret      types.String `tfsdk:"secret"`
	Version     types.String `tfsdk:"version"`
	SecretData  types.String `tfsdk:"secret_data"`
	Name        types.String `tfsdk:"name"`
	CreateTime  types.String `tfsdk:"create_time"`
	DestroyTime types.String `tfsdk:"destroy_time"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func (p *googleEphemeralSecretManagerSecretVersion) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema.Description = "This ephemeral resource provides access to a Google Secret Manager secret version, which can be combined with a write-only attribute."
	resp.Schema.MarkdownDescription = "This ephemeral resource provides access to a Google Secret Manager secret version, which can be combined with a write-only attribute."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Arguments
			"project": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"secret": schema.StringAttribute{
				Description: "The name of the secret to access (e.g. `projects/my-project/secrets/my-secret`)",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the secret to access (e.g. `latest` or `1`). If not specified, the latest version is retrieved.",
				Optional:    true,
			},
			// Attributes
			"secret_data": schema.StringAttribute{
				Description: "The secret data of the specified secret version.",
				Computed:    true,
				Sensitive:   true,
			},
			"name": schema.StringAttribute{
				Description: "The resource name of the secret version.",
				Computed:    true,
			},
			"create_time": schema.StringAttribute{
				Description: "The time at which the secret version was created.",
				Computed:    true,
			},
			"destroy_time": schema.StringAttribute{
				Description: "The time at which the secret version was destroyed, if applicable.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the secret version is enabled.",
				Computed:    true,
			},
		},
	}
}

func (p *googleEphemeralSecretManagerSecretVersion) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	p.providerConfig = pd
}

func (p *googleEphemeralSecretManagerSecretVersion) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralSecretManagerSecretVersionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO(ramon): Implement the logic to retrieve the secret version data.

	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}
