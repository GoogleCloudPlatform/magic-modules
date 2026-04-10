// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
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

type ListResourceMetadata struct {
	ListResourceWithRawV5Schemas

	TypeName           string
	ResourceSchema     *schema.Resource
	Client             *transport_tpg.Config
	ProjectId          string
	Region             string
	Zone               string
	IdentityAttributes []string
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
