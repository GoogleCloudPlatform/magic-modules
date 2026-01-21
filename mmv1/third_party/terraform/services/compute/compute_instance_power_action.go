package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	defaultTimeout = 600 * time.Second
)

var (
	_ action.Action              = (*computeInstancePowerAction)(nil)
	_ action.ActionWithConfigure = (*computeInstancePowerAction)(nil)
)

func NewComputeInstancePowerAction() action.Action {
	return &computeInstancePowerAction{}
}

type computeInstancePowerAction struct {
	config *transport_tpg.Config
	client *compute.Service
}

type computeInstancePowerActionModel struct {
	Project   types.String `tfsdk:"project"`
	Zone      types.String `tfsdk:"zone"`
	Instance  types.String `tfsdk:"instance"`
	Operation types.String `tfsdk:"operation"` // start | stop | restart
	Timeout   types.Int64  `tfsdk:"timeout"`
}

func (a *computeInstancePowerAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "google_compute_instance_power"
}

func (a *computeInstancePowerAction) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	tflog.Info(ctx, "Building schema for google_compute_instance_power")
	resp.Schema = actionschema.Schema{
		Description: "Performs power operations (start, stop, restart) on a Compute Engine instance.",
		Attributes: map[string]actionschema.Attribute{
			"instance": actionschema.StringAttribute{
				Description: "The Compute Engine instance name (e.g., `my-instance`) or full self-link (e.g., `projects/PROJECT/zones/ZONE/instances/INSTANCE`).",
				Required:    true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"project": actionschema.StringAttribute{
				Description: "The project ID of the instance. Defaults to provider configuration if omitted.",
				Optional:    true,
			},
			"zone": actionschema.StringAttribute{
				Description: "The zone of the instance. Defaults to provider configuration if omitted.",
				Optional:    true,
			},
			"operation": actionschema.StringAttribute{
				Description: "The power operation to perform: `start`, `stop`, or `restart`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("start", "stop", "restart"),
				},
			},
			"timeout": actionschema.Int64Attribute{
				Description: "Timeout in seconds for the operation (default: 600).",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(30),
					int64validator.AtMost(3600),
				},
			},
		},
	}
}

func (a *computeInstancePowerAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Action Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T", req.ProviderData))
		return
	}
	a.config = cfg
	a.client = cfg.NewComputeClient("terraform-provider-google-tf-action")
}

func (a *computeInstancePowerAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if a.config == nil {
		resp.Diagnostics.AddError("Provider Not Configured", "Google Provider must be configured before running this action.")
		return
	}

	var data computeInstancePowerActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, zone, name := resolveInstanceProjectZoneName(ctx, a.config, data, data.Instance.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout := defaultTimeout
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() {
		timeout = time.Duration(data.Timeout.ValueInt64()) * time.Second
	}

	operation := data.Operation.ValueString()
	client := a.client

	// Fetch current status
	instance, err := client.Instances.Get(project, zone, name).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Instance", fmt.Sprintf("Failed to read instance %q: %v", name, err))
		return
	}

	tflog.Info(ctx, "Performing power operation", map[string]any{
		"operation": operation, "instance": name, "status": instance.Status,
	})

	var op interface{}
	switch operation {
	case "start":
		if instance.Status == "RUNNING" {
			resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Instance %s is already RUNNING", name)})
			return
		}
		op, err = client.Instances.Start(project, zone, name).Context(ctx).Do()
	case "stop":
		if instance.Status == "TERMINATED" {
			resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Instance %s is already TERMINATED", name)})
			return
		}
		op, err = client.Instances.Stop(project, zone, name).Context(ctx).Do()
	case "restart":
		op, err = client.Instances.Reset(project, zone, name).Context(ctx).Do()
	default:
		resp.Diagnostics.AddError("Invalid Operation", fmt.Sprintf("Unsupported operation: %s", operation))
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error Invoking Operation", fmt.Sprintf("Failed to perform %s on instance %q: %v", operation, name, err))
		return
	}

	if err := ComputeOperationWaitTime(a.config, op, project, fmt.Sprintf("%s Instance", operation), a.config.UserAgent, timeout); err != nil {
		resp.Diagnostics.AddError("Error Waiting for Operation", fmt.Sprintf("Failed while waiting for %s operation: %v", operation, err))
		return
	}

	// Verify final state
	instance, err = client.Instances.Get(project, zone, name).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("Error Verifying Instance State", fmt.Sprintf("Failed to verify instance %q: %v", name, err))
		return
	}

	expected := map[string]string{"start": "RUNNING", "stop": "TERMINATED", "restart": "RUNNING"}[operation]
	if instance.Status != expected {
		resp.Diagnostics.AddError("Unexpected Instance State", fmt.Sprintf("Expected %q, got %q", expected, instance.Status))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Instance %s successfully reached %s state.", name, expected),
	})
	tflog.Info(ctx, "Compute power action completed successfully", map[string]any{
		"instance": name, "operation": operation, "status": instance.Status,
	})
}

func resolveInstanceProjectZoneName(
	ctx context.Context,
	cfg *transport_tpg.Config,
	data computeInstancePowerActionModel,
	instanceRef string,
	diags *diag.Diagnostics,
) (project string, zone string, name string) {
	if strings.HasPrefix(instanceRef, "projects/") {
		parts := strings.Split(instanceRef, "/")
		if len(parts) >= 6 {
			return parts[1], parts[3], parts[5]
		}
		diags.AddError("Invalid Self-Link", fmt.Sprintf("Instance self-link %q is malformed.", instanceRef))
		return "", "", ""
	}

	project = cfg.Project
	zone = cfg.Zone
	name = instanceRef

	if !data.Project.IsNull() && !data.Project.IsUnknown() && data.Project.ValueString() != "" {
		project = data.Project.ValueString()
	}
	if !data.Zone.IsNull() && !data.Zone.IsUnknown() && data.Zone.ValueString() != "" {
		zone = data.Zone.ValueString()
	}

	if project == "" || zone == "" {
		diags.AddError("Missing Project or Zone", "Provide project and zone or use a valid self-link.")
	}

	return project, zone, name
}
