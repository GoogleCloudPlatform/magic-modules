// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// ListResourceMetadata holds provider configuration for Terraform list resources (plugin-framework list package).
// Embed it in list resource implementations and call Defaults from Configure.
type ListResourceMetadata struct {
	Client *transport_tpg.Config
}

// Defaults copies muxed provider metadata into Client. Use in ListResource.Configure
// when ListResourceData is set to *transport_tpg.Config in the framework provider.
func (m *ListResourceMetadata) Defaults(req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*transport_tpg.Config); ok {
		m.Client = c
	}
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
func ApplyListBlockConfig(ctx context.Context, req list.ListRequest, attrSchema map[string]*schema.Schema, rd *schema.ResourceData) diag.Diagnostics {
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
