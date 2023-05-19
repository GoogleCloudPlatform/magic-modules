package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

func DataSourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.DataSourceIamPolicy(parentSpecificSchema, newUpdaterFunc, options...)
}
