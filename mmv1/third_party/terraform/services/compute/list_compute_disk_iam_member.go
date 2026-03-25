// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

// diskIamMemberResource must stay aligned with provider registration for google_compute_disk_iam_member.
var diskIamMemberResource = tpgiamresource.ResourceIamMember(
	ComputeDiskIamSchema,
	ComputeDiskIamUpdaterProducer,
	ComputeDiskIdParseFunc,
	tpgiamresource.IamWithResourceIdentity(ComputeDiskIamResourceIdentityParser),
)

// NewComputeDiskIamMemberListResource returns the list implementation for google_compute_disk_iam_member.
func NewComputeDiskIamMemberListResource() list.ListResource {
	return tpgiamresource.NewIamMemberListResource(
		"google_compute_disk_iam_member",
		diskIamMemberResource,
		ComputeDiskIamSchema,
		ComputeDiskIamUpdaterProducer,
	)
}
