// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

// projectIamMemberResource must stay aligned with provider registration for google_project_iam_member.
var projectIamMemberResource = tpgiamresource.ResourceIamMember(
	IamProjectSchema,
	NewProjectIamUpdater,
	ProjectIdParseFunc,
	tpgiamresource.IamWithBatching,
	tpgiamresource.IamWithResourceIdentity(ProjectIamResourceIdentityParser),
)

// NewProjectIamMemberListResource returns the list implementation for google_project_iam_member.
func NewProjectIamMemberListResource() list.ListResource {
	return tpgiamresource.NewIamMemberListResource(
		"google_project_iam_member",
		projectIamMemberResource,
		IamProjectSchema,
		NewProjectIamUpdater,
	)
}
