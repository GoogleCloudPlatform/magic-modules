package resourcemanager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GoogleProjectFWDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleProjectFWDataSource{}
)

// NewGoogleProjectDataSource is a helper function to simplify the provider implementation.
func NewGoogleProjectFWDataSource() datasource.DataSource {
	return &GoogleProjectFWDataSource{}
}

// GoogleProjectDataSource is the data source implementation.
type GoogleProjectFWDataSource struct {
	client         *cloudresourcemanager.Service
	providerConfig *transport_tpg.Config
}

type GoogleProjectModel struct {
	Id        types.String `tfsdk:"id"`
	ProjectId types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	OrgId     types.String `tfsdk:"org_id"`
	FolderId  types.String `tfsdk:"folder_id"`
	Number    types.String `tfsdk:"number"`
}

// Metadata returns the data source type name.
func (d *GoogleProjectFWDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fwprovider_project"
}

func (d *GoogleProjectFWDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = p.NewResourceManagerClient(p.UserAgent)
	if resp.Diagnostics.HasError() {
		return
	}
	d.providerConfig = p
}

// Schema defines the schema for the data source.
func (d *GoogleProjectFWDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get project details.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Description:         `The project ID. Changing this forces a new project to be created.`,
				MarkdownDescription: `The project ID. Changing this forces a new project to be created.`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         `The display name of the project.`,
				MarkdownDescription: `The display name of the project.`,
				Computed:            true,
			},
			"org_id": schema.StringAttribute{
				Description:         `The numeric ID of the organization this project belongs to.`,
				MarkdownDescription: `The numeric ID of the organization this project belongs to.`,
				Computed:            true,
			},
			"folder_id": schema.StringAttribute{
				Description:         `The numeric ID of the folder this project is created under.`,
				MarkdownDescription: `The numeric ID of the folder this project is created under.`,
				Computed:            true,
			},
			"number": schema.StringAttribute{
				Description:         `The numeric identifier of the project.`,
				MarkdownDescription: `The numeric identifier of the project.`,
				Computed:            true,
			},
			// This is included for backwards compatibility with the original, SDK-implemented data source.
			"id": schema.StringAttribute{
				Description:         "Project identifier",
				MarkdownDescription: "Project identifier",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GoogleProjectFWDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleProjectModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	d.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, d.client.UserAgent)

	projectId := fmt.Sprintf("projects/%s", data.ProjectId.ValueString())
	clientResp, err := d.client.Projects.Get(data.ProjectId.ValueString()).Do()
	if err != nil {
		fwtransport.HandleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceGoogleProject %q", data.ProjectId.ValueString()), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, "read fwprovider google_project data source")

	// Put data in model
	data.Id = types.StringValue(projectId)
	data.ProjectId = types.StringValue(clientResp.ProjectId)
	switch clientResp.Parent.Type {
	case "organization":
		data.OrgId = types.StringValue(clientResp.Parent.Id)
	case "folder":
		data.FolderId = types.StringValue(clientResp.Parent.Id)
	}

	data.Number = types.StringValue(strconv.FormatInt(clientResp.ProjectNumber, 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
