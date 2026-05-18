// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"context"
	"errors"
	"fmt"
	"log"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// ListResource is the base contract for Google provider list resources
// It embeds list.ListResourceWithConfigure and list.ListResourceWithRawV5Schemas
// to extend the plugin-framework list API with Configure and RawV5Schemas.
type ListResource interface {
	list.ListResourceWithConfigure
	list.ListResourceWithRawV5Schemas
}

var _ ListResource = &ListResourceMetadata{}

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
	ListResource

	TypeName string
	// SDKv2Resource is the plugin SDK v2 *schema.Resource (schema, CRUD, Identity, etc.), not only attribute definitions.
	SDKv2Resource    *schema.Resource
	Client           *transport_tpg.Config
	ProjectId        string
	Region           string
	Zone             string
	ListConfigFields []ListConfigField
}

func (listR *ListResourceMetadata) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = listR.TypeName
}

func (listR *ListResourceMetadata) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	resp.ProtoV5Schema = listR.SDKv2Resource.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = listR.SDKv2Resource.ProtoIdentitySchema(ctx)()
}

func (listR *ListResourceMetadata) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	listR.Defaults(req, resp)
}

func (listR *ListResourceMetadata) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	s, err := NewListConfigSchema(listR.ListConfigFields...)
	if err != nil {
		resp.Diagnostics.AddError("Invalid list resource configuration", err.Error())
		resp.Schema = listschema.Schema{Attributes: map[string]listschema.Attribute{}}
		return
	}
	resp.Schema = s
}

func (listR *ListResourceMetadata) Defaults(request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	c, ok := request.ProviderData.(*transport_tpg.Config)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	listR.Client = c
	listR.ProjectId = c.Project
	listR.Region = c.Region
	listR.Zone = c.Zone
}

// GetProject: list config override, else provider default project.
func (listR *ListResourceMetadata) GetProject(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return listR.ProjectId
}

// GetRegion: list config override, else provider default region.
func (listR *ListResourceMetadata) GetRegion(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return listR.Region
}

// GetZone: list config override, else provider default zone.
func (listR *ListResourceMetadata) GetZone(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return listR.Zone
}

// GetLocation: list config override, else provider default region.
func (listR *ListResourceMetadata) GetLocation(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return listR.Region
}

func SetResourceIdentityAttributes(d *schema.ResourceData, attrs map[string]interface{}) error {
	identity, err := d.Identity()
	if err != nil || identity == nil {
		log.Printf("[DEBUG] SetResourceIdentityAttributes: skipping, identity unavailable: %v", err)
	} else {
		for k, v := range attrs {
			if err := identity.Set(k, v); err != nil {
				return fmt.Errorf("error setting identity field %q: %w", k, err)
			}
		}
	}
	return nil
}

// setResourceIdentity copies identity fields from rd using SDKv2Resource.Identity.SchemaMap().
// It panics if SDKv2Resource, Identity, or the identity schema is empty (wiring error).
func (listR *ListResourceMetadata) setResourceIdentity(rd *schema.ResourceData) error {
	idSchema := listR.SDKv2Resource.Identity.SchemaMap()
	attrs := make(map[string]interface{}, len(idSchema))
	for attr := range idSchema {
		attrs[attr] = rd.Get(attr)
	}
	return SetResourceIdentityAttributes(rd, attrs)
}

// ListResultDisplayName returns the first non-empty label from rd for keys in order. Use a
// single key or several for fallbacks (e.g. display_name then email).
// it returns an error if none of the keys yield a non-empty string.
func ListResultDisplayName(rd *schema.ResourceData, keys ...string) (string, error) {
	if rd == nil {
		return "", fmt.Errorf("ListResultDisplayName: ResourceData is nil")
	}
	if len(keys) == 0 {
		return "", fmt.Errorf("ListResultDisplayName: no keys provided")
	}
	for _, k := range keys {
		v, ok := rd.GetOk(k)
		if !ok {
			continue
		}
		if s := fmt.Sprintf("%v", v); s != "" {
			return s, nil
		}
	}
	return "", fmt.Errorf("ListResultDisplayName: no non-empty value among keys %q", keys)
}

// SetResult fills list result identity from rd; if includeResource, also full resource state.
// displayNameKeys lists schema attribute names (in priority order) used to set result.DisplayName
// via ListResultDisplayName when it is still empty; omit or pass no keys to skip. Non-empty keys
// produce an error if no key yields a non-empty display label.
func (listR *ListResourceMetadata) SetResult(ctx context.Context, includeResource bool, result *list.ListResult, rd *schema.ResourceData, displayNameKeys ...string) error {
	if err := listR.setResourceIdentity(rd); err != nil {
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

	if result.DisplayName == "" && len(displayNameKeys) > 0 {
		s, err := ListResultDisplayName(rd, displayNameKeys...)
		if err != nil {
			return err
		}
		result.DisplayName = s
	}

	return nil
}

// SdkSchemaToListSchema converts an SDK schema map to a plugin-framework list schema.
// Required SDK attributes become Required list attributes; everything else becomes Optional.
func SdkSchemaToListSchema(sdkSchema map[string]*schema.Schema) listschema.Schema {
	attrs := make(map[string]listschema.Attribute, len(sdkSchema))
	for name, sch := range sdkSchema {
		attr := listschema.StringAttribute{Description: sch.Description}
		if sch.Required {
			attr.Required = true
		} else {
			attr.Optional = true
		}
		attrs[name] = attr
	}
	return listschema.Schema{Attributes: attrs}
}

// ApplyListBlockConfig reads string attributes from the Terraform list-block config
// and sets them on the given ResourceData.
func ApplyListBlockConfig(ctx context.Context, req list.ListRequest, attrSchema map[string]*schema.Schema, rd *schema.ResourceData) fwdiag.Diagnostics {
	var diags fwdiag.Diagnostics
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
			diags.AddError("Error setting list block attribute", fmt.Sprintf("%s: %v", attrName, err))
			return diags
		}
	}
	return diags
}

// CopyResourceDataFields copies the values of the given schema keys from src to dst.
func CopyResourceDataFields(dst, src *schema.ResourceData, fields map[string]*schema.Schema) {
	for k := range fields {
		_ = dst.Set(k, src.Get(k))
	}
}

// DeriveListSchemas splits a resource's full schema into two maps by stripping out baseSchema
// keys and optionally excluding resourceNameField from the list block.
//   - resourceSchema: all keys in fullSchema that are NOT in baseSchema
//   - listBlockSchema: same as resourceSchema minus resourceNameField (if non-empty)
func DeriveListSchemas(fullSchema map[string]*schema.Schema, baseSchema map[string]*schema.Schema, resourceNameField string) (resourceSchema, listBlockSchema map[string]*schema.Schema) {
	resourceSchema = make(map[string]*schema.Schema, len(fullSchema))
	listBlockSchema = make(map[string]*schema.Schema, len(fullSchema))
	for k, v := range fullSchema {
		if _, isBase := baseSchema[k]; isBase {
			continue
		}
		resourceSchema[k] = v
		if k != resourceNameField {
			listBlockSchema[k] = v
		}
	}
	return
}
