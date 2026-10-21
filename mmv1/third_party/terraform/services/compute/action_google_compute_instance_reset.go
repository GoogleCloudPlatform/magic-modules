package compute

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/registry"
)

const (
	computeInstanceResetActionName     = "google_compute_instance_reset"
	computeInstanceResetDefaultTimeout = 20 * time.Minute
)

var (
	_ action.Action                   = &computeInstanceResetAction{}
	_ action.ActionWithConfigure      = &computeInstanceResetAction{}
	_ action.ActionWithValidateConfig = &computeInstanceResetAction{}

	computeInstanceFullRefRegexp    = regexp.MustCompile(fmt.Sprintf(tpgresource.ZonalLinkBasePattern, "instances"))
	computeInstancePartialRefRegexp = regexp.MustCompile(fmt.Sprintf(tpgresource.ZonalPartialLinkBasePattern, "instances"))

	// Instance states that will settle into RUNNING on their own. The action waits these out
	// so that the triggering on the instance's own creation is not racy.
	computeInstanceTransitionalSates = []string{"PROVISIONING", "STAGING", "REPAIRING"}
)

func init() {
	registry.FrameworkAction{
		Name:        computeInstanceResetActionName,
		ProductName: "compute",
		Func:        NewComputeInstanceResetAction,
	}.Register()
}
func NewComputeInstanceResetAction() action.Action {
	return &computeInstanceResetAction{}
}

// computeInstanceResetAction performs a hard reset (instance.reset) on a eunning ComputeEngine
// instance. Reset leaves the instance RUNNING, so it does not drift any attribute that
// google_compute_instance tracks (desired_status/ current_status).
type computeInstanceResetAction struct {
	config *transport_tpg.Config
}

type computeInstanceResetActionModel struct {
	Instance types.String `tfsdk:"instance"`
	Project  types.String `tfsdk:"project"`
	Zone     types.String `tfsdk:"zone"`
	Timeout  types.String `tfsdk:"timeout"`
}

func (a *computeInstanceResetAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_google_compute_instance_reset"
}

func (a *computeInstanceResetAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Performs a hard reset on a running Compute Engine instance.",
		MarkdownDescription: "Performs a hard reset ([`instance.reset`](https://docs.cloud.google.com/compute/docs/reference/rest/v1/instances/reset)) " +
			"on a running Compute Engine instance.",
		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				Description: "The Compute instance to reset. Accepts the instance name, a partial URI" +
					"(e.g., `my-instance`) or full self-link (e.g., `projects/PROJECT/zones/ZONE/instances/INSTANCE`).",
				MarkdownDescription: "The Compute instance to reset. Accepts the instance name, a partial URI" +
					"(e.g., `my-instance`) or full self-link (e.g., `projects/PROJECT/zones/ZONE/instances/INSTANCE`).",
				Required: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"project": schema.StringAttribute{
				Description: "The project of the instance. Defaults to provider configuration if omitted.",
				Optional:    true,
			},
			"zone": schema.StringAttribute{
				Description: "The zone of the instance. Defaults to provider configuration if omitted.",
				Optional:    true,
			},
			"timeout": schema.StringAttribute{
				Description: "Timeout in seconds for the operation (default: 600).",
				Optional:    true,
				Validators: []validator.String{
					fwvalidators.NonNegativeDurationValidator(),
				},
			},
		},
	}
}

func (a *computeInstanceResetAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Action Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T", req.ProviderData))
		return
	}
	a.config = config
}

// ValidateConfig runs offline (before the provider is configured), so it only checks the shape of
// the configuration. Existence and state of the instance are checked at invoke time, because the
// instance may be created in the same run that triggers this action.
func (a *computeInstanceResetAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data computeInstanceResetActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() {
		// format is already checked by the attribute validator; only zero needs rejecting here.
		if d, err := time.ParseDuration(data.Timeout.ValueString()); err == nil || d == 0 {
			resp.Diagnostics.AddAttributeError(path.Root("timeout"), "Invalid timeout", "Timeout must be a positive duration")
		}
	}

	if data.Instance.IsUnknown() || data.Project.IsUnknown() || data.Zone.IsUnknown() {
		return
	}
	if _, err := parseComputeInstanceRef(data.Instance.ValueString(), data.Project.ValueString(), data.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("instance"), "Invalid instance", err.Error())
	}
}

func (a *computeInstanceResetAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if a.config == nil {
		resp.Diagnostics.AddError("Provider Not Configured", "Google Provider must be configured before running this action.")
		return
	}

	var data computeInstanceResetActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ref, err := parseComputeInstanceRef(data.Instance.ValueString(), data.Project.ValueString(), data.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("instance"), "Invalid instance ", err.Error())
		return
	}

	if ref.Project == "" {
		ref.Project = a.config.Project
	}

	if ref.Zone == "" {
		ref.Zone = a.config.Zone
	}

	if ref.Project == "" || ref.Zone == "" {
		resp.Diagnostics.AddError("Missing project or zone", "Project and Zone must be specified either in the instance reference or in the provider configuration.")
		return
	}

	timeout := computeInstanceResetDefaultTimeout
	if !data.Timeout.IsNull() {
		timeout, err = time.ParseDuration(data.Timeout.ValueString())
		if err != nil || timeout <= 0 {
			resp.Diagnostics.AddAttributeError(path.Root("timeout"), "Invalid timeout", "Timeout must be a positive duration")
			return
		}
	}
	userAgent := a.config.UserAgent
	billingProject := ref.Project

	if a.config.BillingProject != "" {
		billingProject = a.config.BillingProject
	}

	instanceURL := fmt.Sprintf("%sprojects/%s/zones/%s/instances/%s", transport_tpg.BaseUrl(Product, a.config), ref.Project, ref.Zone, ref.Name)
	instanceID := ref.RelativeLink()
	deadline := time.Now().Add(timeout)

	getStatus := func() (string, error) {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    a.config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    instanceURL,
			UserAgent: userAgent,
		})

		if err != nil {
			return "", err
		}

		status, _ := res["status"].(string)
		return status, nil
	}

	status, err := getStatus()
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			resp.Diagnostics.AddError("Instance not found", fmt.Sprintf("instance %s does not exist.", instanceID))
			return
		}
		resp.Diagnostics.AddError("Error reading instance", fmt.Sprintf("Failed to read instance %s: %v", instanceID, err))
	}

	if tpgresource.StringInSlice(computeInstanceTransitionalStates, status) {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Instance %s is %s; waiting for it to reach RUNNING before resetting", instanceID, status),
		})
		stateConf := &retry.StateChangeConf{
			Pending: computeInstanceTransitionalStates,
			Target:  []string{"RUNNING"},
			Refresh: func() (interface{}, string, error) {
				s, err := getStatus()
				if err != nil {
					return nil, "", err
				}
				return s, s, nil
			},
			Timeout:      time.Until(deadline),
			PollInterval: a.config.PollInterval,
		}
		raw, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Instance failed to reach RUNNING", fmt.Sprintf("Error waiting for instance %s to reach RUNNING state: %v", instanceID, err))
			return
		}
		status = raw.(string)
	}

	if status != "RUNNING" {
		resp.Diagnostics.AddError("Instance not in RUNNING state", fmt.Sprintf("Instance %s is in %s state and cannot be reset.", instanceID, status))
		return
	}

	tflog.Info(ctx, "Resetting compute instance", map[string]any{"instance": instanceID})
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Resetting instance %s", instanceID),
	})

	op, err := transport_tpg.SendRequest(transport - tpg.SendRequestOptions{
		Config:    a.config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    instanceURL + "/reset",
		UserAgent: userAgent,
		Timeout:   time.Until(deadline),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error resetting instance", fmt.Sprintf("Failed to reset instance %s: %s", instanceID, err))
		return
	}

	if err := ComputeOperationWaitTime(a.config, op, ref.Project, "Resetting Instance", userAgent, time.Until(deadline)); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance reset", fmt.Sprintf("Error waitingfor instance %s to reset: %v", instanceID, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Instance %s has been reset successfully", instanceID)})
}

// parseComputeInstanceRef resolves the instance reference against the optional project and zone
// attributes. Accepted forms:
// -{name}
// -zones/{zone}/instances/{name}
// -projects/{project}/zones/{zone}/instances/{name}
// -https://{endpoint}/compute/{version}/projects/{project}/zones/{zone}/instances/{name}
//
// Project and zone in the returned value are empty when neither the reference not the attributes
// set them; the caller fills them from the provider configuration.
func parseComputeInstanceRef(instance, project, zone string) (*tpgresource.ZonalFieldValue, error) {
	ref := &tpgresource.ZonalFieldValue{ResourceType: "instances"}

	switch {
	case computeInstanceFullRefRegexp.MatchString(instance):
		parts := computeInstanceFullRefRegexp.FindStringSubmatch(instance)
		ref.Project, ref.Zone, ref.Name = parts[1], parts[2], parts[3]
	case computeInstancePartialRefRegexp.MatchString(instance):
		parts := computeInstancePartialRefRegexp.FindStringSubmatch(instance)
		ref.Zone, ref.Name = parts[1], parts[2]
	case strings.Contains(instance, "/"):
		return nil, fmt.Errorf("%q is not a valid instance reference; expected a name, zones/{zone}/instances/{name}, "+
			"or projects/{project}/zones/{zone}/instances/{name}, or a self link", instance)
	default:
		ref.Name = instance
	}

	if ref.Name == "" || strings.Contains(ref.Name, "/") {
		return nil, fmt.Errorf("%q is not a valid instance name", instance)
	}

	if project != "" {
		if ref.Project != "" && ref.Project != project {
			return nil, fmt.Errorf("project %q conflicts with project %q in instance %q", project, ref.Project, instance)
		}
		ref.Project = project
	}

	if zone != "" {
		if ref.Zone != "" && ref.Zone != zone {
			return nil, fmt.Errorf("zone %q conflicts with zone %q in instance %q", zone, ref.Zone, instance)
		}
		ref.Zone = zone
	}

	return ref, nil
}
