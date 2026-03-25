// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

// folderIamMemberResource must stay aligned with provider registration for google_folder_iam_member.
var folderIamMemberResource = tpgiamresource.ResourceIamMember(
	IamFolderSchema,
	NewFolderIamUpdater,
	FolderIdParseFunc,
	tpgiamresource.IamWithResourceIdentity(FolderIamResourceIdentityParser),
)

// NewFolderIamMemberListResource returns the list implementation for google_folder_iam_member.
func NewFolderIamMemberListResource() list.ListResource {
	return tpgiamresource.NewIamMemberListResource(
		"google_folder_iam_member",
		folderIamMemberResource,
		IamFolderSchema,
		NewFolderIamUpdater,
	)
}
