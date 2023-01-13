package google

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProviderModel describes the provider data model.
type ProviderModel struct {
	Credentials 					   types.String 	 `tfsdk:"credentials"`
	AccessToken 					   types.String 	 `tfsdk:"access_token"`
	ImpersonateServiceAccount 		   types.String 	 `tfsdk:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates types.List   	 `tfsdk:"impersonate_service_account_delegates"`
	Project 						   types.String 	 `tfsdk:"project"`
	BillingProject 					   types.String 	 `tfsdk:"billing_project"`
	Region 							   types.String 	 `tfsdk:"region"`
	Zone 							   types.String 	 `tfsdk:"zone"`
	Scopes 							   types.List   	 `tfsdk:"scopes"`
	Batching 						   ProviderBatching  `tfsdk:"batching"`
	UserProjectOverride				   types.Bool		 `tfsdk:"user_project_override"`
	RequestTimeout					   types.String 	 `tfsdk:"request_timeout"`
	RequestReason					   types.String 	 `tfsdk:"request_reason"`
}

type ProviderBatching struct {
	SendAfter 	   types.String `tfsdk:"send_after"`
	EnableBatching types.Bool 	`tfsdk:"enable_batching"`
}