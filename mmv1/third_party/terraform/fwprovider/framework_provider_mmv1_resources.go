package fwprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

// TODO: LOOK INTO HOW WE'D GENERATE THIS THE LIST OF LISTRESOURCES
// ListResources
var generatedListResources = map[string]list.ListResource{}

var handwrittenListResources = map[string]list.ListResource{
	"google_service_account": resourcemanager.NewGoogleServiceAccountListResource(),
}
