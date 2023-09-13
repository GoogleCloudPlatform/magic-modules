package tpgdclresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

// Returns the terraform representation of a three-state boolean value represented by a pointer to bool in DCL.
func FlattenEnumBool(v interface{}) string {
	b, ok := v.(*bool)
	if !ok || b == nil {
		return ""
	}
	if *b {
		return "TRUE"
	}
	return "FALSE"
}


func FlattenContainerAwsNodePoolManagement(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}

	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	if original["node_repair"] == nil {
		transformed["node_repair"] = false
	} else {
		transformed["node_repair"] = original["node_repair"]
	}

	return []interface{}{transformed}
}