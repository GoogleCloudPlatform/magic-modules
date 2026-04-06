package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
