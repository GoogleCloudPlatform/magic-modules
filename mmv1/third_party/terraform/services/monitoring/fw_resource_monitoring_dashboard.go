package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var (
	_ resource.Resource                = &MonitoringDashboardResource{}
	_ resource.ResourceWithConfigure   = &MonitoringDashboardResource{}
	_ resource.ResourceWithModifyPlan  = &MonitoringDashboardResource{}
	_ resource.ResourceWithImportState = &MonitoringDashboardResource{}
)

func NewMonitoringDashboardResource() resource.Resource {
	return &MonitoringDashboardResource{}
}

type MonitoringDashboardResource struct {
	providerConfig *transport_tpg.Config
}

type MonitoringDashboardResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Project             types.String `tfsdk:"project"`
	DashboardJson       Normalized   `tfsdk:"dashboard_json"`
	DashboardJsonExport types.String `tfsdk:"dashboard_json_export"`
}

func (r *MonitoringDashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fw_monitoring_dashboard"
}

func (r *MonitoringDashboardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.providerConfig = p
}

func (r *MonitoringDashboardResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	fwresource.DefaultProjectModify(ctx, req, resp, r.providerConfig.Project)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MonitoringDashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An alias from a key/cert file.",
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Description: "The ID of the project in which the resource belongs. If it is not provided, the provider project is used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dashboard_json": schema.StringAttribute{
				CustomType:  NormalizedType{},
				Description: "The JSON representation of a dashboard, following the format at https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					FWMonitoringDashboardDiffSuppress(),
				},
			},
			"dashboard_json_export": schema.StringAttribute{
				Description: "The JSON representation of a dashboard, following the format at https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards. This attribute contains computed dashboard fields not contained in the user-supplied `dashboard_json` field",
				Computed:    true,
			},
			// This is included for backwards compatibility with the original, SDK-implemented resource.
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *MonitoringDashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitoringDashboardResourceModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var project types.String
	project = fwresource.GetProjectFramework(data.Project, types.StringValue(r.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = project

	dashboardJsonProp := data.DashboardJson.ValueString()
	//resource is configured to have a supplied request json rather than the resource being configured directly in Terraform
	obj, err := structure.ExpandJsonFromString(dashboardJsonProp)
	if err != nil {
		resp.Diagnostics.AddError("Error expanding supplied JSON:", fmt.Sprintf("%s", dashboardJsonProp))
		return
	}

	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, r.providerConfig.UserAgent)
	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, r.providerConfig, "{{MonitoringBasePath}}v1/projects/{{project}}/dashboards")
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout := time.Duration(20) * time.Minute

	tflog.Trace(ctx, "Creating Monitoring Dashboard", map[string]interface{}{"url": url})
	res, err := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:               r.providerConfig,
		Method:               "POST",
		Project:              project.ValueString(),
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              createTimeout,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Successfully created Monitoring Dashboard", map[string]interface{}{"response": res})

	id := res["name"].(string)
	data.Id = types.StringValue(id)

	r.Refresh(ctx, &data, &resp.State, req, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitoringDashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonitoringDashboardResourceModel
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

	r.Refresh(ctx, &data, &resp.State, req, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitoringDashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan MonitoringDashboardResourceModel
	var metaData *fwmodels.ProviderMetaModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var project types.String
	project = fwresource.GetProjectFramework(plan.Project, types.StringValue(r.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = project

	dashboardJsonProp := plan.DashboardJson.ValueString()

	//resource is configured to have a supplied request json rather than the resource being configured directly in Terraform
	obj, err := structure.ExpandJsonFromString(dashboardJsonProp)
	if err != nil {
		resp.Diagnostics.AddError("Error expanding supplied JSON:", fmt.Sprintf("%s", dashboardJsonProp))
		return
	}

	dashboardJsonExportProp := state.DashboardJsonExport.ValueString()

	//etag must be obtained from the export object
	eObj, err := structure.ExpandJsonFromString(dashboardJsonExportProp)
	if err != nil {
		resp.Diagnostics.AddError("Error expanding supplied JSON:", fmt.Sprintf("%s", dashboardJsonProp))
		return
	}
	obj["etag"] = eObj["etag"]

	// Use provider_meta to set User-Agent
	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, r.providerConfig.UserAgent)
	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, r.providerConfig, "{{MonitoringBasePath}}"+"v1/"+state.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updating Monitoring Dashboard", map[string]interface{}{"url": url})

	updateTimeout := time.Duration(20) * time.Minute

	res, err := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:               r.providerConfig,
		Method:               "PATCH",
		Project:              project.ValueString(),
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              updateTimeout,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Successfully updated Monitoring Dashboard", map[string]interface{}{"response": res})

	r.Refresh(ctx, &plan, &resp.State, req, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitoringDashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitoringDashboardResourceModel
	var metaData *fwmodels.ProviderMetaModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var project types.String
	project = fwresource.GetProjectFramework(data.Project, types.StringValue(r.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = project

	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, r.providerConfig.UserAgent)
	url := fwtransport.ReplaceVars(ctx, req, &resp.Diagnostics, schemaDefaultVals, r.providerConfig, "{{MonitoringBasePath}}"+"v1/"+data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Deleting Monitoring Dashboard", map[string]interface{}{"url": url})

	deleteTimeout := time.Duration(20) * time.Minute

	_, _ = fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:    r.providerConfig,
		Method:    "DELETE",
		Project:   project.ValueString(),
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   deleteTimeout,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Error deleting Monitoring Domain:", fmt.Sprintf("%s", data.Id.ValueString()))
		return
	}

	tflog.Trace(ctx, "Successfully deleted Monitoring Domain.")
}

func (r *MonitoringDashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idRegexes := []string{
		"projects/(?P<project>[^/]+)/dashboards/(?P<id>[^/]+)",
		"(?P<id>[^/]+)",
	}

	var resourceSchemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resourceSchemaResp)
	if resourceSchemaResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resourceSchemaResp.Diagnostics...)
		return
	}

	parsedAttributes, diags := fwresource.ParseImportId(ctx, req, resourceSchemaResp.Schema, r.providerConfig, idRegexes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for name, value := range parsedAttributes {
		//manually construct resource id for state
		if name == "id" {
			var proj string
			pVal, _ := parsedAttributes["project"].ToTerraformValue(ctx)
			_ = pVal.As(&proj)
			var dashId string
			nVal, _ := value.ToTerraformValue(ctx)
			_ = nVal.As(&dashId)
			id := fmt.Sprintf("projects/%s/dashboards/%s",
				proj,
				dashId,
			)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), types.StringValue(id))...)
		} else {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), value)...)
		}
	}
}

func (r *MonitoringDashboardResource) Refresh(ctx context.Context, data *MonitoringDashboardResourceModel, state *tfsdk.State, req interface{}, diags *diag.Diagnostics) {
	var metaData *fwmodels.ProviderMetaModel

	var project types.String
	project = fwresource.GetProjectFramework(data.Project, types.StringValue(r.providerConfig.Project), diags)
	if diags.HasError() {
		return
	}
	var schemaDefaultVals fwtransport.DefaultVars
	schemaDefaultVals.Project = project

	userAgent := fwtransport.GenerateFrameworkUserAgentString(metaData, r.providerConfig.UserAgent)
	url := fwtransport.ReplaceVars(ctx, req, diags, schemaDefaultVals, r.providerConfig, "{{MonitoringBasePath}}"+"v1/"+data.Id.ValueString())
	if diags.HasError() {
		return
	}

	tflog.Trace(ctx, "Reading Monitoring Dashboard", map[string]interface{}{"url": url})

	res, err := fwtransport.SendRequest(fwtransport.SendRequestOptions{
		Config:               r.providerConfig,
		Method:               "GET",
		Project:              project.ValueString(),
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	}, diags)
	if diags.HasError() {
		fwtransport.HandleNotFoundError(ctx, err, state, fmt.Sprintf("MonitoringDashboard %s", data.Id.ValueString()), diags)
		return
	}

	tflog.Trace(ctx, "Successfully read Monitoring Dashboard", map[string]interface{}{"response": res})

	str, err := structure.FlattenJsonToString(res)
	if err != nil {
		diags.AddError("Error converting Dashboard:", fmt.Sprintf("%s", err))
		return
	}
	if data.DashboardJson.IsNull() || data.DashboardJson.IsUnknown() {
		data.DashboardJson = NewNormalizedValue(str)
	}
	exportStr, _ := structure.NormalizeJsonString(str)
	data.DashboardJsonExport = types.StringValue(exportStr)
}
