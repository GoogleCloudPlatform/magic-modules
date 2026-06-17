package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func IamMemberListResources() []func() list.ListResource {
	return []func() list.ListResource{
		resourcemanager.NewProjectIamMemberListResource,
	}
}
