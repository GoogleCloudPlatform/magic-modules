package resourcemanager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/pubsublite/v1"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &GooglePubsubLiteReservationFWResource{}
	_ resource.ResourceWithConfigure = &GooglePubsubLiteReservationFWResource{}
)

// NewGooglePubsubLiteReservationResource is a helper function to simplify the provider implementation.
func NewGooglePubsubLiteReservationFWResource() resource.Resource {
	return &GooglePubsubLiteReservationFWResource{}
}

// GooglePubsubLiteReservationResource is the data source implementation.
type GooglePubsubLiteReservationFWResource struct {
	client         *pubsublite.Service
	providerConfig *transport_tpg.Config
}

type GooglePubsubLiteReservationModel struct {
	Id        types.String `tfsdk:"id"`
	Project   types.String `tfsdk:"project"`
	Region    types.String `tfsdk:"region"`
	Name      types.String `tfsdk:"name"`
	ThroughputCapacity    types.String `tfsdk:"throughput_capacity"`
}

// Metadata returns the resource type name.
func (d *GooglePubsubLiteReservationFWResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fwprovider_pubsub_lite_reservation"
}

func (d *GooglePubsubLiteReservationFWResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	clientPubsubLite, err := pubsublite.NewService(p.Context, option.WithHTTPClient(p.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client pubsublite: %s", err)
		return nil
	}
	clientPubsubLite.UserAgent = userAgent
	clientPubsubLite.BasePath =  resourceManagerBasePath

	d.client = clientPubsubLite

	if resp.Diagnostics.HasError() {
		return
	}
	d.providerConfig = p
}

// Schema defines the schema for the data source.
func (d *GooglePubsubLiteReservationFWResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pubsub Lite Reservation resource description",

		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Description:         "The project id of the Pubsub Lite Reservation.",
				MarkdownDescription: "The project id of the Pubsub Lite Reservation.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         `The display name of the project.`,
				MarkdownDescription: `The display name of the project.`,
				Required:            true,
			},
			"throughput_capacity": schema.Int64Attribute{
				Description:         `The reserved throughput capacity. Every unit of throughput capacity is equivalent to 1 MiB/s of published messages or 2 MiB/s of subscribed messages.`,
				MarkdownDescription: `The reserved throughput capacity. Every unit of throughput capacity is equivalent to 1 MiB/s of published messages or 2 MiB/s of subscribed messages.`,
				Required:            true,
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
func (d *GooglePubsubLiteReservationFWResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GooglePubsubLiteReservationModel
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

	data.Project = fwresource.GetProjectFramework(data.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// GET Request
	service := pubsublite.NewAdminProjectsLocationsReservationsService(d.client)
	appName := fmt.Sprintf("projects/%s/locations/%s/reservations/%s", data.Project.ValueString(), data.AppId.ValueString())
	clientResp, err := service.GetConfig(appName).Do()
	if err != nil {
		fwtransport.HandleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceFirebaseAndroidAppConfig %q", data.AppId.ValueString()), &resp.Diagnostics)
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
