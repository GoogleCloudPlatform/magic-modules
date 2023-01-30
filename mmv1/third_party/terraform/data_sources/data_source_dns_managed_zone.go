package google

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &GoogleDnsManagedZoneDataSource{}

func NewGoogleDnsManagedZoneDataSource() datasource.DataSource {
	return &GoogleDnsManagedZoneDataSource{}
}

// GoogleDnsManagedZoneDataSource defines the data source implementation
type GoogleDnsManagedZoneDataSource struct {
	client  *dns.Service
	project types.String
}

type GoogleDnsManagedZoneModel struct {
	Id            types.String `tfsdk:"id"`
	DnsName       types.String `tfsdk:"dns_name"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ManagedZoneId types.Int64  `tfsdk:"managed_zone_id"`
	NameServers   types.List   `tfsdk:"name_servers"`
	Visibility    types.String `tfsdk:"visibility"`
	Project       types.String `tfsdk:"project"`
}

func (d *GoogleDnsManagedZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_managed_zone"
}

func (d *GoogleDnsManagedZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides access to a zone's attributes within Google Cloud DNS",

		Attributes: map[string]schema.Attribute{
			"dns_name": schema.StringAttribute{
				Computed: true,
			},

			"name": schema.StringAttribute{
				Required: true,
			},

			"description": schema.StringAttribute{
				Computed: true,
			},

			"managed_zone_id": schema.Int64Attribute{
				Computed:    true,
				Description: `Unique identifier for the resource; defined by the server.`,
			},

			"name_servers": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},

			"visibility": schema.StringAttribute{
				Computed: true,
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.
			"project": schema.StringAttribute{
				MarkdownDescription: "The ID of the project for the Google Cloud.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "DNS managed zone identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GoogleDnsManagedZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*frameworkProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *frameworkProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = p.NewDnsClient(p.userAgent, &resp.Diagnostics)
	d.project = p.project
}

func (d *GoogleDnsManagedZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleDnsManagedZoneModel
	var metaData *ProviderMetaModel
	var diags diag.Diagnostics

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client.UserAgent = generateFrameworkUserAgentString(metaData, d.client.UserAgent)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Project = getProjectFramework(data.Project, d.project, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/managedZones/%s", data.Project.ValueString(), data.Name.ValueString()))
	clientResp, err := d.client.ManagedZones.Get(data.Project.ValueString(), data.Name.ValueString()).Do()
	if err != nil {
		handleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceDnsManagedZone %q", data.Name.ValueString()), &resp.Diagnostics)
	}

	tflog.Trace(ctx, "read dns record set data source")

	data.DnsName = types.StringValue(clientResp.DnsName)
	data.Description = types.StringValue(clientResp.Description)
	data.ManagedZoneId = types.Int64Value(int64(clientResp.Id))
	data.Visibility = types.StringValue(clientResp.Visibility)
	data.NameServers, diags = types.ListValueFrom(ctx, types.StringType, clientResp.NameServers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
