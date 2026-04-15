package fwprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func listResourceFunc(lr list.ListResource) func() list.ListResource {
	return func() list.ListResource { return lr }
}

// TODO: LOOK INTO HOW WE'D GENERATE THIS THE LIST OF LISTRESOURCES
// ListResources
var generatedListResources = []func() list.ListResource{}

var handwrittenListResources = []func() list.ListResource{
	listResourceFunc(resourcemanager.NewGoogleServiceAccountListResource()),
	listResourceFunc(resourcemanager.NewGoogleProjectServiceListResource()),
}
