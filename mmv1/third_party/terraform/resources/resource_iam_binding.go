package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func ResourceIamBinding(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.ResourceIamBinding(parentSpecificSchema, newUpdaterFunc, resourceIdParser, options...)
}

func expandIamCondition(v interface{}) *cloudresourcemanager.Expr {
	return tpgiamresource.ExpandIamCondition(v)
}

func flattenIamCondition(condition *cloudresourcemanager.Expr) []map[string]interface{} {
	return tpgiamresource.FlattenIamCondition(condition)
}
