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

	Client    *transport_tpg.Config
	ProjectId string
	Region    string
	Zone      string
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

// ResolveProject returns the project from override if it is set and non-empty,
// otherwise falls back to the provider-level default.
func (r *ListResourceMetadata) GetProject(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.ProjectId
}

// ResolveRegion returns the region from override if it is set and non-empty,
// otherwise falls back to the provider-level default.
func (r *ListResourceMetadata) GetRegion(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.Region
}

// ResolveZone returns the zone from override if it is set and non-empty,
// otherwise falls back to the provider-level default.
func (r *ListResourceMetadata) GetZone(override types.String) string {
	if !override.IsNull() && !override.IsUnknown() {
		if v := override.ValueString(); v != "" {
			return v
		}
	}
	return r.Zone
}

func SetIdentityFields(ctx context.Context, result *list.ListResult, rd *schema.ResourceData, fields map[string]string) error {
	identity, err := rd.Identity()
	if err != nil {
		return fmt.Errorf("error getting identity: %s", err)
	}
	for k, v := range fields {
		if err := identity.Set(k, v); err != nil {
			return fmt.Errorf("error setting identity field %q: %s", k, err)
		}
	}
	tfTypeIdentity, err := rd.TfTypeIdentityState()
	if err != nil {
		return err
	}
	if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
		return errors.New("error setting identity")
	}
	return nil
}
