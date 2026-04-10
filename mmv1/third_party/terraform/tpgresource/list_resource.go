// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type ListResource interface {
	list.ListResourceWithConfigure
}

type ListResourceWithRawV5Schemas interface {
	ListResource

	list.ListResourceWithRawV5Schemas
}

var _ ListResourceWithRawV5Schemas = &ListResourceMetadata{}

// ListConfigFieldKind selects the Terraform type for one attribute in a list resource config block.
type ListConfigFieldKind uint8

const (
	ListConfigKindString ListConfigFieldKind = iota
	ListConfigKindBool
	ListConfigKindInt64
)

// ListConfigField describes one list-block attribute for [NewListConfigSchema].
// Define [ListResourceMetadata.ListConfigFields] explicitly; add a separate model struct
// for req.Config.Get with tfsdk tags matching Name and types matching Kind (e.g. string → types.String).
type ListConfigField struct {
	Name     string
	Kind     ListConfigFieldKind
	Optional bool // true = Optional list attribute; false = Required
}

// NewListConfigSchema builds listschema.Schema Attributes from field descriptors.
func NewListConfigSchema(fields ...ListConfigField) (listschema.Schema, error) {
	attrs := make(map[string]listschema.Attribute, len(fields))
	for _, f := range fields {
		if f.Name == "" {
			return listschema.Schema{}, fmt.Errorf("list config field has empty name")
		}
		if _, dup := attrs[f.Name]; dup {
			return listschema.Schema{}, fmt.Errorf("duplicate list config field %q", f.Name)
		}
		opt := f.Optional
		req := !opt
		switch f.Kind {
		case ListConfigKindString:
			attrs[f.Name] = listschema.StringAttribute{Optional: opt, Required: req}
		case ListConfigKindBool:
			attrs[f.Name] = listschema.BoolAttribute{Optional: opt, Required: req}
		case ListConfigKindInt64:
			attrs[f.Name] = listschema.Int64Attribute{Optional: opt, Required: req}
		default:
			return listschema.Schema{}, fmt.Errorf("unsupported list config kind for field %q", f.Name)
		}
	}
	return listschema.Schema{Attributes: attrs}, nil
}

type ListResourceMetadata struct {
	ListResourceWithRawV5Schemas

	TypeName           string
	ResourceSchema     *schema.Resource
	Client             *transport_tpg.Config
	ProjectId          string
	Region             string
	Zone               string
	IdentityAttributes []string
	ListConfigFields   []ListConfigField
}

func (r *ListResourceMetadata) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ListResourceMetadata) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	resp.ProtoV5Schema = r.ResourceSchema.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = r.ResourceSchema.ProtoIdentitySchema(ctx)()
}

func (r *ListResourceMetadata) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *ListResourceMetadata) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	s, err := NewListConfigSchema(r.ListConfigFields...)
	if err != nil {
		panic("tpgresource.ListResourceConfigSchema: " + err.Error())
	}
	resp.Schema = s
}

func (r *ListResourceMetadata) Defaults(request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	c, ok := request.ProviderData.(*transport_tpg.Config)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	r.Client = c
	r.ProjectId = c.Project
	r.Region = c.Region
	r.Zone = c.Zone
}

// GetProject: list config override, else provider default project.
func (r *ListResourceMetadata) GetProject(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.ProjectId
}

// GetRegion: list config override, else provider default region.
func (r *ListResourceMetadata) GetRegion(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.Region
}

// GetZone: list config override, else provider default zone.
func (r *ListResourceMetadata) GetZone(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.Zone
}

// GetLocation: list config override, else provider default region.
func (r *ListResourceMetadata) GetLocation(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.Region
}

// setResourceIdentity copies IdentityAttributes from rd into resource identity.
func (r *ListResourceMetadata) setResourceIdentity(rd *schema.ResourceData) error {
	identity, err := rd.Identity()
	if err != nil {
		return fmt.Errorf("error getting identity: %w", err)
	}
	for _, attr := range r.IdentityAttributes {
		if v, ok := rd.GetOk(attr); ok {
			if err := identity.Set(attr, v); err != nil {
				return fmt.Errorf("error setting identity field %q: %w", attr, err)
			}
		}
	}
	return nil
}

// SetResult fills list result identity from rd; if includeResource, also full resource state.
func (r *ListResourceMetadata) SetResult(ctx context.Context, includeResource bool, result *list.ListResult, rd *schema.ResourceData) error {
	if err := r.setResourceIdentity(rd); err != nil {
		return err
	}

	tfTypeIdentity, err := rd.TfTypeIdentityState()
	if err != nil {
		return fmt.Errorf("error converting identity state: %w", err)
	}
	if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
		return errors.New("error setting identity on list result")
	}

	if includeResource {
		tfTypeResource, err := rd.TfTypeResourceState()
		if err != nil {
			return fmt.Errorf("error converting resource state: %w", err)
		}
		if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
			return errors.New("error setting resource on list result")
		}
	}

	return nil
}
