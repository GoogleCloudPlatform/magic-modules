package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

func ResourceIamMember(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.ResourceIamMember(parentSpecificSchema, newUpdaterFunc, resourceIdParser, options...)
}
