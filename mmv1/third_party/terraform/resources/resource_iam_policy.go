package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

func ResourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.ResourceIamPolicy(parentSpecificSchema, newUpdaterFunc, resourceIdParser, options...)
}
