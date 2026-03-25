// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

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
