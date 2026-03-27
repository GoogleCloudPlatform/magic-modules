// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// errIamListStop signals normal stop (limit reached or consumer closed the stream).
var errIamListStop = errors.New("iam list stop")

var (
	_ list.ListResource                 = &ComputeInstanceIamMemberListResource{}
	_ list.ListResourceWithRawV5Schemas = &ComputeInstanceIamMemberListResource{}
	_ list.ListResourceWithConfigure    = &ComputeInstanceIamMemberListResource{}
)

// ComputeInstanceIamMemberListResource lists IAM bindings for google_compute_instance_iam_member
// by walking all instances matched by project, zone, and filter (via instances.list),
// then listing IAM members on each instance.
type ComputeInstanceIamMemberListResource struct {
	tpgresource.ListResourceMetadata

	memberResource *schema.Resource
}

// NewComputeInstanceIamMemberListResource returns the list resource for google_compute_instance_iam_member.
func NewComputeInstanceIamMemberListResource() list.ListResource {
	return &ComputeInstanceIamMemberListResource{
		memberResource: tpgiamresource.ResourceIamMember(
			ComputeInstanceIamSchema,
			ComputeInstanceIamUpdaterProducer,
			ComputeInstanceIdParseFunc,
			tpgiamresource.IamWithResourceIdentity(ComputeInstanceIamResourceIdentityParser),
		),
	}
}

func (r *ComputeInstanceIamMemberListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "google_compute_instance_iam_member"
}

func (r *ComputeInstanceIamMemberListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	resp.ProtoV5Schema = r.memberResource.ProtoSchema(ctx)()
	if fn := r.memberResource.ProtoIdentitySchema(ctx); fn != nil {
		resp.ProtoV5IdentitySchema = fn()
	}
}

func (r *ComputeInstanceIamMemberListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *ComputeInstanceIamMemberListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project": listschema.StringAttribute{
				Optional:    true,
				Description: "Project ID. Defaults to the provider project when unset.",
			},
			"zone": listschema.StringAttribute{
				Optional:    true,
				Description: "Zone for instances.list. Defaults to the provider zone when unset.",
			},
			"filter": listschema.StringAttribute{
				Optional:    true,
				Description: "Filter for instances.list (Compute Engine filter expression).",
			},
		},
	}
}

// ComputeInstanceIamMemberListModel holds list block attributes.
type ComputeInstanceIamMemberListModel struct {
	Project types.String `tfsdk:"project"`
	Zone    types.String `tfsdk:"zone"`
	Filter  types.String `tfsdk:"filter"`
}

func (r *ComputeInstanceIamMemberListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	if r.Client == nil {
		stream.Results = list.ListResultsStreamDiagnostics(diag.Diagnostics{
			diag.NewErrorDiagnostic("Provider not configured", "ListResource received no provider metadata; configure the provider before listing."),
		})
		return
	}

	var data ComputeInstanceIamMemberListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	project := stringAttr(data.Project, r.Client.Project)
	zone := stringAttr(data.Zone, r.Client.Zone)
	filter := stringAttr(data.Filter, "")

	stream.Results = func(push func(list.ListResult) bool) {
		var count int64
		err := ListInstanceIdentifiers(r.Client, project, zone, filter, func(instRd *schema.ResourceData) error {
			name := instRd.Get("name").(string)
			z, _ := instRd.Get("zone").(string)
			p, _ := instRd.Get("project").(string)
			if z == "" {
				z = zone
			}
			if p == "" {
				p = project
			}
			return r.listForSingleInstance(ctx, req, p, z, name, &count, push)
		})
		if err != nil && !errors.Is(err, errIamListStop) {
			res := req.NewListResult(ctx)
			res.Diagnostics.AddError("API Error", err.Error())
			push(res)
		}
	}
}

// ListInstanceIdentifiers lists compute instances and only extracts fields needed by IAM list
// enumeration (project, zone, name). This avoids full compute-instance flattening.
func ListInstanceIdentifiers(config *transport_tpg.Config, project, zone, filter string, callback func(rd *schema.ResourceData) error) error {
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
		Flattener: func(item map[string]interface{}, d *schema.ResourceData, _ *transport_tpg.Config) error {
			name, _ := item["name"].(string)
			if name == "" {
				return fmt.Errorf("instance list item missing name")
			}
			if err := d.Set("name", name); err != nil {
				return err
			}
			if v, ok := item["zone"].(string); ok && v != "" {
				if err := d.Set("zone", tpgresource.GetResourceNameFromSelfLink(v)); err != nil {
					return err
				}
			}
			if v, ok := item["selfLink"].(string); ok && v != "" {
				if m := regexp.MustCompile(`projects/([^/]+)/`).FindStringSubmatch(v); m != nil {
					if err := d.Set("project", m[1]); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Callback: callback,
	}
	return transport_tpg.ListCall(opts)
}

func stringAttr(v types.String, def string) string {
	if v.IsNull() || v.IsUnknown() {
		return def
	}
	s := v.ValueString()
	if s == "" {
		return def
	}
	return s
}

func (r *ComputeInstanceIamMemberListResource) listForSingleInstance(ctx context.Context, req list.ListRequest, project, zone, instanceName string, count *int64, push func(list.ListResult) bool) error {
	baseRd := r.memberResource.TestResourceData()
	if err := baseRd.Set("project", project); err != nil {
		return err
	}
	if err := baseRd.Set("zone", zone); err != nil {
		return err
	}
	if err := baseRd.Set("instance_name", instanceName); err != nil {
		return err
	}
	updater, err := ComputeInstanceIamUpdaterProducer(baseRd, r.Client)
	if err != nil {
		return err
	}
	p, err := tpgiamresource.IamPolicyReadWithRetry(updater)
	if err != nil {
		return err
	}
	return yieldIamMemberRows(ctx, req, r.memberResource, baseRd, updater, p, count, push)
}

func yieldIamMemberRows(ctx context.Context, req list.ListRequest, memberResource *schema.Resource, baseRd *schema.ResourceData, updater tpgiamresource.ResourceIamUpdater, p *cloudresourcemanager.Policy, count *int64, push func(list.ListResult) bool) error {
	for _, binding := range p.Bindings {
		for _, mem := range binding.Members {
			if strings.HasPrefix(mem, "deleted:") {
				continue
			}
			if req.Limit > 0 && count != nil && *count >= req.Limit {
				return errIamListStop
			}
			rd := memberResource.TestResourceData()
			for k := range ComputeInstanceIamSchema {
				_ = rd.Set(k, baseRd.Get(k))
			}
			normalized := tpgresource.NormalizeIamPrincipalCasing(mem)
			if err := rd.Set("role", binding.Role); err != nil {
				return err
			}
			if err := rd.Set("member", normalized); err != nil {
				return err
			}
			if err := rd.Set("condition", tpgiamresource.FlattenIamCondition(binding.Condition)); err != nil {
				return err
			}
			if err := rd.Set("etag", p.Etag); err != nil {
				return err
			}
			rd.SetId(tpgiamresource.IamMemberListItemID(updater, binding.Role, mem, binding.Condition))

			res := req.NewListResult(ctx)
			if memberResource.ProtoIdentitySchema(ctx) != nil {
				identity, err := rd.Identity()
				if err != nil {
					res.Diagnostics.AddError("identity", err.Error())
					if !push(res) {
						return errIamListStop
					}
					continue
				}
				identity.Set("project", rd.Get("project"))
				identity.Set("zone", rd.Get("zone"))
				identity.Set("instance_name", rd.Get("instance_name"))
				identity.Set("role", binding.Role)
				identity.Set("member", normalized)
				if binding.Condition != nil && binding.Condition.Title != "" {
					identity.Set("condition_title", binding.Condition.Title)
				}
				tfIdent, err := rd.TfTypeIdentityState()
				if err != nil {
					res.Diagnostics.AddError("identity state", err.Error())
					if !push(res) {
						return errIamListStop
					}
					continue
				}
				if err := res.Identity.Set(ctx, *tfIdent); err != nil {
					res.Diagnostics.AddError("identity state", "error setting identity")
					if !push(res) {
						return errIamListStop
					}
					continue
				}
			}
			if req.IncludeResource {
				tfRes, err := rd.TfTypeResourceState()
				if err != nil {
					res.Diagnostics.AddError("resource state", err.Error())
					if !push(res) {
						return errIamListStop
					}
					continue
				}
				if err := res.Resource.Set(ctx, *tfRes); err != nil {
					res.Diagnostics.AddError("resource state", "error setting resource state")
					if !push(res) {
						return errIamListStop
					}
					continue
				}
			}
			res.DisplayName = fmt.Sprintf("%s %s %s", updater.DescribeResource(), binding.Role, normalized)
			if count != nil {
				*count++
			}
			if !push(res) {
				return errIamListStop
			}
		}
	}
	return nil
}
