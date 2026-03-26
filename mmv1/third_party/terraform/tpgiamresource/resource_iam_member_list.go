// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgiamresource

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

var _ list.ListResource = &IamMemberListResource{}
var _ list.ListResourceWithRawV5Schemas = &IamMemberListResource{}
var _ list.ListResourceWithConfigure = &IamMemberListResource{}

// IamMemberListResource lists IAM member instances for one parent resource by reading its IAM policy.
// It reuses the same SDK schemas and updater as google_*_iam_member; the managed resource must define
// Terraform resource identity (IamWithResourceIdentity) so list results can expose identity state.
type IamMemberListResource struct {
	tpgresource.ListResourceMetadata

	typeName       string
	memberResource *schema.Resource
	parentSchema   map[string]*schema.Schema
	newUpdater     NewResourceIamUpdaterFunc
}

// NewIamMemberListResource returns a list.ListResource for the given IAM member schema resource.
// typeName must equal the Terraform resource type (e.g. "google_compute_disk_iam_member").
// memberResource must be the *schema.Resource returned by ResourceIamMember with identity enabled.
func NewIamMemberListResource(typeName string, memberResource *schema.Resource, parentSchema map[string]*schema.Schema, newUpdater NewResourceIamUpdaterFunc) list.ListResource {
	if memberResource.Identity == nil {
		panic("tpgiamresource: NewIamMemberListResource requires a memberResource built with ResourceIamMember and IamWithResourceIdentity")
	}
	return &IamMemberListResource{
		typeName:       typeName,
		memberResource: memberResource,
		parentSchema:   parentSchema,
		newUpdater:     newUpdater,
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
	attrs := make(map[string]listschema.Attribute, len(r.parentSchema))
	for name, sch := range r.parentSchema {
		if sch.Type != schema.TypeString {
			panic(fmt.Sprintf("tpgiamresource: list parent attribute %q must be TypeString for IAM list resources", name))
		}
		desc := sch.Description
		if sch.Required {
			attrs[name] = listschema.StringAttribute{
				Required:    true,
				Description: desc,
			}
			continue
		}
		attrs[name] = listschema.StringAttribute{
			Optional:    true,
			Description: desc,
		}
	}
	resp.Schema = listschema.Schema{Attributes: attrs}
}

func applyListParentConfig(ctx context.Context, req list.ListRequest, parentSchema map[string]*schema.Schema, rd *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	for attrName := range parentSchema {
		var v types.String
		diags.Append(req.Config.GetAttribute(ctx, path.Root(attrName), &v)...)
		if diags.HasError() {
			return diags
		}
		if v.IsNull() || v.IsUnknown() {
			continue
		}
		if err := rd.Set(attrName, v.ValueString()); err != nil {
			diags.AddError("Error setting IAM parent field", fmt.Sprintf("%s: %v", attrName, err))
			return diags
		}
	}
	return diags
}

func copyParentFields(dst, src *schema.ResourceData, parentSchema map[string]*schema.Schema) {
	for k := range parentSchema {
		_ = dst.Set(k, src.Get(k))
	}
}

func (r *IamMemberListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	baseRd := r.memberResource.TestResourceData()
	diags := applyListParentConfig(ctx, req, r.parentSchema, baseRd)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	updater, err := r.newUpdater(baseRd, r.Client)
	if err != nil {
		diags.AddError("API Error", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	p, err := iamPolicyReadWithRetry(updater)
	if err != nil {
		diags.AddError("API Error", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(yield func(list.ListResult) bool) {
		var count int64
		for _, binding := range p.Bindings {
			for _, mem := range binding.Members {
				if strings.HasPrefix(mem, "deleted:") {
					continue
				}
				if req.Limit > 0 && count >= req.Limit {
					return
				}
				rd := r.memberResource.TestResourceData()
				copyParentFields(rd, baseRd, r.parentSchema)

				normalized := tpgresource.NormalizeIamPrincipalCasing(mem)
				if err := rd.Set("role", binding.Role); err != nil {
					res := req.NewListResult(ctx)
					res.Diagnostics.AddError("internal", fmt.Errorf("set role: %w", err).Error())
					yield(res)
					return
				}
				if err := rd.Set("member", normalized); err != nil {
					res := req.NewListResult(ctx)
					res.Diagnostics.AddError("internal", fmt.Errorf("set member: %w", err).Error())
					yield(res)
					return
				}
				if err := rd.Set("condition", FlattenIamCondition(binding.Condition)); err != nil {
					res := req.NewListResult(ctx)
					res.Diagnostics.AddError("internal", fmt.Errorf("set condition: %w", err).Error())
					yield(res)
					return
				}
				if err := rd.Set("etag", p.Etag); err != nil {
					res := req.NewListResult(ctx)
					res.Diagnostics.AddError("internal", fmt.Errorf("set etag: %w", err).Error())
					yield(res)
					return
				}

				id := updater.GetResourceId() + "/" + binding.Role + "/" + normalized
				if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
					id = id + "/" + k.String()
				}
				rd.SetId(id)

				identity, err := rd.Identity()
				if err != nil {
					res := req.NewListResult(ctx)
					res.Diagnostics.AddError("identity", err.Error())
					if !yield(res) {
						return
					}
					continue
				}
				ctitle := ""
				if binding.Condition != nil {
					ctitle = binding.Condition.Title
				}
				setIamMemberResourceIdentity(identity, rd, r.parentSchema, binding.Role, mem, ctitle)

				res := req.NewListResult(ctx)
				tfIdent, err := rd.TfTypeIdentityState()
				if err != nil {
					res.Diagnostics.AddError("identity state", err.Error())
					if !yield(res) {
						return
					}
					continue
				}
				if err := res.Identity.Set(ctx, *tfIdent); err != nil {
					res.Diagnostics.AddError("identity state", errors.New("error setting identity").Error())
					if !yield(res) {
						return
					}
					continue
				}

				if req.IncludeResource {
					tfRes, err := rd.TfTypeResourceState()
					if err != nil {
						res.Diagnostics.AddError("resource state", err.Error())
						if !yield(res) {
							return
						}
						continue
					}
					if err := res.Resource.Set(ctx, *tfRes); err != nil {
						res.Diagnostics.AddError("resource state", errors.New("error setting resource").Error())
						if !yield(res) {
							return
						}
						continue
					}
				}

				res.DisplayName = fmt.Sprintf("%s %s %s", updater.DescribeResource(), binding.Role, normalized)
				count++
				if !yield(res) {
					return
				}
			}
		}
	}
}
