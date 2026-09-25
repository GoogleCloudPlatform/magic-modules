package apikeys

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ ephemeral.EphemeralResource = &googleEphemeralApikeysKey{}

func init() {
	registry.FrameworkEphemeralResource{
		Name:        "google_apikeys_key",
		ProductName: "apikeys",
		Func:        GoogleEphemeralApikeysKey,
	}.Register()
}

func GoogleEphemeralApikeysKey() ephemeral.EphemeralResource {
	return &googleEphemeralApikeysKey{}
}

type googleEphemeralApikeysKey struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralApikeysKey) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikeys_key"
}

type ephemeralApikeysKeyModel struct {
	Name      types.String `tfsdk:"name"`
	Project   types.String `tfsdk:"project"`
	Id        types.String `tfsdk:"id"`
	KeyString types.String `tfsdk:"key_string"`
}

func (p *googleEphemeralApikeysKey) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This ephemeral resource provides access to an API key string.",
		MarkdownDescription: "This ephemeral resource provides access to an API key string.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The API key resource name or key id. This can be a full resource name in the format `projects/{{project}}/locations/global/keys/{{key}}`, or the final key id when `project` is set or available from provider configuration.",
				Required:    true,
			},
			"project": schema.StringAttribute{
				Description: "The project to get the API key for. If it is not provided, the provider project is used. This field is inferred when `name` is a full resource name.",
				Optional:    true,
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The full API key resource name. Format: `projects/{{project}}/locations/global/keys/{{key}}`.",
				Computed:    true,
			},
			"key_string": schema.StringAttribute{
				Description: "The encrypted and signed value held by this API key.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *googleEphemeralApikeysKey) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

	p.providerConfig = pd
}

func (p *googleEphemeralApikeysKey) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralApikeysKeyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := p.providerConfig
	project := data.Project.ValueString()
	name := data.Name.ValueString()

	fullName, project, err := canonicalAPIKeyName(name, project, config.Project)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing API key name", err.Error())
		return
	}

	url := fmt.Sprintf("%s%s/keyString", transport_tpg.BaseUrl(registry.GetProduct("apikeys"), config), fullName)
	keyResp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving API key string", err.Error())
		return
	}

	keyString, ok := keyResp["keyString"].(string)
	if !ok {
		resp.Diagnostics.AddError("Error retrieving API key string", "Response did not contain a keyString value.")
		return
	}

	data.Project = types.StringValue(project)
	data.Id = types.StringValue(fullName)
	data.KeyString = types.StringValue(keyString)

	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}

func canonicalAPIKeyName(name, project, providerProject string) (string, string, error) {
	if strings.HasPrefix(name, "projects/") {
		parts := strings.Split(name, "/")
		if len(parts) != 6 || parts[2] != "locations" || parts[3] != "global" || parts[4] != "keys" || parts[1] == "" || parts[5] == "" {
			return "", "", fmt.Errorf("expected name to match projects/{project}/locations/global/keys/{key}, got %q", name)
		}
		if project != "" && project != parts[1] {
			return "", "", fmt.Errorf("project %q does not match project %q in name", project, parts[1])
		}
		return name, parts[1], nil
	}

	if strings.Contains(name, "/") {
		return "", "", fmt.Errorf("expected name to be a key id or projects/{project}/locations/global/keys/{key}, got %q", name)
	}

	if project == "" {
		project = providerProject
	}
	if project == "" {
		return "", "", fmt.Errorf("project must be set when name is not a full resource name")
	}

	return fmt.Sprintf("projects/%s/locations/global/keys/%s", project, name), project, nil
}
