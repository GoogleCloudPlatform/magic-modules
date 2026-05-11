package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
)

func IamMemberListResources() []func() list.ListResource {
	return []func() list.ListResource{
		resourcemanager.newprojectIamMemberListresource,
	}
}
