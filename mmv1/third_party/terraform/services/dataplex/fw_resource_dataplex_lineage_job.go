package dataplex

// fw_resource_dataplex_lineage_job.go
//
// Handwritten Plugin Framework resource for google_dataplex_lineage_job.
//
// This resource emits OpenLineage COMPLETE events to the GCP Data Lineage API
// (ProcessOpenLineageRunEvent) on every Create and Update.  A Dataplex Process
// entity is created on the first apply and reused on subsequent applies; each
// apply creates a new Run and LineageEvent.
//
// Key design points:
//   - ol.GenerateJobSchema()  is called directly in Schema() to produce the
//     full OL-compliant schema (job_type, ownership, inputs/outputs with all
//     dataset facet blocks, column lineage, etc.).  ~300 schema lines come
//     from the openlineage-base-resource library at zero cost.
//   - ol.BuildRunEvent() is called directly in Create() / Update().
//   - ol.BaseJobResource is NOT embedded — CRUD is written inline (MM style).
//   - ol.JobResourceBackend interface is NOT implemented — capability is a
//     private package-level function.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	lineage "cloud.google.com/go/datacatalog/lineage/apiv1"
	"cloud.google.com/go/datacatalog/lineage/apiv1/lineagepb"
	"github.com/hashicorp/terraform-provider-google/google/registry"

	"github.com/OpenLineage/openlineage/byool/terraform/openlineage-base-resource/ol"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Compile-time interface checks.
var (
	_ resource.Resource              = &DataplexLineageJobResource{}
	_ resource.ResourceWithConfigure = &DataplexLineageJobResource{}
)

const (
	producer = "https://github.com/hashicorp/terraform-provider-google/google/services/dataplex/dataplex_lineage_job"
)

func init() {
	registry.FrameworkResource{
		Name:        "dataplex_lineage_job",
		ProductName: "dataplex",
		Func:        NewDataplexLineageJobResource,
	}.Register()
}

// NewDataplexLineageJobResource is the constructor registered in the framework
// provider's Resources() function.
func NewDataplexLineageJobResource() resource.Resource {
	return &DataplexLineageJobResource{}
}

// DataplexLineageJobResource is the MM-style (no BaseJobResource embedding)
// implementation of google_dataplex_lineage_job.
//
// Compare with the standalone provider's DataplexJobResource which embeds
// ol.BaseJobResource and implements ol.JobResourceBackend.  Here we write the
// CRUD methods explicitly (~80 lines of boilerplate) while reusing the two most
// valuable ol functions: GenerateJobSchema and BuildRunEvent.
type DataplexLineageJobResource struct {
	lineageClient *lineage.Client
	project       string // default project from provider config
	location      string // default region from provider config
}

// DataplexLineageJobModel is the Terraform state struct.
//
// ol.JobResourceModel is embedded (no tfsdk tag) so the Plugin Framework
// promotes all its fields — namespace, name, description, job_type, ownership,
// source_code, source_code_location, sql, tags, inputs, outputs — directly
// into this struct's attribute namespace.  No manual field duplication needed.
type DataplexLineageJobModel struct {
	ol.JobResourceModel // embedded — ~30 tfsdk fields promoted automatically

	// MM-standard: user-configured or defaulted from provider config.
	Project  types.String `tfsdk:"project"`
	Location types.String `tfsdk:"location"`

	// Computed: populated from Dataplex API response after emission.
	ProcessName      types.String `tfsdk:"process_name"`
	RunName          types.String `tfsdk:"run_name"`
	LineageEventName types.String `tfsdk:"lineage_event_name"`
	UpdateTime       types.String `tfsdk:"update_time"`
}

// ── resource.Resource ─────────────────────────────────────────────────────────

func (r *DataplexLineageJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplex_lineage_job"
}

// Configure extracts the provider-level config (*transport_tpg.Config) and
// creates the GCP Data Lineage gRPC client authenticated via the provider's
// token source — no credentials file needed.
func (r *DataplexLineageJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.project = p.Project
	r.location = p.Region

	client, err := lineage.NewClient(ctx, option.WithTokenSource(p.TokenSource))
	if err != nil {
		resp.Diagnostics.AddError(
			"Lineage Client Error",
			fmt.Sprintf("Unable to create Data Lineage API client: %s", err),
		)
		return
	}
	r.lineageClient = client
}

// Schema calls ol.GenerateJobSchema() to produce the full OL-compliant schema,
// then merges in MM-standard (project, location) and Dataplex-computed attrs.
//
// This is the primary reuse point: the entire OL facet schema — job_type,
// ownership, source_code, source_code_location, sql, tags, plus inputs/outputs
// with symlinks, catalog, schema, column_lineage, and all other dataset facets
// — is generated by a single library call rather than being written by hand.
func (r *DataplexLineageJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// ol.GenerateJobSchema builds schema for all OL facets, enabled or not.
	// Disabled facets are included as Optional+Computed stubs so a config that
	// was written for a different provider is accepted without error.
	s := ol.GenerateJobSchema(dataplexLineageJobCapability())

	// MM-standard: project is Optional+Computed so it defaults from the provider.
	s.Attributes["project"] = schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The GCP project ID. Defaults to the project configured on the provider.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// location is Required and ForceNew — the Lineage API parent is immutable.
	s.Attributes["location"] = schema.StringAttribute{
		Required:    true,
		Description: "GCP region or multi-region where the Data Lineage API operates (e.g. 'us', 'eu', 'us-central1').",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}

	// Dataplex-computed state.
	s.Attributes["process_name"] = schema.StringAttribute{
		Computed:    true,
		Description: "Full resource name of the Dataplex Process entity created for this job.",
		PlanModifiers: []planmodifier.String{
			// Stable across applies — the same process is reused for all runs.
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	s.Attributes["run_name"] = schema.StringAttribute{
		Computed:    true,
		Description: "Full resource name of the Dataplex Run created for the most recent apply.",
	}
	s.Attributes["lineage_event_name"] = schema.StringAttribute{
		Computed:    true,
		Description: "Full resource name of the Dataplex LineageEvent created for the most recent apply.",
	}
	s.Attributes["update_time"] = schema.StringAttribute{
		Computed:    true,
		Description: "End time of the most recent Run in RFC 3339 format (or start time if the run is still in progress).",
	}

	resp.Schema = s
}

// Create reads the plan into a DataplexLineageJobModel, calls ol.BuildRunEvent
// to construct the OL RunEvent, emits it to the Dataplex Lineage API, and
// persists the resulting process/run/event names to state.
func (r *DataplexLineageJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model DataplexLineageJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ol.BuildRunEvent is reused without change.  It takes *ol.JobResourceModel
	// (promoted into DataplexLineageJobModel via embedding) and the capability.
	event := ol.NewJobEventBuilder(&resp.Diagnostics, dataplexLineageJobCapability(), producer).BuildRunEvent(&model.JobResourceModel)

	if err := r.emitEvent(ctx, &model, event); err != nil {
		resp.Diagnostics.AddError("Emission Error",
			fmt.Sprintf("Unable to emit OpenLineage event to Dataplex: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Update is identical to Create — each apply emits a new COMPLETE RunEvent.
// The same Process is reused (its name is preserved in state via UseStateForUnknown).
func (r *DataplexLineageJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model DataplexLineageJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	event := ol.NewJobEventBuilder(&resp.Diagnostics, dataplexLineageJobCapability(), producer).BuildRunEvent(&model.JobResourceModel)

	if err := r.emitEvent(ctx, &model, event); err != nil {
		resp.Diagnostics.AddError("Emission Error",
			fmt.Sprintf("Unable to emit OpenLineage event to Dataplex: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Read checks whether the Dataplex Process still exists.  If the process has
// been deleted outside Terraform (drift), RemoveResource signals a re-create.
func (r *DataplexLineageJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model DataplexLineageJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	processName := model.ProcessName.ValueString()
	if processName == "" {
		// Not yet created — nothing to read.
		return
	}

	exists, err := r.processExists(ctx, r.parent(model), processName)
	if err != nil {
		resp.Diagnostics.AddError("Read Error",
			fmt.Sprintf("Unable to check Dataplex process existence: %s", err))
		return
	}
	if !exists {
		tflog.Warn(ctx, "Dataplex process no longer exists — scheduling re-create",
			map[string]any{"process_name": processName})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Delete deletes the Dataplex Process (and all its Runs / LineageEvents) via a
// long-running operation.
func (r *DataplexLineageJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model DataplexLineageJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	processName := model.ProcessName.ValueString()
	if processName == "" || r.lineageClient == nil {
		return
	}

	tflog.Info(ctx, "Deleting Dataplex process", map[string]any{"process_name": processName})
	op, err := r.lineageClient.DeleteProcess(ctx,
		&lineagepb.DeleteProcessRequest{Name: processName})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			// Already gone — treat as success.
			return
		}
		resp.Diagnostics.AddError("Delete Error",
			fmt.Sprintf("Unable to delete Dataplex process %s: %s", processName, err))
		return
	}
	if err := op.Wait(ctx); err != nil {
		resp.Diagnostics.AddError("Delete Error",
			fmt.Sprintf("Waiting for process deletion failed: %s", err))
	}
}

// ── Private helpers ───────────────────────────────────────────────────────────

// dataplexLineageJobCapability declares which OL facets this resource supports.
// Replaces ol.JobResourceBackend.Capability() — no interface needed in MM style.
func dataplexLineageJobCapability() ol.JobCapability {
	return ol.EmptyJobCapability().
		WithFacetEnabled(
			ol.JobTypeJobFacet,
			ol.OwnershipJobFacet,
		).
		WithDatasetFacetEnabled(
			ol.SymlinksDatasetFacet,
			ol.CatalogDatasetFacet,
			ol.ColumnLineageDatasetFacet,
		)
}

// parent returns the Dataplex parent resource path for API calls.
func (r *DataplexLineageJobResource) parent(model DataplexLineageJobModel) string {
	project := model.Project.ValueString()
	if project == "" {
		project = r.project
	}
	location := model.Location.ValueString()
	return fmt.Sprintf("projects/%s/locations/%s", project, location)
}

// emitEvent marshals the OL event to a protobuf Struct, calls
// ProcessOpenLineageRunEvent, and populates the computed fields in model.
func (r *DataplexLineageJobResource) emitEvent(ctx context.Context, model *DataplexLineageJobModel, event any) error {
	// Marshal the OL RunEvent to JSON then to a protobuf Struct — the API
	// accepts a generic JSON payload, not a typed proto message.
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal OL event: %w", err)
	}
	payload := map[string]any{}
	if err := json.Unmarshal(eventJSON, &payload); err != nil {
		return fmt.Errorf("unmarshal OL event to map: %w", err)
	}

	// The Dataplex Lineage API requires the OpenLineage top-level schemaURL field.
	// ol.BuildRunEvent does not set it, so we inject the canonical OL 1.0.5 URL here.
	if _, ok := payload["schemaURL"]; !ok {
		payload["schemaURL"] = "https://openlineage.io/spec/1-0-5/OpenLineage.json#/definitions/RunEvent"
	}

	// Ensure eventType is set (must be COMPLETE for a finished job run).
	if _, ok := payload["eventType"]; !ok {
		payload["eventType"] = "COMPLETE"
	}

	// Ensure eventTime is set in ISO 8601 / RFC 3339 format.
	if _, ok := payload["eventTime"]; !ok {
		payload["eventTime"] = time.Now().UTC().Format(time.RFC3339)
	}

	// Ensure the `job` block with namespace and name is present.
	if jobRaw, ok := payload["job"]; !ok || jobRaw == nil {
		payload["job"] = map[string]any{
			"namespace": model.Namespace.ValueString(),
			"name":      model.Name.ValueString(),
		}
	}

	s, err := structpb.NewStruct(payload)
	if err != nil {
		return fmt.Errorf("convert payload to structpb.Struct: %w", err)
	}

	parent := r.parent(*model)
	tflog.Info(ctx, "Emitting OL event to Dataplex", map[string]any{
		"parent":    parent,
		"namespace": model.Namespace.ValueString(),
		"name":      model.Name.ValueString(),
	})

	resp, err := r.lineageClient.ProcessOpenLineageRunEvent(ctx,
		&lineagepb.ProcessOpenLineageRunEventRequest{
			Parent:      parent,
			OpenLineage: s,
		})
	if err != nil {
		return fmt.Errorf("ProcessOpenLineageRunEvent: %w", err)
	}

	model.ProcessName = types.StringValue(resp.Process)
	model.RunName = types.StringValue(resp.Run)
	if len(resp.LineageEvents) > 0 {
		model.LineageEventName = types.StringValue(resp.LineageEvents[0])
	} else {
		model.LineageEventName = types.StringValue("")
	}

	tflog.Info(ctx, "OL event emitted", map[string]any{
		"process_name":        resp.Process,
		"run_name":            resp.Run,
		"lineage_event_count": len(resp.LineageEvents),
	})

	r.refreshUpdateTime(ctx, model)
	return nil
}

// processExists returns true when processName is found in the parent's process list.
func (r *DataplexLineageJobResource) processExists(ctx context.Context, parent, processName string) (bool, error) {
	it := r.lineageClient.ListProcesses(ctx,
		&lineagepb.ListProcessesRequest{Parent: parent})
	for {
		p, err := it.Next()
		if err == iterator.Done {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if p.Name == processName {
			return true, nil
		}
	}
}

// refreshUpdateTime fetches the latest Run and updates UpdateTime in model.
func (r *DataplexLineageJobResource) refreshUpdateTime(ctx context.Context, model *DataplexLineageJobModel) {
	runName := model.RunName.ValueString()
	if runName == "" {
		return
	}
	run, err := r.lineageClient.GetRun(ctx, &lineagepb.GetRunRequest{Name: runName})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return
		}
		tflog.Warn(ctx, "Failed to refresh run update time",
			map[string]any{"error": err.Error(), "run_name": runName})
		return
	}
	if run.EndTime != nil {
		model.UpdateTime = types.StringValue(run.EndTime.AsTime().Format(time.RFC3339))
	} else if run.StartTime != nil {
		model.UpdateTime = types.StringValue(run.StartTime.AsTime().Format(time.RFC3339))
	}
}
