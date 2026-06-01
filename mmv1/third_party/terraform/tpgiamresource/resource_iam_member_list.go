// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// IAM list resources enumerate rows for google_*_iam_member instances by reading
// IAM policies on one or more GCP resources (policy targets).
//
// When IamMemberListCallConfig.ListUrlFunc is set, List() uses transport.ListCall to
// discover multiple targets (e.g. all disks in a zone), then reads IAM for each.
// Otherwise a single target is built from the list block.

package tpgiamresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ list.ListResource = &IamMemberListResource{}
var _ list.ListResourceWithRawV5Schemas = &IamMemberListResource{}
var _ list.ListResourceWithConfigure = &IamMemberListResource{}

// IamMemberListCallConfig holds resource-specific pieces for transport.ListCall.
type IamMemberListCallConfig struct {
	ListPagesOptions    transport_tpg.ListPagesOptions
	ListURLFunc         func(rd *schema.ResourceData, config *transport_tpg.Config) (string, error)
	ParentResourceField string
	EnableRoleFilter    bool
	EnableMemberFilter  bool
}

// IamMemberListResource lists IAM member rows by reading IAM policies on one or more policy targets.
type IamMemberListResource struct {
	tpgresource.ListResourceMetadata

	typeName          string
	memberResource    *schema.Resource
	iamResourceSchema map[string]*schema.Schema // parent-identifying fields (project, zone, name, …)
	listBlockSchema   listschema.Schema
	listCallConfig    IamMemberListCallConfig
	newUpdater        NewResourceIamUpdaterFunc
	Client            *transport_tpg.Config
}

func NewIamMemberListResource(typeName string, memberResource *schema.Resource, newUpdater NewResourceIamUpdaterFunc, listCallConfig IamMemberListCallConfig) list.ListResource {
	if memberResource.Identity == nil {
		panic("tpgiamresource: NewIamMemberListResource requires a memberResource with identity (use IamWithResourceIdentity)")
	}

	listConfigFields := []tpgresource.ListConfigField{
		{
			Name: listCallConfig.ParentResourceField,
			Kind: tpgresource.ListConfigKindString,
		},
	}

	if listCallConfig.EnableRoleFilter {
		listConfigFields = append(listConfigFields, tpgresource.ListConfigField{
			Name:     "role",
			Kind:     tpgresource.ListConfigKindString,
			Optional: true,
		})
	}

	if listCallConfig.EnableMemberFilter {
		listConfigFields = append(listConfigFields, tpgresource.ListConfigField{
			Name:     "member",
			Kind:     tpgresource.ListConfigKindString,
			Optional: true,
		})
	}

	return &IamMemberListResource{
		ListResourceMetadata: tpgresource.ListResourceMetadata{
			TypeName:         typeName,
			SDKv2Resource:    memberResource,
			ListConfigFields: listConfigFields,
		},
		typeName:       typeName,
		memberResource: memberResource,
		listCallConfig: listCallConfig,
		newUpdater:     newUpdater,
	}
}

func (r *IamMemberListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.typeName
}

func (r *IamMemberListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)

	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data type",
			fmt.Sprintf("Expected *transport_tpg.Config, got %T", req.ProviderData),
		)

		return
	}
	r.Client = config
}

func (r *IamMemberListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	resp.ProtoV5Schema = r.memberResource.ProtoSchema(ctx)()
	if fn := r.memberResource.ProtoIdentitySchema(ctx); fn != nil {
		resp.ProtoV5IdentitySchema = fn()
	}
}

// discoverPolicyTargets returns one ResourceData per GCP resource whose IAM policy should be read.
func (r *IamMemberListResource) discoverPolicyTargets(ctx context.Context, req list.ListRequest) ([]*schema.ResourceData, error) {
	baseRd := r.memberResource.TestResourceData()

	var parent types.String

	diags := req.Config.GetAttribute(ctx, path.Root(r.listCallConfig.ParentResourceField), &parent)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	if !parent.IsNull() && !parent.IsUnknown() {
		if err := baseRd.Set(r.listCallConfig.ParentResourceField, parent.ValueString()); err != nil {
			return nil, fmt.Errorf("setting %s: %w", r.listCallConfig.ParentResourceField, err)
		}
	}

	var targets []*schema.ResourceData

	if r.listCallConfig.ListURLFunc == nil {
		return []*schema.ResourceData{baseRd}, nil
	}

	if r.Client == nil {
		return nil, fmt.Errorf("provider client nil")
	}

	listUrl, err := r.listCallConfig.ListURLFunc(baseRd, r.Client)
	if err != nil {
		return nil, fmt.Errorf("building list URL: %w", err)
	}
	listOpts := r.listCallConfig.ListPagesOptions
	listOpts.Config = r.Client
	listOpts.TempData = baseRd
	listOpts.Resource = r.memberResource
	listOpts.ListURL = listUrl
	listOpts.UserAgent = r.Client.UserAgent

	if listOpts.ItemName == "" {
		listOpts.ItemName = "items"
	}

	listOpts.Callback = func(rd *schema.ResourceData) error {
		targetRd := r.memberResource.TestResourceData()

		targets = append(targets, targetRd)
		return nil
	}

	if err := transport_tpg.ListPages(listOpts); err != nil {
		return nil, fmt.Errorf("listing Iam policy targets: %w", err)
	}
	return targets, nil
}

func (r *IamMemberListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	policyTargets, err := r.discoverPolicyTargets(ctx, req)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Error discovering policy targets", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	roleFilter, memberFilter, diags := r.readFilters(ctx, req)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(yield func(list.ListResult) bool) {
		var yielded int64
		for _, targetRd := range policyTargets {
			if req.Limit > 0 && yielded >= req.Limit {
				return
			}
			updater, err := r.newUpdater(targetRd, r.Client)
			if err != nil {
				res := req.NewListResult(ctx)
				res.Diagnostics.AddError("API Error", err.Error())
				if !yield(res) {
					return
				}
				continue
			}
			p, err := iamPolicyReadWithRetry(updater)
			if err != nil {
				res := req.NewListResult(ctx)
				res.Diagnostics.AddError("API Error", err.Error())
				if !yield(res) {
					return
				}
				continue
			}

			if !r.yieldPolicyMembers(ctx, req, targetRd, updater, p, roleFilter, memberFilter, &yielded, yield) {
				return
			}
		}
	}
}

// yieldPolicyMembers yields one list result per IAM binding member for a single policy target.
func (r *IamMemberListResource) yieldPolicyMembers(ctx context.Context, req list.ListRequest, targetRd *schema.ResourceData, updater ResourceIamUpdater, p *cloudresourcemanager.Policy, roleFilter string, memberFilter string, yielded *int64, yield func(list.ListResult) bool) bool {
	for _, binding := range p.Bindings {
		if roleFilter != "" && binding.Role != roleFilter {
			continue
		}
		for _, mem := range binding.Members {
			normalized := tpgresource.NormalizeIamPrincipalCasing(mem)

			if memberFilter != "" && normalized != memberFilter {
				continue
			}
			if strings.HasPrefix(mem, "deleted:") {
				continue
			}
			if req.Limit > 0 && *yielded >= req.Limit {
				return true
			}
			res, err := r.buildMemberResult(ctx, req, targetRd, updater, binding, normalized, p.Etag)
			if err != nil {
				res = req.NewListResult(ctx)
				res.Diagnostics.AddError("Error building IAM member result", err.Error())
			}
			*yielded++
			if !yield(res) {
				return false
			}
		}
	}
	return true
}

// buildMemberResult populates a ResourceData for one binding member and converts it to a ListResult.
func (r *IamMemberListResource) buildMemberResult(ctx context.Context, req list.ListRequest, targetRd *schema.ResourceData, updater ResourceIamUpdater, binding *cloudresourcemanager.Binding, member, etag string) (list.ListResult, error) {
	rd := r.memberResource.TestResourceData()
	for k := range r.iamResourceSchema {
		if v, ok := rd.GetOk(k); ok {
			if err := targetRd.Set(k, v); err != nil {
				return list.ListResult{}, fmt.Errorf("setting %s: %w", k, err)
			}
		}
	}

	normalized := tpgresource.NormalizeIamPrincipalCasing(member)
	for k, v := range map[string]interface{}{
		"role":      binding.Role,
		"member":    normalized,
		"condition": FlattenIamCondition(binding.Condition),
		"etag":      etag,
	} {
		if err := rd.Set(k, v); err != nil {
			return list.ListResult{}, fmt.Errorf("set %s: %w", k, err)
		}
	}

	id := updater.GetResourceId() + "/" + binding.Role + "/" + normalized
	if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
		id += "/" + k.String()
	}
	rd.SetId(id)

	identity, err := rd.Identity()
	if err != nil {
		return list.ListResult{}, fmt.Errorf("identity: %w", err)
	}
	condTitle := ""
	if binding.Condition != nil {
		condTitle = binding.Condition.Title
	}
	setIamMemberResourceIdentity(identity, rd, r.iamResourceSchema, binding.Role, member, condTitle)

	res := req.NewListResult(ctx)
	tfIdent, err := rd.TfTypeIdentityState()
	if err != nil {
		return list.ListResult{}, fmt.Errorf("identity state: %w", err)
	}
	if err := res.Identity.Set(ctx, *tfIdent); err != nil {
		return list.ListResult{}, fmt.Errorf("set identity: %v", err)
	}

	if req.IncludeResource {
		tfRes, err := rd.TfTypeResourceState()
		if err != nil {
			return list.ListResult{}, fmt.Errorf("resource state: %w", err)
		}
		if err := res.Resource.Set(ctx, *tfRes); err != nil {
			return list.ListResult{}, fmt.Errorf("set resource: %v", err)
		}
	}

	res.DisplayName = fmt.Sprintf("%s %s %s", updater.DescribeResource(), binding.Role, normalized)
	return res, nil
}

func (r *IamMemberListResource) readFilters(ctx context.Context, req list.ListRequest) (string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	roleFilter := ""
	memberFilter := ""

	if r.listCallConfig.EnableRoleFilter {
		var v types.String
		diags.Append(req.Config.GetAttribute(ctx, path.Root("role"), &v)...)
		if !v.IsNull() && !v.IsUnknown() {
			roleFilter = v.ValueString()
		}
	}

	if r.listCallConfig.EnableMemberFilter {
		var v types.String
		diags.Append(req.Config.GetAttribute(ctx, path.Root("member"), &v)...)
		if !v.IsNull() && !v.IsUnknown() {
			memberFilter = v.ValueString()
		}
	}

	return roleFilter, memberFilter, diags
}
