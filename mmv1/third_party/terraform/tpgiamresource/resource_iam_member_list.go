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
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func NewIamMemberListResource(typeName string, memberResource *schema.Resource, newUpdater NewResourceIamUpdaterFunc, listCallConfig IamMemberListCallConfig) list.ListResource {
	if memberResource.Identity == nil {
		panic("tpgiamresource: NewIamMemberListResource requires a memberResource with identity (use IamWithResourceIdentity)")
	}
	iamResourceSchema, listBlockSchema := tpgresource.DeriveListSchemas(memberResource.Schema, IamMemberBaseSchema, listCallConfig.ResourceNameField)
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
	resp.Schema = tpgresource.SdkSchemaToListSchema(r.listBlockSchema)
}

// discoverPolicyTargets returns one ResourceData per GCP resource whose IAM policy should be read.
func (r *IamMemberListResource) discoverPolicyTargets(ctx context.Context, req list.ListRequest) ([]*schema.ResourceData, error) {
	baseRd := r.memberResource.TestResourceData()
	if diags := tpgresource.ApplyListBlockConfig(ctx, req, r.listBlockSchema, baseRd); diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	if r.listCallConfig.ListUrlFunc == nil {
		return []*schema.ResourceData{baseRd}, nil
	}

	listUrl, err := r.listCallConfig.ListUrlFunc(baseRd, r.Client)
	if err != nil {
		return nil, fmt.Errorf("building list URL: %w", err)
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
			tpgresource.CopyResourceDataFields(rd, temp, r.iamResourceSchema)
			targets = append(targets, rd)
			return nil
		},
	}); err != nil {
		return nil, fmt.Errorf("listing resources: %w", err)
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
	tpgresource.CopyResourceDataFields(rd, targetRd, r.iamResourceSchema)

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
