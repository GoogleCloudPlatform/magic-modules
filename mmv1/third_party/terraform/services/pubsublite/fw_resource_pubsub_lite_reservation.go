package pubsublite

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
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

// GooglePubsubLiteReservationResource is the resource implementation.
type GooglePubsubLiteReservationFWResource struct {
	client         *pubsublite.Service
	providerConfig *transport_tpg.Config
}

type GooglePubsubLiteReservationModel struct {
	Id                 types.String `tfsdk:"id"`
	Project            types.String `tfsdk:"project"`
	Region             types.String `tfsdk:"region"`
	Name               types.String `tfsdk:"name"`
	ThroughputCapacity types.Int64  `tfsdk:"throughput_capacity"`
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
			"region": schema.StringAttribute{
				Description:         "The region of the Pubsub Lite Reservation.",
				MarkdownDescription: "The region of the Pubsub Lite Reservation.",
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

func (d *GooglePubsubLiteReservationFWResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GooglePubsubLiteReservationModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, d.providerConfig.UserAgent)

	obj := make(map[string]interface{})

	obj["throughputCapacity"] = data.ThroughputCapacity.ValueInt64()

	data.Project = fwresource.GetProjectFramework(data.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Region = fwresource.GetRegionFramework(data.Region, types.StringValue(d.providerConfig.Region), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	billingProject := data.Project

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = data.Project
	schemaDefaultVals.Region = data.Region

	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, d.providerConfig, "{{PubsubLiteBasePath}}projects/{{project}}/locations/{{region}}/reservations?reservationId={{name}}")
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("[DEBUG] Creating new Reservation: %#v", obj))

	headers := make(http.Header)
	res := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:    d.providerConfig,
		Method:    "POST",
		Project:   billingProject.ValueString(),
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Headers:   headers,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "create fwprovider google_pubsub_lite resource")

	// Put data in model
	data.Id = types.StringValue(fmt.Sprintf("projects/%s/locations/%s/reservations/%s", data.Project.ValueString(), data.Region.ValueString(), data.Name.ValueString()))
	data.ThroughputCapacity = types.Int64Value(res["throughputCapacity"].(int64))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, d.providerConfig.UserAgent)

	data.Project = fwresource.GetProjectFramework(data.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Region = fwresource.GetRegionFramework(data.Region, types.StringValue(d.providerConfig.Region), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	billingProject := data.Project

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = data.Project
	schemaDefaultVals.Region = data.Region

	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, d.providerConfig, "{{PubSubLiteBasePath}}projects/{{project}}/locations/{{region}}/instances/{{name}}")

	if resp.Diagnostics.HasError() {
		return
	}

	headers := make(http.Header)
	res := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:    d.providerConfig,
		Method:    "GET",
		Project:   billingProject.ValueString(),
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read fwprovider google_pubsub_lite resource")

	// Put data in model
	data.Id = types.StringValue(fmt.Sprintf("projects/%s/locations/%s/instances/%s", data.Project.ValueString(), data.Region.ValueString(), data.Name.ValueString()))
	data.ThroughputCapacity = types.Int64Value(res["throughputCapacity"].(int64))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *GooglePubsubLiteReservationFWResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GooglePubsubLiteReservationModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, d.providerConfig.UserAgent)

	obj := make(map[string]interface{})

	obj["throughputCapacity"] = plan.ThroughputCapacity.ValueInt64()

	plan.Project = fwresource.GetProjectFramework(plan.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Region = fwresource.GetRegionFramework(plan.Region, types.StringValue(d.providerConfig.Region), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	billingProject := plan.Project

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = plan.Project
	schemaDefaultVals.Region = plan.Region

	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, d.providerConfig, "{{PubSubLiteBasePath}}projects/{{project}}/locations/{{region}}/instances/{{name}}")

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("[DEBUG] Updating Reservation: %#v", obj))

	headers := make(http.Header)

	updateMask := []string{}
	if !plan.ThroughputCapacity.Equal(state.ThroughputCapacity) {
		updateMask = append(updateMask, "throughputCapacity")
	}

	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	var err error
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		resp.Diagnostics.AddError("Error when sending HTTP request: ", err.Error())
		return
	}

	res := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:    d.providerConfig,
		Method:    "PATCH",
		Project:   billingProject.ValueString(),
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Headers:   headers,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "update fwprovider google_pubsub_lite resource")

	// Put data in model
	plan.Id = types.StringValue(fmt.Sprintf("projects/%s/locations/%s/instances/%s", plan.Project.ValueString(), plan.Region.ValueString(), plan.Name.ValueString()))
	plan.ThroughputCapacity = types.Int64Value(res["throughputCapacity"].(int64))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (d *GooglePubsubLiteReservationFWResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GooglePubsubLiteReservationModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, d.providerConfig.UserAgent)

	obj := make(map[string]interface{})

	data.Project = fwresource.GetProjectFramework(data.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Region = fwresource.GetRegionFramework(data.Region, types.StringValue(d.providerConfig.Region), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	billingProject := data.Project

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = data.Project
	schemaDefaultVals.Region = data.Region

	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, d.providerConfig, "{{PubSubLiteBasePath}}projects/{{project}}/locations/{{region}}/instances/{{name}}")

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("[DEBUG] Deleting Reservation: %#v", obj))

	headers := make(http.Header)
	res := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:    d.providerConfig,
		Method:    "DELETE",
		Project:   billingProject.ValueString(),
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Headers:   headers,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[DEBUG] Deleted Reservation: %#v", res))
}
