package compute

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ tpgresource.ListResourceWithRawV5Schemas = &ComputeInstanceListResource{}

type ComputeInstanceListResource struct {
	tpgresource.ListResourceMetadata
}

func NewComputeInstanceListResource() list.ListResource {
	return &ComputeInstanceListResource{}
}

func (r *ComputeInstanceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "google_compute_instance"
}

func (r *ComputeInstanceListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	computeInstance := ResourceComputeInstance()
	resp.ProtoV5Schema = computeInstance.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = computeInstance.ProtoIdentitySchema(ctx)()
}

func (r *ComputeInstanceListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *ComputeInstanceListResource) ListResourceConfigSchema(ctx context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project": listschema.StringAttribute{
				Optional: true,
			},
			"zone": listschema.StringAttribute{
				Optional: true,
			},
			"filter": listschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type ComputeInstanceListModel struct {
	Project types.String `tfsdk:"project"`
	Zone    types.String `tfsdk:"zone"`
	Filter  types.String `tfsdk:"filter"`
}

func (r *ComputeInstanceListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data ComputeInstanceListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	var project string
	if !data.Project.IsNull() && !data.Project.IsUnknown() {
		project = data.Project.ValueString()
	}

	if project == "" {
		project = r.Client.Project
	}

	var zone string
	if !data.Zone.IsNull() && !data.Zone.IsUnknown() {
		zone = data.Zone.ValueString()
	}
	if zone == "" {
		zone = r.Client.Zone
	}

	filterString := ""
	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		filterString = data.Filter.ValueString()
	}

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListInstances(r.Client, project, zone, filterString, func(rd *schema.ResourceData) error {
			result := req.NewListResult(ctx)

			// flatten using the instance from the LIST call
			identity, err := rd.Identity()
			if err != nil {
				return fmt.Errorf("Error getting identity: %s", err)
			}
			err = identity.Set("name", rd.Get("name").(string))
			if err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
			instanceZone := rd.Get("zone").(string)
			err = identity.Set("zone", instanceZone)
			if err != nil {
				return fmt.Errorf("Error setting zone: %s", err)
			}
			err = identity.Set("project", project)
			if err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				return err
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				return errors.New("error setting identity")
			}
			if req.IncludeResource {
				tfTypeResource, err := rd.TfTypeResourceState()
				if err != nil {
					return err
				}
				if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
					return errors.New("error setting resource")
				}
			}
			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})
		if err != nil {
			res := req.NewListResult(ctx)
			res.Diagnostics.AddError("API Error", err.Error())
			push(res)
		}
	}
}

// ListInstances lists compute instances in a zone (instances.list) and invokes callback for each after flattenComputeInstance.
func ListInstances(config *transport_tpg.Config, project, zone, filter string, callback func(rd *schema.ResourceData) error) error {
	resourceData := ResourceComputeInstance().Data(&terraform.InstanceState{})
	if err := resourceData.Set("project", project); err != nil {
		return err
	}
	if err := resourceData.Set("zone", zone); err != nil {
		return err
	}
	url, err := tpgresource.ReplaceVars(resourceData, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances")
	if err != nil {
		return err
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(resourceData, config); err == nil {
		billingProject = bp
	}

	userAgent, err := tpgresource.GenerateUserAgentString(resourceData, config.UserAgent)
	if err != nil {
		return err
	}

	opts := transport_tpg.ListCallOptions{
		Config:         config,
		TempData:       resourceData,
		Url:            url,
		BillingProject: billingProject,
		UserAgent:      userAgent,
		Filter:         filter,
		Flattener:      flattenComputeInstance,
		Callback:       callback,
	}

	return transport_tpg.ListCall(opts)
}
