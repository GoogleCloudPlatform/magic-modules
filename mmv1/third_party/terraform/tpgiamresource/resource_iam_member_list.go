// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	ListUrlFunc       func(rd *schema.ResourceData, config *transport_tpg.Config) (string, error)
	Flattener         func(item map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) error
	ItemName          string // JSON key for items array (default "items")
	ResourceNameField string // identity key filled by list API, excluded from list block
}

// IamMemberListResource lists IAM member rows by reading IAM policies on one or more policy targets.
type IamMemberListResource struct {
	tpgresource.ListResourceMetadata

	typeName          string
	memberResource    *schema.Resource
	iamResourceSchema map[string]*schema.Schema // parent-identifying fields (project, zone, name, …)
	listBlockSchema   map[string]*schema.Schema // iamResourceSchema minus ResourceNameField
	listCallConfig    IamMemberListCallConfig
	newUpdater        NewResourceIamUpdaterFunc
}

// deriveSchemas extracts parent-identifying fields from memberResource.Schema (everything
// except IamMemberBaseSchema). listBlockSchema is the same minus resourceNameField.
func deriveSchemas(memberResource *schema.Resource, resourceNameField string) (iamResourceSchema, listBlockSchema map[string]*schema.Schema) {
	iamResourceSchema = make(map[string]*schema.Schema, len(memberResource.Schema))
	listBlockSchema = make(map[string]*schema.Schema, len(memberResource.Schema))
	for k, v := range memberResource.Schema {
		if _, isBase := IamMemberBaseSchema[k]; isBase {
			continue
		}
		iamResourceSchema[k] = v
		if k != resourceNameField {
			listBlockSchema[k] = v
		}
	}
	return
}

func NewIamMemberListResource(typeName string, memberResource *schema.Resource, newUpdater NewResourceIamUpdaterFunc, listCallConfig IamMemberListCallConfig) list.ListResource {
	if memberResource.Identity == nil {
		panic("tpgiamresource: NewIamMemberListResource requires a memberResource with identity (use IamWithResourceIdentity)")
	}
	iamResourceSchema, listBlockSchema := deriveSchemas(memberResource, listCallConfig.ResourceNameField)
	return &IamMemberListResource{
		typeName:          typeName,
		memberResource:    memberResource,
		iamResourceSchema: iamResourceSchema,
		listBlockSchema:   listBlockSchema,
		listCallConfig:    listCallConfig,
		newUpdater:        newUpdater,
	}
}

func (r *IamMemberListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.typeName
}

func (r *IamMemberListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *IamMemberListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	resp.ProtoV5Schema = r.memberResource.ProtoSchema(ctx)()
	if fn := r.memberResource.ProtoIdentitySchema(ctx); fn != nil {
		resp.ProtoV5IdentitySchema = fn()
	}
}

func (r *IamMemberListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	attrs := make(map[string]listschema.Attribute, len(r.listBlockSchema))
	for name, sch := range r.listBlockSchema {
		attr := listschema.StringAttribute{Description: sch.Description}
		if sch.Required {
			attr.Required = true
		} else {
			attr.Optional = true
		}
		attrs[name] = attr
	}
	resp.Schema = listschema.Schema{Attributes: attrs}
}

// applyListBlockConfig copies list-block attributes from the Terraform config into rd.
func applyListBlockConfig(ctx context.Context, req list.ListRequest, attrSchema map[string]*schema.Schema, rd *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	for attrName := range attrSchema {
		var v types.String
		diags.Append(req.Config.GetAttribute(ctx, path.Root(attrName), &v)...)
		if diags.HasError() {
			return diags
		}
		if v.IsNull() || v.IsUnknown() {
			continue
		}
		if err := rd.Set(attrName, v.ValueString()); err != nil {
			diags.AddError("Error setting IAM resource field", fmt.Sprintf("%s: %v", attrName, err))
			return diags
		}
	}
	return diags
}

func copyIamResourceFields(dst, src *schema.ResourceData, fields map[string]*schema.Schema) {
	for k := range fields {
		_ = dst.Set(k, src.Get(k))
	}
}

// discoverPolicyTargets returns one ResourceData per GCP resource whose IAM policy should be read.
func (r *IamMemberListResource) discoverPolicyTargets(ctx context.Context, req list.ListRequest) ([]*schema.ResourceData, diag.Diagnostics) {
	baseRd := r.memberResource.TestResourceData()
	diags := applyListBlockConfig(ctx, req, r.listBlockSchema, baseRd)
	if diags.HasError() {
		return nil, diags
	}

	if r.listCallConfig.ListUrlFunc == nil {
		return []*schema.ResourceData{baseRd}, diags
	}

	listUrl, err := r.listCallConfig.ListUrlFunc(baseRd, r.Client)
	if err != nil {
		diags.AddError("Error building list URL", err.Error())
		return nil, diags
	}

	var targets []*schema.ResourceData
	if err := transport_tpg.ListCall(transport_tpg.ListCallOptions{
		Config:    r.Client,
		TempData:  baseRd,
		Url:       listUrl,
		UserAgent: r.Client.UserAgent,
		ItemName:  r.listCallConfig.ItemName,
		Flattener: r.listCallConfig.Flattener,
		Callback: func(temp *schema.ResourceData) error {
			rd := r.memberResource.TestResourceData()
			copyIamResourceFields(rd, temp, r.iamResourceSchema)
			targets = append(targets, rd)
			return nil
		},
	}); err != nil {
		diags.AddError("Error listing resources", err.Error())
		return nil, diags
	}
	return targets, diags
}

func (r *IamMemberListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	policyTargets, diags := r.discoverPolicyTargets(ctx, req)
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
			if !r.yieldPolicyMembers(ctx, req, targetRd, updater, p, &yielded, yield) {
				return
			}
		}
	}
}

// yieldPolicyMembers yields one list result per IAM binding member for a single policy target.
func (r *IamMemberListResource) yieldPolicyMembers(ctx context.Context, req list.ListRequest, targetRd *schema.ResourceData, updater ResourceIamUpdater, p *cloudresourcemanager.Policy, yielded *int64, yield func(list.ListResult) bool) bool {
	for _, binding := range p.Bindings {
		for _, mem := range binding.Members {
			if strings.HasPrefix(mem, "deleted:") {
				continue
			}
			if req.Limit > 0 && *yielded >= req.Limit {
				return true
			}
			res, err := r.buildMemberResult(ctx, req, targetRd, updater, binding, mem, p.Etag)
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
	copyIamResourceFields(rd, targetRd, r.iamResourceSchema)

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
		return list.ListResult{}, fmt.Errorf("set identity: %w", err)
	}

	if req.IncludeResource {
		tfRes, err := rd.TfTypeResourceState()
		if err != nil {
			return list.ListResult{}, fmt.Errorf("resource state: %w", err)
		}
		if err := res.Resource.Set(ctx, *tfRes); err != nil {
			return list.ListResult{}, fmt.Errorf("set resource: %w", err)
		}
	}

	res.DisplayName = fmt.Sprintf("%s %s %s", updater.DescribeResource(), binding.Role, normalized)
	return res, nil
}
