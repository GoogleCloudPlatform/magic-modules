package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setStateForCollapsedObject sets nested fields as top-level state values for an object that should be collapsed.
// It takes in the flattened value of the object (i.e. converted from *obj.FooBar to []interface{}{map[string]interface{}}).
// Example Usage:
// err := setStateForFieldsInFlattenedObject(d, flattenMyTopLevelNestedObject(obj.MyTopLevelNestedObject))
func setStateForCollapsedObject(d *schema.ResourceData, v interface{}) error {
	if v == nil {
		return nil
	}

	ls, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("expected nested object value to be flattened to []interface{}")
	}
	if len(ls) == 0 {
		return nil
	}

	nestedObj := ls[0].(map[string]interface{})
	for k, kv := range nestedObj {
		if err := d.Set(k, kv); err != nil {
			return fmt.Errorf("error setting %s in state: %s", k, err)
		}
	}
	return nil
}

